import dataclasses
import json
import socket
import time
from argparse import ArgumentParser
from dataclasses import dataclass
from typing import List

import cv2
import torch
from confluent_kafka import Producer
from paho.mqtt import client as mqtt_client

import function.helper as helper
import function.utils_rotate as utils_rotate


@dataclass
class LicensePlateDetection:
    parking_lot_id: int
    license_plate: str
    entrance: bool


class LicensePlateNotifier:
    def notify(self, detection: LicensePlateDetection):
        raise NotImplementedError


class LicensePlateHandler:
    def __init__(self, record_entrance: bool, parking_lot_id: int, notifiers: List[LicensePlateNotifier],
                 min_req_duplicates=5):
        self.parking_lot_id = parking_lot_id
        self.record_entrance = record_entrance
        self.notifiers = notifiers
        self.min_req_duplicates = min_req_duplicates
        self.known_lp = {}

    def add_all(self, license_plates: list):
        for lp in license_plates:
            # ignore license plate if it doesn't meet the formal requirements:
            if len(lp) not in [9, 10]:
                continue

            if lp not in self.known_lp:
                dup_count = 0
                notified = False
            else:
                dup_count = self.known_lp[lp]['dup_count'] + 1
                notified = self.known_lp[lp]['notified']
            self.known_lp[lp] = {'lastseen': time.time(), 'dup_count': dup_count, 'notified': notified}

        self.handle_known_lps()

    def handle_known_lps(self):
        known_lp_updated = self.known_lp.copy()
        for lp, info in self.known_lp.items():
            if info['notified'] is False and info['dup_count'] >= self.min_req_duplicates:
                # notify new record:
                for notifier in self.notifiers:
                    detection = LicensePlateDetection(parking_lot_id=self.parking_lot_id, license_plate=lp,
                                                      entrance=self.record_entrance)
                    notifier.notify(detection)
                info['notified'] = True

            # cleanup old records:
            if info['lastseen'] < time.time() - 60:
                known_lp_updated.pop(lp)
            else:
                known_lp_updated[lp] = info

        self.known_lp = known_lp_updated


# noinspection DuplicatedCode
class LicensePlateDetector:

    def __init__(self, video_src: int | str, record_entrance: bool, parking_lot_id: int,
                 notifiers: List[LicensePlateNotifier], frames_per_second=10):
        """
        :param video_src specify the source of video input.
        Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), or video file (vid.mp4)
        """
        self.video_src = video_src
        self.notifiers = notifiers
        self.frames_per_second = frames_per_second
        self.vid = None
        self.yolo_LP_detect = None
        self.yolo_license_plate = None
        self.load_model()
        self.lp_handler = LicensePlateHandler(record_entrance, parking_lot_id, notifiers)

    def load_model(self):
        self.yolo_LP_detect = torch.hub.load('yolov5', 'custom', path='model/LP_detector_nano_61.pt', force_reload=True,
                                             source='local')
        self.yolo_license_plate = torch.hub.load('yolov5', 'custom', path='model/LP_ocr_nano_62.pt', force_reload=True,
                                                 source='local')
        self.yolo_license_plate.conf = 0.60

    def run(self):
        self.init_video()
        print(f"Starting detection from video source '{self.video_src}'")
        try:
            self.process_video()
        except KeyboardInterrupt:
            print("Program terminated by KeyboardInterrupt")
        finally:
            self.vid.release()

    def init_video(self):
        self.vid = cv2.VideoCapture(self.video_src)

    def process_video(self):
        ct = 0
        skip_frames = 5
        while True:
            ct += 1
            self.vid.grab()
            if ct % skip_frames != 0:
                continue
            frame = self.read_video_frame()
            detected_license_plates = self.extract_license_plates(frame)
            self.lp_handler.add_all(detected_license_plates)

    def extract_license_plates(self, frame):
        result = []
        try:
            plates = self.yolo_LP_detect(frame, size=640)
            list_plates = plates.pandas().xyxy[0].values.tolist()
        except Exception as e:
            print("Exception during detection", e)
            return result
        for plate in list_plates:
            flag = 0
            x = int(plate[0])  # xmin
            y = int(plate[1])  # ymin
            w = int(plate[2] - plate[0])  # xmax - xmin
            h = int(plate[3] - plate[1])  # ymax - ymin
            crop_img = frame[y:y + h, x:x + w]
            for cc in range(0, 2):
                for ct in range(0, 2):
                    lp = helper.read_plate(self.yolo_license_plate, utils_rotate.deskew(crop_img, cc, ct))
                    if lp != "unknown":
                        result.append(lp)
                        flag = 1
                        break
                if flag == 1:
                    break

        return result

    def read_video_frame(self):
        retries = 0
        while True:
            try:
                ret, frame = self.vid.read()
                if frame is not None:
                    return frame
            except Exception as e:
                print("Exception while reading video:", e)

            time.sleep(0.01)
            retries += 1
            if retries > 10:
                print("Reinitializing video capture due to too many failed read operations.")
                self.vid.release()
                self.init_video()
                retries = 0


class MqttNotifier(LicensePlateNotifier):
    TOPIC = "/vn/hust/iot/lp-detector"

    def __init__(self, broker, port, client_id):
        self.client = self.create_client(broker, client_id, port)

    @staticmethod
    def create_client(broker, client_id, port):
        """
        Code in this method is inspired by https://www.emqx.com/en/blog/how-to-use-mqtt-in-python
        """
        def on_connect(this_client, userdata, flags, rc):
            if rc == 0:
                print("Connected to MQTT Broker!")
            else:
                print(f"Failed to connect, return code {rc}")

        def on_disconnect(this_client, userdata, rc):
            print(f"Disconnected with result code: {rc}")
            reconnect_count, reconnect_delay = 0, 0.5
            while reconnect_count < 5:
                print(f"Reconnecting in {reconnect_delay} seconds...")
                time.sleep(reconnect_delay)

                try:
                    this_client.reconnect()
                    print("Reconnected successfully!")
                    return
                except Exception as err:
                    print(f"{err}. Reconnect failed. Retrying...")

                reconnect_delay *= 1.5
                reconnect_delay = min(reconnect_delay, 5)
                reconnect_count += 1
            print(f"Reconnect failed after {reconnect_count} attempts. Exiting...")

        client = mqtt_client.Client(client_id)
        client.on_connect = on_connect
        client.on_disconnect = on_disconnect
        client.connect(broker, port)
        return client

    def notify(self, detection: LicensePlateDetection):
        payload = json.dumps(dataclasses.asdict(detection))
        result = self.client.publish(self.TOPIC, payload)
        status = result[0]
        if status == 0:
            print(f"Send `{payload}` to topic `{self.TOPIC}`")
        else:
            print(f"Failed to send message to topic {self.TOPIC}")


class KafkaNotifier(LicensePlateNotifier):
    TOPIC = "quickstart-events"

    def __init__(self, server):
        self.server = server

    def notify(self, detection: LicensePlateDetection):
        conf = {'bootstrap.servers': self.server,
                'client.id': socket.gethostname()}

        payload = json.dumps(dataclasses.asdict(detection))
        print(f"[Kafka] Send to {self.TOPIC}: {payload}")
        producer = Producer(conf)
        producer.produce(self.TOPIC, payload)
        producer.flush()


class ConsoleNotifier(LicensePlateNotifier):
    def notify(self, detection):
        print(f"Detected {detection.license_plate}")


def parse_arguments():
    arg_parser = ArgumentParser()
    arg_parser.add_argument('-v', '--video-src', default=0,
                            help="Source of video input. "
                                 "Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), "
                                 "or video file (vid.mp4)")
    arg_parser.add_argument('-ks', '--kafka-server', help="Kafka Server")
    arg_parser.add_argument('-mb', '--mqtt-broker', help="Hostname of MQTT broker",
                            default="broker.hivemq.com")
    arg_parser.add_argument('-mp', '--mqtt-port', help="Port of MQTT broker",
                            default=1883, type=int)
    arg_parser.add_argument('-mi', '--mqtt-client-id', help="MQTT client-ID",
                            default="iot-lp-detector")
    arg_parser.add_argument('-e', '--record-entrance', help="Record entrance (true) or exit (false)",
                            default=True, type=bool)
    arg_parser.add_argument('-p', '--parking-lot-id', help="ID of the corresponding parking lot",
                            default=0, type=int)
    return arg_parser.parse_args()


if __name__ == '__main__':
    args = parse_arguments()
    console_notifier = ConsoleNotifier()
    notifiers = [console_notifier]
    if args.kafka_server:
        notifiers.append(KafkaNotifier(server=args.kafka_server))
    elif args.mqtt_broker:
        notifiers.append(MqttNotifier(broker=args.mqtt_broker, port=args.mqtt_port, client_id=args.mqtt_client_id))
    print("Using notifiers:", notifiers)
    detector = LicensePlateDetector(video_src=args.video_src,
                                    record_entrance=args.record_entrance,
                                    parking_lot_id=args.parking_lot_id,
                                    notifiers=notifiers)
    detector.run()

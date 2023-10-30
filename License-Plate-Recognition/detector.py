import time
from argparse import ArgumentParser
from typing import List

import cv2
import torch
from paho.mqtt import client as mqtt_client

import function.helper as helper
import function.utils_rotate as utils_rotate


class LicensePlateNotifier:
    def notify(self, licence_plate):
        raise NotImplementedError


# noinspection DuplicatedCode
class LicensePlateDetector:

    def __init__(self, video_src: int | str, notifiers: List[LicensePlateNotifier], frames_per_second=10):
        """
        :param video_src specify the source of video input.
        Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), or video file (vid.mp4)
        """
        self.video_src = video_src
        self.notifiers = notifiers
        self.frames_per_second = frames_per_second
        self.yolo_LP_detect = None
        self.yolo_license_plate = None
        self.load_model()

    def load_model(self):
        self.yolo_LP_detect = torch.hub.load('yolov5', 'custom', path='model/LP_detector_nano_61.pt', force_reload=True,
                                             source='local')
        self.yolo_license_plate = torch.hub.load('yolov5', 'custom', path='model/LP_ocr_nano_62.pt', force_reload=True,
                                                 source='local')
        self.yolo_license_plate.conf = 0.60

    def run(self):
        vid = cv2.VideoCapture(self.video_src)
        print(f"Starting detection from video source '{self.video_src}'")
        try:
            self.process_video(vid)
        except KeyboardInterrupt:
            print("Program terminated by KeyboardInterrupt")
        finally:
            vid.release()

    def process_video(self, vid):
        while vid.isOpened():
            time.sleep(1 / self.frames_per_second)
            ret, frame = vid.read()

            plates = self.yolo_LP_detect(frame, size=640)
            list_plates = plates.pandas().xyxy[0].values.tolist()
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
                            for n in self.notifiers:
                                n.notify(lp)
                            flag = 1
                            break
                    if flag == 1:
                        break
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break


class MqttNotifier(LicensePlateNotifier):
    TOPIC = "/vn/hust/iot/lp-detector"

    def __init__(self, broker, port, client_id):
        self.client = self.create_client(broker, client_id, port)

    @staticmethod
    def create_client(broker, client_id, port):
        def on_connect(this_client, userdata, flags, rc):
            if rc == 0:
                print("Connected to MQTT Broker!")
            else:
                print("Failed to connect, return code %d\n", rc)

        client = mqtt_client.Client(client_id)
        client.on_connect = on_connect
        # TODO add auto reconnect mechanism, see https://www.emqx.com/en/blog/how-to-use-mqtt-in-python
        client.connect(broker, port)
        return client

    def notify(self, license_plate: str):
        result = self.client.publish(topic=self.TOPIC, payload=license_plate)
        status = result[0]
        if status == 0:
            print(f"Send `{license_plate}` to topic `{self.TOPIC}`")
        else:
            print(f"Failed to send message to topic {self.TOPIC}")


class ConsoleNotifier(LicensePlateNotifier):
    def notify(self, licence_plate):
        print(f"Detected {licence_plate}")


def parse_arguments():
    arg_parser = ArgumentParser()
    arg_parser.add_argument('-v', '--video-src', default=0,
                            help="Source of video input. "
                                 "Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), "
                                 "or video file (vid.mp4)")
    arg_parser.add_argument('-mb', '--mqtt-broker', help="Hostname of MQTT broker",
                            default="broker.hivemq.com")
    arg_parser.add_argument('-mp', '--mqtt-port', help="Port of MQTT broker",
                            default=1883)
    arg_parser.add_argument('-mi', '--mqtt-client-id', help="MQTT client-ID",
                            default="iot-lp-detector")
    return arg_parser.parse_args()


if __name__ == '__main__':
    args = parse_arguments()
    mqtt_notifier = MqttNotifier(broker=args.mqtt_broker, port=args.mqtt_port, client_id=args.mqtt_client_id)
    console_notifier = ConsoleNotifier()
    detector = LicensePlateDetector(video_src=args.video_src, notifiers=[console_notifier, mqtt_notifier])
    detector.run()

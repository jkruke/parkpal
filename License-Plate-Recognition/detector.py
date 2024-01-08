import dataclasses
import json
import socket
import time
from argparse import ArgumentParser
from dataclasses import dataclass
from typing import List

import cv2
import imutils
import numpy as np
import torch
from confluent_kafka import Producer

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
    def __init__(self, record_entrance: bool, parking_lot_id: int, notifiers: List[LicensePlateNotifier]=[],
                 min_req_duplicates=0):
        self.parking_lot_id = parking_lot_id
        self.record_entrance = record_entrance
        self.notifiers = notifiers
        self.min_req_duplicates = min_req_duplicates
        self.known_lp = {}

    def add_all(self, license_plates: list):
        for lp in license_plates:
            # ignore license plate if it doesn't meet the formal requirements:
            if len(lp) not in [9, 10, 11]:
                continue

            if lp not in self.known_lp:
                dup_count = 0
                notified = False
            else:
                dup_count = self.known_lp[lp]['dup_count'] + 1
                notified = self.known_lp[lp]['notified']
            self.known_lp[lp] = {'lastseen': time.time(), 'dup_count': dup_count, 'notified': notified}

        return self.handle_known_lps()

    def handle_known_lps(self):
        detections = []
        known_lp_updated = self.known_lp.copy()
        for lp, info in self.known_lp.items():
            if info['notified'] is False and info['dup_count'] >= self.min_req_duplicates:
                # notify new record:
                for notifier in self.notifiers:
                    detection = LicensePlateDetection(parking_lot_id=self.parking_lot_id, license_plate=lp,
                                                      entrance=self.record_entrance)
                    detections.append(detection)
                    notifier.notify(detection)
                info['notified'] = True

            # cleanup old records:
            if info['lastseen'] < time.time() - 60:
                known_lp_updated.pop(lp)
            else:
                known_lp_updated[lp] = info

        self.known_lp = known_lp_updated
        return detections


# noinspection DuplicatedCode
class LicensePlateDetector:

    FPS = 2
    FRAME_DIFF_THRESHOLD = 0

    def __init__(self, video_src: int | str, record_entrance: bool, parking_lot_id: int,
                 notifiers: List[LicensePlateNotifier], show_video=False):
        """
        :param video_src specify the source of video input.
        Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), or video file (vid.mp4)
        """
        self.video_src = video_src
        self.notifiers = notifiers
        self.vid = None
        self.yolo_LP_detect = None
        self.yolo_license_plate = None
        self.load_model()
        self.lp_handler = LicensePlateHandler(record_entrance, parking_lot_id, notifiers, min_req_duplicates=self.FPS)
        self.show_video = show_video

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
            cv2.destroyAllWindows()

    def init_video(self):
        self.vid = cv2.VideoCapture(self.video_src)

    def process_video(self):
        latest_frame_ts = time.time()
        timeout = 1 / self.FPS
        latest_inspected_frame = None
        while True:
            frame = self.read_video_frame()
            if time.time() - latest_frame_ts < timeout:
                continue
            latest_frame_ts = time.time()

            if self.show_video:
                cv2.imshow('frame', frame)
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    return

            frame = imutils.resize(frame, width=250)
            frame = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            frame = np.dstack([frame, frame, frame])

            if latest_inspected_frame is None:
                latest_inspected_frame = frame

            # only apply frame differencing if desired:
            if self.FRAME_DIFF_THRESHOLD > 0:
                # calculate difference between current frame and the last inspected frame
                diff = cv2.absdiff(frame, latest_inspected_frame)
                movement_indicator = np.mean(diff)
                # ignore frame if the difference is not significant enough:
                if movement_indicator < self.FRAME_DIFF_THRESHOLD:
                    continue
                latest_inspected_frame = frame

            # inspect frame:
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

            retries += 1
            if retries > 10:
                print("Reinitializing video capture due to too many failed read operations.")
                self.vid.release()
                self.init_video()
                retries = 0


class KafkaNotifier(LicensePlateNotifier):
    TOPIC = "quickstart-events"

    def __init__(self, server):
        conf = {'bootstrap.servers': server,
                'client.id': socket.gethostname()}
        self.producer = Producer(conf)

    def notify(self, detection: LicensePlateDetection):
        payload = json.dumps(dataclasses.asdict(detection))
        print(f"[Kafka] Send to {self.TOPIC}: {payload}")
        self.producer.produce(self.TOPIC, payload)
        self.producer.flush()
        print("Flushed Kafka message")


class ConsoleNotifier(LicensePlateNotifier):
    def notify(self, detection):
        print(f"Detected {detection.license_plate}")


def parse_arguments():
    arg_parser = ArgumentParser()
    arg_parser.add_argument('-v', '--video-src', default=0,
                            help="Source of video input. "
                                 "Can be device number (0), RTSP video stream (rtsp://user:pwd@host.local:8081), "
                                 "or video file (vid.mp4)")
    arg_parser.add_argument('--show-video', help="Show video", default=False, action="store_true")
    arg_parser.add_argument('-ks', '--kafka-server', help="Kafka Server")
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
    print("Using notifiers:", notifiers)
    detector = LicensePlateDetector(video_src=args.video_src,
                                    record_entrance=args.record_entrance,
                                    parking_lot_id=args.parking_lot_id,
                                    notifiers=notifiers,
                                    show_video=args.show_video)
    detector.run()

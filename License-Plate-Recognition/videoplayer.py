import logging

import cv2


class VideoCam():
    def __init__(self, url):
        self.url = url
        self.cap = cv2.VideoCapture(self.url)
        self.get_frame()
        self.get_frame_read()
        logging.basicConfig(format='%(asctime)s %(message)s', level=logging.INFO)

    def check_camera(self, cap):
        logging.info('Camera {} status: {}'.format(self.url, cap.isOpened()))

    def show_frame(self, frame, name_fr='NAME'):
        cv2.imshow(name_fr, frame)
        # cv2.imshow(name_fr, cv2.resize(frame, (0, 0), fx=0.4, fy=0.4))
        cv2.waitKey(1)

    def get_frame(self):
        return self.cap.retrieve()

    def get_frame_read(self):
        return self.cap.read()

    def close_cam(self):
        self.cap.release()
        cv2.destroyAllWindows()

    def restart_capture(self, cap):
        cap.release()
        self.cap = cv2.VideoCapture(self.url)


if __name__ == '__main__':
    SKIPFRAME = 0
    v1 = VideoCam(url="rtmp://localhost/live/livestream")
    # v1 = VideoCam(url="rtsp://admin:admin@jk.local:8554/live")
    v1.check_camera(v1.cap)
    ct = 0
    while True:
        ct += 1
        try:
            ret = v1.cap.grab()
            if SKIPFRAME == 0 or ct % SKIPFRAME == 0:  # skip some frames
                ret, frame = v1.get_frame()
                if not ret:
                    v1.restart_capture(v1.cap)
                    v1.check_camera(v1.cap)
                    continue

                # frame HERE
                v1.show_frame(frame, 'frame')
        except KeyboardInterrupt:
            v1.close_cam()
            exit(0)

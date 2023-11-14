from detector import LicensePlateHandler, LicensePlateNotifier, LicensePlateDetection


class TestNotifier(LicensePlateNotifier):
    notified = []

    def notify(self, detection: LicensePlateDetection):
        self.notified.append(detection.license_plate)


def test_license_plate_handler():
    def notify_instant():
        notifier = TestNotifier()
        notifier.notified = []
        handler = LicensePlateHandler([notifier], min_req_duplicates=0)
        handler.add_all(['A1'])
        assert notifier.notified == ['A1']
        handler.add_all(['A1'])
        # A1 won't be notified twice:
        assert notifier.notified == ['A1']
        print("[PASSED] notify_instant")

    def notify_2_req_duplicates():
        notifier = TestNotifier()
        notifier.notified = []
        handler = LicensePlateHandler([notifier], min_req_duplicates=2)
        handler.add_all(['A1'])
        assert notifier.notified == []
        handler.add_all(['A1'])
        assert notifier.notified == []
        handler.add_all(['A1'])
        # A1 will be notified on the second duplicate:
        assert notifier.notified == ['A1']
        print("[PASSED] notify_2_req_duplicates")

    notify_instant()
    notify_2_req_duplicates()


if __name__ == '__main__':
    test_license_plate_handler()

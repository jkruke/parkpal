FROM python:3.11.6-slim

WORKDIR iot

# opencv dependencies:
RUN apt-get update && apt-get install python3-opencv -y

COPY requirement.txt requirements.txt
RUN pip install -r requirements.txt

COPY function/ function/
COPY model/ model/
COPY yolov5/ yolov5/
COPY detector.py .

CMD ["-h"]
ENTRYPOINT ["python", "detector.py"]
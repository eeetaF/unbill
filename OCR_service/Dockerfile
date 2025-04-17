FROM nvidia/cuda:12.2.0-runtime-ubuntu22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

RUN apt-get update && apt-get install -y \
    python3-pip python3-opencv libglib2.0-0 libsm6 libxext6 libxrender-dev git

RUN pip3 install easyocr flask opencv-python-headless torch torchvision torchaudio requests dotenv

WORKDIR /app
COPY app /app

CMD ["python3", "ocr_server.py"]

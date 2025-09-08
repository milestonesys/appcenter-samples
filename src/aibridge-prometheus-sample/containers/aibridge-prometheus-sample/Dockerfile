FROM python:3.12-slim

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app
COPY app/requirements.txt /app/requirements.txt
COPY app/.env /app/.env
COPY app/src/ /app/src

RUN rm -rf /root/.cache/pip && \
    pip3 install --upgrade pip && \
    pip3 install -r /app/requirements.txt --no-cache-dir && \
    rm -rf /root/.cache/pip

CMD [ "/bin/bash", "-c",  "uvicorn src.main:app --host 0.0.0.0 --port 9090" ]
#!/bin/bash

IMAGE=starland-backend:latest
CONTAINER_NAME=starland-backend
HTTP_PORT=8083
IMAGE_DIR=/data/starland/image
VOICE_DIR=/data/starland/voice
start() {
  docker run -d -p ${HTTP_PORT}:8083 -v ${PWD}/conf:/app/conf \
    -v ${PWD}/logfile:/app/logfile \
    -v ${IMAGE_DIR}:/app/image \
    -v ${VOICE_DIR}:/app/voice \
    --name ${CONTAINER_NAME} ${IMAGE}
}

stop() {
  docker rm ${CONTAINER_NAME} --force
}

case C"$1" in
C)
  echo "Usage: $0 {start|stop|restart}"
  ;;
Cstart)
  start
  echo "Start Done!"
  ;;
Cstop)
  stop
  echo "Stop Done!"
  ;;
Crestart)
  stop
  start
  echo "Restart Done!"
  ;;
C*)
  echo "Usage: $0 {start|stop|restart}"
  ;;
esac

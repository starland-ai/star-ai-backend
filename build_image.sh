#!/bin/bash

IMAGE=starland-backend:latest
docker rmi ${IMAGE}
docker build --label project=starland-backend -t ${IMAGE} .
docker push ${IMAGE}

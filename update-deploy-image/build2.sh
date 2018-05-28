#!/usr/bin/env bash

docker build  -t build-update-deploy-image2 . && \
docker run --rm -v/usr/local/bin:/out build-update-deploy-image2 cp /go/bin/update-deploy-image /out/update-deploy-image
#!/usr/bin/env bash

cd ..
docker build -f Dockerfile_ingress \
-t registry.cn-hangzhou.aliyuncs.com/mhc_dev/ingress-watcher:latest . \
&& docker push registry.cn-hangzhou.aliyuncs.com/mhc_dev/ingress-watcher:latest
#!/bin/bash

docker run -d --name=k8s-dev \
--restart=always \
-v`pwd`:/go/src/k8s-dev \
-v /tmp/common:/common \
-p 50012:8080 \
registry.cn-hangzhou.aliyuncs.com/mhc_dev/k8s-app-build-env:go-1.10  \
tail -f /dev/null
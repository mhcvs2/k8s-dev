#!/usr/bin/env bash

docker build -f ../Dockerfile_flex --build-arg http_proxy=http://109.105.4.17:3128 --build-arg https_proxy=http://109.105.4.17:3128 \
-t registry.bst-1.cns.bstjpc.com:5000/flex-provisioner:20180529 . \
&& docker push registry.bst-1.cns.bstjpc.com:5000/flex-provisioner:20180529
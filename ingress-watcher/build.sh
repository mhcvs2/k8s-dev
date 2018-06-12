#!/usr/bin/env bash

cd ..
docker build -f Dockerfile_ingress \
-t registry.bst-1.cns.bstjpc.com:5000/ingress-watcher:20180612 . \
&& docker push registry.bst-1.cns.bstjpc.com:5000/ingress-watcher:20180612
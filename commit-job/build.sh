#!/usr/bin/env bash

cd ..

docker build --build-arg http_proxy=http://109.105.4.17:8119 --build-arg https_proxy=http://109.105.4.17:8119 \
    -f Dockerfile_commit_job -t registry.bst-1.cns.bstjpc.com:5000/docker-commit-push:20180828 .

docker push registry.bst-1.cns.bstjpc.com:5000/docker-commit-push:20180828
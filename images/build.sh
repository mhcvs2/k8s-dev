#!/usr/bin/env bash

docker build -t registry.bst-1.cns.bstjpc.com:5000/k8s-dev:v1.10 .

docker push registry.bst-1.cns.bstjpc.com:5000/k8s-dev:v1.10
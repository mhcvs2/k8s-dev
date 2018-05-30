#!/usr/bin/env bash

cd ..

docker build -f Dockerfile_commit_job -t registry.bst-1.cns.bstjpc.com:5000/docker-commit-push:20180530 .

docker push registry.bst-1.cns.bstjpc.com:5000/docker-commit-push:20180530
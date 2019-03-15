#!/bin/bash

ds=(
k8s.io/client-go
k8s.io/apimachinery
github.com/google/gofuzz
k8s.io/klog
k8s.io/utils
k8s.io/code-generator/...
github.com/golang/glog
github.com/googleapis/gnostic
github.com/golang/protobuf
k8s.io/kube-openapi/pkg/util/proto
)


for i in ${ds[@]}; do
  echo "$i---------------------------------------------"
  go get $i
done


echo "cd $GOPATH/src/k8s.io"
cd $GOPATH/src/k8s.io
echo "git clone https://github.com/kubernetes/api.git"
git clone https://github.com/kubernetes/api.git

#!/bin/bash

ds=(
k8s.io/client-go/...
k8s.io/apimachinery/...
github.com/google/gofuzz
k8s.io/api
k8s.io/klog
k8s.io/utils
k8s.io/code-generator/...
github.com/golang/glog
)


for i in ${ds[@]}; do
  echo $i
  go get $i
done

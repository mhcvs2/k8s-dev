FROM registry.bst-1.cns.bstjpc.com:5000/k8s-client-go:v1.10

WORKDIR /go/src
COPY k8s.io/utils /go/src/k8s.io/utils
COPY k8s.io/kubernetes /go/src/k8s.io/kubernetes
COPY k8s.io/kube-openapi /go/src/k8s.io/kube-openapi

COPY src/github.com/golang/protobuf go/src/github.com/golang/protobuf
COPY src/github.com/googleapis go/src/github.com/googleapis
COPY src/gopkg.in     go/src/gopkg.in
FROM registry.bst-1.cns.bstjpc.com:5000/k8s-client-go:v1.10

WORKDIR /go/src

COPY src/github.com/* /go/src/github.com
COPY k8s.io/utils /go/src/k8s.io/utils
COPY k8s.io/kubernetes /go/src/k8s.io/kubernetes
COPY k8s.io/kube-openapi /go/src/k8s.io/kube-openapi
COPY src/gopkg.in /go/src/gopkg.in
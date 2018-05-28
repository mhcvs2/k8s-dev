FROM registry.bst-1.cns.bstjpc.com:5000/golang:1.10.2

RUN mkdir /go/src/update-deploy-image

WORKDIR /go/src/update-deploy-image

RUN go get -d  k8s.io/client-go/tools/clientcmd

RUN go get -d k8s.io/api/apps/v1beta1

COPY main.go .

RUN go install -v update-deploy-image
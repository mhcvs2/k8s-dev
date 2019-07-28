FROM alpine:3.7

RUN apk add --no-cache ca-certificates e2fsprogs findmnt

ADD do-csi-plugin /bin/

ENTRYPOINT ["/bin/do-csi-plugin"]

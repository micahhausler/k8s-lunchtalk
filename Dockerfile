# To build:
# $ docker run --rm -v $(pwd):/go/src/github.com/micahhausler/k8s-lunchtalk -w /go/src/github.com/micahhausler/k8s-lunchtalk golang:1.8-alpine go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'
# $ docker build -t micahhausler/k8s-lunchtalk .
#
# To run:
# $ docker run micahhausler/k8s-lunchtalk

FROM alpine

RUN apk -U add ca-certificates

MAINTAINER Micah Hausler, <hausler.m@gmail.com>

COPY k8s-lunchtalk /bin/k8s-lunchtalk
RUN chmod 755 /bin/k8s-lunchtalk

EXPOSE 3000

ENTRYPOINT ["/bin/k8s-lunchtalk"]

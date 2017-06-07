FROM golang:1.7-alpine
MAINTAINER W32Blaster

ADD . /go/src/fyzon

RUN set -ex && \
    apk add --no-cache git gcc g++ && \
    cd /go/src/fyzon && \
    go get -u -v github.com/gin-gonic/gin && \
    go get -u -v github.com/mattn/go-sqlite3 && \
    go get -u -v github.com/stretchr/testify && \
    go build && \
    go install .
   
ADD db/trans.sqlite3 /go/src/fyzon/db/trans.sqlite3
#ADD templates/* /templates/

ENV GIN_MODE=release


WORKDIR /go/src/fyzon

VOLUME /go/src/fyzon/db
EXPOSE 8080

ENTRYPOINT /go/bin/fyzon

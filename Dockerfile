FROM golang:1.7-alpine
MAINTAINER W32Blaster

ADD . /go/src/fyzon
ADD ./templates /go/src/fyzon/templates

    # install SQlite3 to set up a new database
RUN apk add --no-cache sqlite && \
    rm -f src/fyzon/db/trans.sqlite3 && \
    sqlite3 src/fyzon/db/trans.sqlite3 < src/fyzon/db/schema.sql

RUN set -ex && \
    apk add --no-cache git gcc g++ && \
    cd /go/src/fyzon && \
    go get -u -v github.com/gin-gonic/gin && \
    go get -u -v github.com/mattn/go-sqlite3 && \
    go get -u -v github.com/stretchr/testify && \
    go build && \
    go install . && \

    # remove sqlite and git, because we don't need it at runtime
    apk del sqlite git && \
    rm -rf /var/cache/apk/* && \

    # remove sources as well, because we already compiled them
    rm -rf src/*

ENV GIN_MODE=release

WORKDIR /go/src/fyzon

VOLUME /go/src/fyzon/db
EXPOSE 8080

ENTRYPOINT /go/bin/fyzon

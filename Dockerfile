FROM golang:1.7-alpine

ADD . /go/src/github.com/w32blaster/fyzon

# add templates to the WORKDIR
ADD ./templates /go/bin/templates && \
    ./assets/go/bin/assets

    # install SQlite3 to set up a new database
RUN apk add --no-cache sqlite && \
    rm -f src/github.com/w32blaster/fyzon/db/trans.sqlite3 && \
    mkdir -p /go/bin/db && \
    # install database to the WORKDIR
    sqlite3 /go/bin/db/trans.sqlite3 < src/github.com/w32blaster/fyzon/db/schema.sql && \
    cp src/github.com/w32blaster/fyzon/db/schema.sql /go/bin/ && \

    # copy DB import script for those folks who might want to keep database outside of the container
    cp src/github.com/w32blaster/fyzon/db/importDb.sh /go/bin/ && \
    chmod +x /go/bin/importDb.sh

RUN set -ex && \
    apk add --no-cache git gcc g++ && \
    cd /go/src/github.com/w32blaster/fyzon && \
    go get -u -v github.com/kardianos/govendor && \
    
    # install dependencies    
    govendor fetch -v +out  && \

    # build the project
    CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" && \
    go install . && \

    # remove git, because we don't need it at runtime
    apk del git && \
    rm -rf /var/cache/apk/* && \

    # remove sources as well, because we already compiled at the moment and we don't need them on runtime
    rm -rf /go/src

ENV GIN_MODE=release

WORKDIR /go/bin

VOLUME /go/bin/db
EXPOSE 8080

ENTRYPOINT /go/bin/fyzon

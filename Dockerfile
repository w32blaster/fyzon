FROM golang:1.7-alpine

ADD . /go/src/github.com/w32blaster/fyzon

# add templates to the WORKDIR
ADD ./templates /go/bin/templates
ADD ./assets /go/bin/assets

RUN echo '@edge http://nl.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories && \
    apk update && apk upgrade && \
    
    # install SQlite3 to set up a new database
    apk add --no-cache sqlite nodejs@edge nodejs-npm@edge && \
    mkdir -p /go/bin/db && \

    # install database to the WORKDIR
    sqlite3 /go/bin/db/trans.sqlite3 < src/github.com/w32blaster/fyzon/db/schema.sql && \
    cp src/github.com/w32blaster/fyzon/db/schema.sql /go/bin/ && \

    # copy DB import script for those folks who might want to keep database outside of the container
    cp src/github.com/w32blaster/fyzon/db/importDb.sh /go/bin/ && \
    chmod +x /go/bin/importDb.sh

RUN set -ex && \
    apk add --no-cache git gcc g++ && \
    npm -version 

RUN cd /go/src/github.com/w32blaster/fyzon && \
    go get -u -v github.com/kardianos/govendor && \
    
    # install dependencies    
    govendor fetch -v +out  && \

    # Build the Semantic UI
    npm install && \
    npm install semantic-ui --save && \
    npm install gulp -g && \

    # Build current theme for our Fyzon using gulp
    cd semantic && \
    gulp build && \
    cd .. && \

    # and copy freshly built Semantic files
    cp -r semantic /go/bin/ && \

    # build the project
    CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" && \
    go install . && \

    # remove git and nodeJS, because we don't need it at runtime
    apk del git nodejs nodejs-npm && \
    rm -rf /var/cache/apk/* && \

    # remove sources as well, because we already compiled at the moment and we don't need them on runtime
    rm -rf /go/src

ENV GIN_MODE=release

WORKDIR /go/bin

VOLUME /go/bin/db
EXPOSE 8080

ENTRYPOINT /go/bin/fyzon

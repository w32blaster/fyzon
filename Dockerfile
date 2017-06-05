FROM alpine:latest
MAINTAINER W32Blaster

ADD fyzon-release /fyzon
ADD db/schema.sql /db/schema.sql

ENV GIN_MODE=release

# Install SQLite3 to generate fresh database, then delete because we don't need it at runtime
RUN set -ex && \
    chmod +x /fyzon && \

    apk upgrade --update && \
    apk add --no-cache sqlite && \

    # Now create demo database from schema file
    sqlite3 /db/trans.sqlite3 < /db/schema.sql && \

    # ...and delete sqlite completely
    apk del sqlite && \
    rm -rf /var/cache/apk/*


VOLUME /db
EXPOSE 8080

ENTRYPOINT ["/fyzon"]

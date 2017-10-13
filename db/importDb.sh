#!/bin/sh

# just import the schema to the database.
# This script is supposed to be run inside the container
sqlite3 /go/bin/db/trans.sqlite3 < /go/bin/schema.sql
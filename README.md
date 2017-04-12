# Translation tool

Lightweight. Simple. Free.

1) import database

```
sqlite3 trans.sqlite3 < db/schema.sql
```

2) download and install all dependencies with one command

```
go get -u -v github.com/gin-gonic/gin
go get -u -v github.com/mattn/go-sqlite3

```

3) rebuild SemanticUI theme

```
npm install semantic-ui --save
cd semantic
gulp build
```

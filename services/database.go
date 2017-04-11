package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "bytes"
)

const (
    dbFile = "./trans.sqlite3"
)

func getProjects() (list string) {

    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("select id, name from projects")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var buffer bytes.Buffer
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        if err != nil {
            log.Fatal(err)
        }
        buffer.WriteString(name)
        buffer.WriteString(" ,")
    }

    return buffer.String();
}
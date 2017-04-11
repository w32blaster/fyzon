package services

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

const (
    dbFile = "./trans.sqlite3"
)

type Project struct {
    ID      int
    Name    string
}
type Projects []Project

/*
  Get all the projects as map
*/
func GetProjects() Projects {

    // connect to a database
    var db, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // make a request
    rows, err := db.Query("select id, name from projects")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var projects Projects
    for rows.Next() {
        var p Project
        err = rows.Scan(&p.ID, &p.Name)
        
        if err != nil {
            log.Fatal(err)
        }
        projects = append(projects, p)
    }

    return projects;
}
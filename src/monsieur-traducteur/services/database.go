package services

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

const (
    dbFile = "./trans.sqlite3"
)

type Term struct {
  ID int
  Code string
  Comment string
}

type Project struct {
    ID      int
    Name    string
    Terms   []Term
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

/*
  Get one project data
*/
func GetOneProject(id int) Project {

    // connect to a database
    var db, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // make a request
    rows, err := db.Query("select id, code, comment from terms where project_id = ?", id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var arrTerms []Term
    for rows.Next() {
        var t Term
        err = rows.Scan(&t.ID, &t.Code, &t.Comment)
        if err != nil {
            log.Fatal(err)
        }

        arrTerms = append(arrTerms, t)
    }

    p := Project{Terms: arrTerms}
    return p;
}

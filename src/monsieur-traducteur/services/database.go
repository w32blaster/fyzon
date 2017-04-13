package services

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

const (
    dbFile = "./trans.sqlite3"
)

type Translation struct {
  ID int
  Translation string
  LanguageCode string
  IsDefault bool
}

type Term struct {
  ID int
  Code string
  Comment string
  Translations []Translation
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
   Get term
*/
func GetOneTerm(id int) Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // make a request for one term
  rows, err := db.Query("select id, code, comment from terms where id = ? limit 1", id)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var t Term
  rows.Next()
  rows.Scan(&t.ID, &t.Code, &t.Comment)

  return t
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

    // make a request for all the terms
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

    // Make a request for a project info
    stmt, err := db.Query("select id, name from projects where id = ? limit 1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    stmt.Next()
    _ = stmt.Scan(&p.ID, &p.Name)

    return p;
}

/*
  Find Term with all the translations
*/
func GetTerm(termId int) Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // find one Term
  stmt, err := db.Query("select id, code, comment from terms where id = ? limit 1", termId)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  stmt.Next()
  var t Term
  _ = stmt.Scan(&t.ID, &t.Code, &t.Comment)

  // find all the translations
  rows, err := db.Query("select id, translation, language_code, is_default from translations where term_id = ?", termId)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var translations []Translation
  for rows.Next() {
    var tr Translation
    err = rows.Scan(&tr.ID, &tr.Translation, &tr.LanguageCode, &tr.IsDefault)

    if err != nil {
      log.Fatal(err)
    }

    translations = append(translations, tr)
  }

  t.Translations = translations
  return t
}

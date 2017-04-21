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
  CountryCode string
  IsDefault bool
}

type Term struct {
  ID int
  Code string
  Comment string
  Translations []Translation
  ProjectId int
}

type Project struct {
    ID      int
    Name    string
    Terms   []Term
    TermsCount int
    CountryCodes []string
}
type Projects []Project

/*
  Get all the projects as map
*/
func GetProjects() *Projects {

    // connect to a database
    var db, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // collect map project_id <=> list of language codes
    lRows, err := db.Query("SELECT country_code, project_id FROM project_languages")
    if err != nil {
        log.Fatal(err)
    }
    defer lRows.Close()

    mapLangs := make(map[int][]string)
    for lRows.Next() {
        var project_id int
        var country_code string
        err = lRows.Scan(&country_code, &project_id)

        if err != nil {
            log.Fatal(err)
        }

        mapLangs[project_id] = append(mapLangs[project_id], country_code)
    }

    // make a request
    rows, err := db.Query("SELECT p.id, p.name, COUNT(p.id) as cnt FROM projects AS p INNER JOIN terms AS t ON t.project_id = p.id GROUP BY p.id")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var projects Projects
    for rows.Next() {
        var p Project
        err = rows.Scan(&p.ID, &p.Name, &p.TermsCount)

        if err != nil {
            log.Fatal(err)
        }
        p.CountryCodes = mapLangs[p.ID]
        projects = append(projects, p)
    }

    return &projects;
}

/*
  Get one project data
*/
func GetOneProject(id int) *Project {

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
    p := Project{Terms: arrTerms, TermsCount: len(arrTerms), ID: id}

    // Make a request for a project info
    stmt, err := db.Query("select id, name from projects where id = ? limit 1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    stmt.Next()
    _ = stmt.Scan(&p.ID, &p.Name)

    return &p;
}

/*
  Find Term with all the translations
*/
func GetTerm(termId int) *Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // find one Term
  stmt, err := db.Query("select id, code, comment, project_id from terms where id = ? limit 1", termId)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  stmt.Next()
  var t Term
  _ = stmt.Scan(&t.ID, &t.Code, &t.Comment, &t.ProjectId)

  // find all the translations
  rows, err := db.Query("select id, translation, country_code, is_default from translations where term_id = ?", termId)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var translations []Translation
  existingLangs := make(map[string]bool) // store which languages we already have

  for rows.Next() {
    var tr Translation
    err = rows.Scan(&tr.ID, &tr.Translation, &tr.CountryCode, &tr.IsDefault)

    if err != nil {
      log.Fatal(err)
    }

    translations = append(translations, tr)
    existingLangs[tr.CountryCode] = true
  }

  // check if there are some languages missing, then add empty field
  langs := getAvailableLanguagesForProject(&t.ProjectId, db)
  for _, lang := range *langs {
      if (!existingLangs[lang]) {
        translations = append(translations, Translation{ID: -1, IsDefault: false, CountryCode: lang} )
      }
  }

  t.Translations = translations
  return &t
}


/**
 * Update one translation
 */
func UpdateTranslation(value string, termId int, countryCode string) {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // firstly, check whether is already exists (create or update?)
  row, _ := db.Query("SELECT count(*) FROM translations WHERE term_id=? AND country_code=?", termId, countryCode)
  var count int
	for row.Next() {
		row.Scan(&count)
  }

  if (count > 0) {

    // Update it
    _, err = db.Exec("UPDATE translations SET translation=? where term_id=?", value, termId)
  	if err != nil {
  		log.Fatal("Failed to update record:", err)
    }
  } else {

    // Create new translation
    _, err = db.Exec(`INSERT INTO translations(translation, country_code, is_default, term_id) VALUES (?, ?, ?, ?)`, value, countryCode, false, termId)
  	if err != nil {
  		log.Fatal("Failed to update record:", err)
    }
  }
}

/**
 * Add new language to given project
 */
func AddNewLanguage(projectId *int, countryCode *string) {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  if (!isLanguageAlreadyExists(projectId, countryCode, db)) {

    _, err = db.Exec("INSERT INTO project_languages(project_id,country_code) values(?, ?)", projectId, countryCode)
    if err != nil {
      log.Fatal("Failed to update record:", err)
    }

  } else {
    log.Print("This project has already this language");
  }
}

/**
 * Add new Term to the given project
 */
func AddNewTerm(termKey *string, termDescr *string, projectId *int) *Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  res, err := db.Exec("INSERT INTO terms(code,comment,project_id) values(?, ?, ?)", termKey, termDescr, projectId)
  if err != nil {
    log.Fatal("Failed to update record:", err)
  }

  addedId64,_ := res.LastInsertId()
  var addedId int
  addedId = int(addedId64)
  return GetTerm(addedId)
}

/**
 * Get list of available languages for a given project
 */
func getAvailableLanguagesForProject(projectId *int, db *sql.DB) *[]string {

  rows, err := db.Query("SELECT country_code FROM project_languages WHERE project_id=?", projectId)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var langs []string
  for rows.Next() {
    var lang string
    err = rows.Scan(&lang)

    if err != nil {
      log.Fatal(err)
    }

    langs = append(langs, lang)
  }

  return &langs
}

/**
 * Check whether the given language already exists for the given project
 */
func isLanguageAlreadyExists(projectId *int, countryCode *string, db *sql.DB) bool {

  row, _ := db.Query("SELECT count(*) FROM project_languages WHERE project_id=? AND country_code=?", projectId, countryCode)
  var count int
  for row.Next() {
    row.Scan(&count)
  }

  return count > 0
}

package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "strings"
    "bytes"
)

/**
 * Represents a language that project may expect
 */
type ProjectLanguage struct {
  CountryCode string
  IsDefault bool
}

/**
 * One translation in one given language for a selected term
 */
type Translation struct {
  ID int
  Translation string
  CountryCode string
  IsDefault bool
  TermId int
}

/**
 * Term (a key, that can have many translations on other languages)
 */
type Term struct {
  ID int
  Code string
  Comment string
  Translations []Translation
  ProjectId int
  HasDefault bool // whether default language has Translation for this term or not
}

type Project struct {
    ID      int
    Name    string
    Terms   []Term
    TermsCount int
    CountryCodes []string
    DefaultCountryCode string
}
type Projects []Project

/*
 *  Tuple pair, used via a file parsing
 */
type ImportTranslation struct {
  Translation string
  Comment string
}
/**
 * Creates new project
 */
func CreateNewProject(dbFilePath string, name string, defaultLanguage string) *Project {

  var db, err = sql.Open("sqlite3", dbFilePath)
  checkErr(err)
  defer db.Close()

  stmt, err := db.Prepare("INSERT INTO projects(name, default_country_code) values(?, ?)")
  checkErr(err)
  defer stmt.Close()

  res, err := stmt.Exec(name, defaultLanguage)
  checkErr(err)

  id64, err := res.LastInsertId()
  id := int(id64)
  checkErr(err)

  AddNewLanguage(dbFilePath, id, defaultLanguage)

  return FindOneProject(int(id), "")
}

/*
  Get all the projects as map
*/
func GetProjects(dbFilePath string) *Projects {

    // connect to a database
    var db, err = sql.Open("sqlite3", dbFilePath)
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
 * Get one project data
 *
 * @id - project id
 * @countryCode - language terms that doesn't have any translations yet. If empty (""), then
 *                show all the terms.
 */
func FindOneProject(id int, countryCode string) *Project {

    // connect to a database
    var db, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    p := Project{ID: id}

    // Make a request for a project info
    stmt, err := db.Query("select id, name, default_country_code from projects where id = ? limit 1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    stmt.Next()
    _ = stmt.Scan(&p.ID, &p.Name, &p.DefaultCountryCode)

    // make a request for all the terms
    var rows *sql.Rows
    if(countryCode == "") {

      // if no untranslated lang specified, return all
      sqlQuery := "select t.id, t.code, t.comment, ifnull(GROUP_CONCAT(tr.country_code), '') AS codes from terms AS t " +
              "LEFT JOIN translations AS tr ON tr.term_id = t.id " +
              "where t.project_id = ? GROUP BY t.code ORDER BY t.code"
      rows, _ = db.Query(sqlQuery, id);

    } else {

      // if untranslated lang is set, show only these terms
      sqlQuery := "select t.id, t.code, t.comment, ifnull(GROUP_CONCAT(tr.country_code), '') AS codes " +
              "FROM terms AS t " +
              "INNER JOIN project_languages AS pl ON pl.project_id = t.project_id " +
              "LEFT JOIN translations AS tr ON tr.term_id = t.id AND pl.country_code = tr.country_code " +
              "WHERE t.project_id = ? AND tr.id IS NULL AND pl.country_code = ? GROUP BY t.id ORDER BY code"
      rows, _ = db.Query(sqlQuery, id, countryCode);
    }

    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var arrTerms []Term
    for rows.Next() {
        var t Term
        var codes string
        err = rows.Scan(&t.ID, &t.Code, &t.Comment, &codes)
        if err != nil {
            log.Fatal(err)
        }
        t.HasDefault = isContainingDefaultLanguage(codes, p.DefaultCountryCode)

        arrTerms = append(arrTerms, t)
    }

    // remember how many terms we have in this project
    p.Terms = arrTerms
    p.TermsCount = len(arrTerms)

    // get all available languages for this project
    projectLanguages := getAvailableLanguagesForProject(id, p.DefaultCountryCode, db)
    p.CountryCodes = *asStringArray(projectLanguages)

    return &p;
}

/*
  Find Term with all the translations
*/
func GetTerm(dbFilePath string, termId int) *Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFilePath)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // find one Term
  stmt, err := db.Query("select t.id, t.code, t.comment, t.project_id, p.default_country_code from terms AS t " +
                        "INNER JOIN projects AS p ON p.id = t.project_id " +
                        "where t.id = ? GROUP BY t.project_id limit 1", termId)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  stmt.Next()
  var t Term
  var default_country_code string
  _ = stmt.Scan(&t.ID, &t.Code, &t.Comment, &t.ProjectId, &default_country_code)

  // find all the translations
  rows, err := db.Query("SELECT t.id, t.translation, t.country_code FROM translations AS t " +
    "INNER JOIN project_languages AS pl ON pl.country_code = t.country_code " +
    "WHERE t.term_id = ? GROUP BY t.country_code", termId)

  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var translations []Translation
  existingLangs := make(map[string]bool) // store which languages we already have

  for rows.Next() {
    tr := Translation{TermId: termId}
    err = rows.Scan(&tr.ID, &tr.Translation, &tr.CountryCode)

    if err != nil {
      log.Fatal(err)
    }

    tr.IsDefault = (default_country_code == tr.CountryCode)

    translations = append(translations, tr)
    existingLangs[tr.CountryCode] = true
  }

  // check if there are some languages missing, then add empty field
  langs := getAvailableLanguagesForProject(t.ProjectId, default_country_code, db)
  for _, lang := range *langs {
      if (!existingLangs[lang.CountryCode]) {
        translations = append(translations, Translation{ID: -1, IsDefault: lang.IsDefault, CountryCode: lang.CountryCode, TermId: termId} )
      }
  }

  t.Translations = translations
  return &t
}


/**
 * Update one translation
 */
func UpdateTranslation(dbFilePath string, value string, termId int, countryCode string) {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFilePath)
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
    _, err = db.Exec("UPDATE translations SET translation=? WHERE term_id=? AND country_code=?", value, termId, countryCode)
  	if err != nil {
  		log.Fatal("Failed to update record:", err)
    }
  } else {

    // Create new translation
    _, err = db.Exec("INSERT INTO translations(translation, country_code, term_id) VALUES (?, ?, ?)", value, countryCode, termId)
  	if err != nil {
  		log.Fatal("Failed to update record:", err)
    }
  }
}

/**
 * Add new language to given project
 */
func AddNewLanguage(dbFilePath string, projectId int, countryCode string) {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFilePath)
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
func AddNewTerm(dbFilePath string, termKey string, termDescr string, projectId int) *Term {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFilePath)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  res, err := db.Exec("INSERT INTO terms(code,comment,project_id) values(?, ?, ?)", termKey, termDescr, projectId)
  if err != nil {
    log.Fatal("Failed to update record:", err)
  }

  addedId64,_ := res.LastInsertId()
  return GetTerm(dbFilePath, int(addedId64))
}

/**
 * Delete one term and its translations
 */
func DeleteTerm(dbFilePath string, termId int) bool {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFilePath)
  if err != nil {
      log.Fatal(err)
      return false
  }
  defer db.Close()

  // Delete translations
  _, err = db.Exec("DELETE FROM translations WHERE term_id = ?", termId)
	if err != nil {
		log.Fatal(err)
    return false
  }

  // and Delete term itself
  _, err = db.Exec("DELETE FROM terms WHERE id = ?", termId)
  if err != nil {
    log.Fatal(err)
    return false
  }

  return true
}

/**
 * Recursively delete the project
 */
func DeleteProject(dbFilePath string, projectId int) bool {

  db, err := sql.Open("sqlite3", dbFilePath)
	checkErr(err)
  defer db.Close()

	// Turning on Forgein key support (for cascading deleting)
	stmt, err := db.Prepare("PRAGMA foreign_keys = ON;")
	checkErr(err)

	_, err = stmt.Exec()
  checkErr(err)

  // we expect cascade deleting of terms, project_languages and transations, please refer to constrains it the schema.sql
  stmt, err = db.Prepare("DELETE FROM projects WHERE id = ?")
	checkErr(err)

	res, err := stmt.Exec(projectId)
	checkErr(err)

	affect, err := res.RowsAffected()
  checkErr(err)

  return affect > 0
}

/**
 * Get list of available languages for a given project
 */
func getAvailableLanguagesForProject(projectId int, project_default_lang string, db *sql.DB) *[]ProjectLanguage {

  rows, err := db.Query("SELECT country_code FROM project_languages WHERE project_id=?", projectId)
  if err != nil {
      log.Fatal(err)
  }
  defer rows.Close()

  var langs []ProjectLanguage
  for rows.Next() {
    var lang ProjectLanguage
    err = rows.Scan(&lang.CountryCode)
    lang.IsDefault = (lang.CountryCode == project_default_lang)

    if err != nil {
      log.Fatal(err)
    }

    langs = append(langs, lang)
  }

  return &langs
}

/**
 * Save the imported terms within a transaction.
 */
func SaveImportedTermsForProject(dbFilePath string, terms map[string]ImportTranslation, countryCode string, projectId int) error {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  project := _getPorjectById(projectId, db)

  // insert languages to the project, if is not used still
  _insertLanguage(dbFilePath, countryCode, project, db)

  _insertAllTerms(terms, projectId, db)

  _insertAllTranslations(terms, projectId, countryCode, db)

  return nil
}

/**
 * Insert country code if not used
 */
func _insertLanguage(dbFilePath string, countryCode string, project *Project, db *sql.DB) {

  existingLangs := getAvailableLanguagesForProject(project.ID, project.DefaultCountryCode, db)
  isFound := false
  for _, lang := range *existingLangs {
      if (!isFound && lang.CountryCode == countryCode) {
        isFound = true
        break
      }
  }

  if (!isFound) {
      AddNewLanguage(dbFilePath, project.ID, countryCode)
  }

}

func _insertAllTerms(terms map[string]ImportTranslation, projectId int, db *sql.DB) error {

  // begin transaction
  tx, err := db.Begin()
    if err != nil {
      log.Fatal(err)
  }

  // prepare one common statement
  stmt, err := tx.Prepare("INSERT INTO terms(project_id, code, comment) values(?, ?, ?)")
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer stmt.Close()

  for key, value := range terms {

    if termId := getTermIdFor(projectId, key, db); termId == -1 {

      // insert and get fresh term
      _, err = stmt.Exec(projectId, key, value.Comment)
      if err != nil {
        log.Fatal(err)
      }

    }
  }

  // Commit transaction
  tx.Commit()
  return nil
}

/**
 * Insert all the translations within one transaction
 */
func _insertAllTranslations(terms map[string]ImportTranslation, projectId int, countryCode string, db *sql.DB) error {
  tx, err := db.Begin()
  	if err != nil {
  		log.Fatal(err)
  }

  stmtTranslation, err := tx.Prepare("insert into translations(translation, country_code, term_id) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
    return err
	}
	defer stmtTranslation.Close()

  // now iterate over all the items and insert all of them
  for key, value := range terms {

    // insert and get fresh term
    if termId := getTermIdFor(projectId, key, db); termId != -1 {

      // add translation for this term:
      _, err = stmtTranslation.Exec(value.Translation, countryCode, termId)
  		if err != nil {
  			log.Fatal(err)
  		}
    }

	}

  // Commit transaction
  tx.Commit()
  return nil
}

/**
 * Check whether the given language already exists for the given project
 */
func isLanguageAlreadyExists(projectId int, countryCode string, db *sql.DB) bool {

  row, _ := db.Query("SELECT count(*) FROM project_languages WHERE project_id=? AND country_code=?", projectId, countryCode)
  var count int
  for row.Next() {
    row.Scan(&count)
  }

  return count > 0
}

/**
 * Check whether this term is already exists in the given project
 */
func getTermIdFor(projectId int, termCode string, db *sql.DB) int {

  row, _ := db.Query("SELECT id FROM terms WHERE project_id=? AND code=?", projectId, termCode)
  termId := -1
  for row.Next() {
    row.Scan(&termId)
  }

  return termId
}

/**
 * Generate the content of a file with translations
 */
func GenerateFile(projectId int, countryCode string, delimeter string) (string, error) {

  // connect to a database
  var db, err = sql.Open("sqlite3", dbFile)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  // get all the translations
  lRows, err := db.Query("SELECT t.code,tr.translation,t.comment FROM terms AS t INNER JOIN translations AS tr ON tr.term_id = t.id WHERE t.project_id = ? AND tr.country_code = ?", projectId, countryCode)
  if err != nil {
      log.Fatal(err)
      return "", err
  }
  defer lRows.Close()

  var buffer bytes.Buffer

  for lRows.Next() {
      var code string
      var translation string
      var comment string

      err = lRows.Scan(&code, &translation, &comment)

      if err != nil {
          log.Fatal(err)
          return "", err
      }
      if (len(comment) > 0) {
          buffer.WriteString("# " + comment + "\n")
      }

      buffer.WriteString(code + delimeter + translation + "\n")
  }

  return buffer.String(), nil
}

/**
 * Simply turns array of ProjectLanguage to array of Strings, having
 */
func asStringArray(langs *[]ProjectLanguage) *[]string {
  var arrCountryCodes = make([]string, len(*langs))
  for i := range *langs {
      arrCountryCodes[i] = (*langs)[i].CountryCode
  }
  return &arrCountryCodes
}

func _getPorjectById(projectId int, db *sql.DB) *Project {
  stmt, err := db.Query("select id, name, default_country_code from projects where id = ? limit 1", projectId)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  var p Project
  stmt.Next()
  _ = stmt.Scan(&p.ID, &p.Name, &p.DefaultCountryCode)
  return &p
}

/**
* Takes the list of codes separated by comma and searches for the default code
*/
func isContainingDefaultLanguage(codes string, defaultCode string) bool {
  for _, code := range strings.Split(codes, ",") {
        if code == defaultCode {
            return true
        }
    }
  return false
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

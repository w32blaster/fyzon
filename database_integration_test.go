package main

import (
  "database/sql"
  "testing"
	"github.com/stretchr/testify/assert"
  "os"
	"log"
  "io/ioutil"
)

const (
    dbTestFile = "/tmp/testDb.sqlite3"
)

func check(e error) {
    if e != nil {
        os.Remove(dbTestFile)
        panic(e)
    }
}

/**
 * After and Before.
 * Set up the database for testing
 */
func TestMain(m *testing.M) {

	log.Println("Set up the database Before all the tests")

  // read SQL script with schema
  dat, err := ioutil.ReadFile("./db/schema.sql")
  check(err)
  tableDDL := string(dat)

  db := _connectDb()
  defer db.Close()

  // create tables and insert some data into it
  _, err3 := db.Exec(tableDDL)
	check(err3)

  // assert there are nine translation (see schema.sql for details)
  count := _getCountOf(db, "translations")
  if (count != 9) {
    panic("Database doesn't contain exactly 9 translations as expected from file schema.sql")
  }

  /*
   run all the tests here:
  */
	exitVal := m.Run()

  // clean up (remove the database) and exit
  os.Remove(dbTestFile)
	os.Exit(exitVal)
}

/**
 * We create new item and verify it is returned
 * and the database contains new item
 */
func TestWeAddNewItem(t *testing.T) {
  log.Println("test WeAddNewItem")

  // Given:
  expectedId := 7 // because we have 6 testing terms in the database and next ID is 7th

  // When:
  createdTerm := AddNewTerm(dbTestFile, "test.term.key", "Some description", 1)

  // Then:
  assert.NotNil(t, createdTerm)

  // and:
  assert.Equal(t, "test.term.key", createdTerm.Code)
  assert.Equal(t, "Some description", createdTerm.Comment)
  assert.Equal(t, 1, createdTerm.ProjectId)
  assert.Equal(t, expectedId, createdTerm.ID)

  // after:
  _ = DeleteTerm(dbTestFile, createdTerm.ID)

}

/**
 * Insert a term and few translations. Delete it. Make sure that the term
 * and its translations are gone
 */
func TestDeleteOneTerm(t *testing.T) {
  log.Println("test DeleteOneTerm")

  // before:
  db := _connectDb()
  defer db.Close()

  // Given:
  createdTerm := AddNewTerm(dbTestFile, "test.term.to.be.deleted", "The term that will be deleted", 1)

  // and:
  UpdateTranslation(dbTestFile, "Translation in English", createdTerm.ID, "gb")
  UpdateTranslation(dbTestFile, "Translation in Russian", createdTerm.ID, "ru")
  UpdateTranslation(dbTestFile, "Translation in Greece", createdTerm.ID, "gr")

  // Before start, assert there are 12 translation (see schema.sql for details)
  translationsCount := _getCountOf(db, "translations")
  assert.Equal(t, 12, translationsCount, "After we insert more three translations, we expect 12 (9+3) of them altogether")

  // ...and 7 terms (6 from schema.sql and one added above)
  termsCount := _getCountOf(db, "terms")
  assert.Equal(t, 7, termsCount, "After we insert more one term, we expect 7 terms")

  // When:
  result := DeleteTerm(dbTestFile, createdTerm.ID)

  // Then:
  assert.True(t, result)

  // and:
  translationsCount = _getCountOf(db, "translations")
  assert.Equal(t, translationsCount, 9, "After we deleted one term, we expect again 9 (12-3) translations in the database")

  // and:
  termsCount = _getCountOf(db, "terms")
  assert.Equal(t, termsCount, 6, "After we deleted one term, we expect to have 6 terms again")
}

/**
 * Delete project
 *
 */
func TestDeleteProject(t *testing.T) {
  log.Println("test DeleteProject")

  /*
     Populate the database
  */

  // before:
  db := _connectDb()
  defer db.Close()

  // Given:
  createdProject := CreateNewProject(dbTestFile, "Test Project", "en")
  AddNewLanguage(dbTestFile, createdProject.ID, "ru")
  AddNewLanguage(dbTestFile, createdProject.ID, "de")

  // and:
  createdTermOne := AddNewTerm(dbTestFile, "term.one.project.to.be.deleted", "", createdProject.ID)
  createdTermTwo := AddNewTerm(dbTestFile, "term.two.project.to.be.deleted", "", createdProject.ID)

  // and:
  UpdateTranslation(dbTestFile, "Translation in English", createdTermOne.ID, "en")
  UpdateTranslation(dbTestFile, "Translation in Russian", createdTermOne.ID, "ru")
  UpdateTranslation(dbTestFile, "Translation in German", createdTermOne.ID, "de")

  UpdateTranslation(dbTestFile, "Translation2 in English", createdTermTwo.ID, "en")
  UpdateTranslation(dbTestFile, "Translation2 in Russian", createdTermTwo.ID, "ru")

  // assert, that we have three projects in the database (two from schema.sql and one created above)
  projectsCount := _getCountOf(db, "projects")
  assert.Equal(t, projectsCount, 3, "We have 3 projects after we created a project")

  // assert we have 5 (initially) + 1 (default language when we created new project) + 2 (in addition) = 8 languages
  langsCount := _getCountOf(db, "project_languages")
  assert.Equal(t, langsCount, 8, "We have 8 languages after we created a project")

  // assert there are 14 translation (9 (schema.sql) + 5 added above)
  translationsCount := _getCountOf(db, "translations")
  assert.Equal(t, 14, translationsCount, "After we insert more three translations, we expect 14 (9+5) of them altogether")

  /*
     Here test begins
  */

  // When:
  result := DeleteProject(dbTestFile, createdProject.ID)

  // Then:
  assert.True(t, result);

  // and:
  _assertDatabaseInOriginalCondition(t, db)
}

/**
 * Get count of records in the given table
 */
func _getCountOf(db *sql.DB, tableName string) int {
  row, _ := db.Query("SELECT count(*) FROM " + tableName)
  var count int
  for row.Next() {
    row.Scan(&count)
  }
  return count
}

func _connectDb() *sql.DB {
  db, err2 := sql.Open("sqlite3", dbTestFile)
  check(err2)
  return db
}

/**
 * Helping method to check that the database in its original state
 * and all the added items are gone.
*/
func _assertDatabaseInOriginalCondition(t *testing.T, db *sql.DB)  {

  projectsCount := _getCountOf(db, "projects")
  assert.Equal(t, projectsCount, 2, "We have 2 projects as before test started (as in the schema.sql)")

  // and
  langsCount := _getCountOf(db, "project_languages")
  assert.Equal(t, langsCount, 5, "We have 5 languages as before test started (as in the schema.sql)")

  // and
  translationsCount := _getCountOf(db, "translations")
  assert.Equal(t, 9, translationsCount, "We expect 9 of them altogether")

}

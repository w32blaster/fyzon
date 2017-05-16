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

  db, err2 := sql.Open("sqlite3", dbTestFile)
  check(err2)
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

  // before:
  db, err2 := sql.Open("sqlite3", dbTestFile)
  check(err2)
  defer db.Close()

  // Given:
  createdTerm := AddNewTerm(dbTestFile, "test.term.to.be.deleted", "The term that will be deleted", 1)

  // and:
  UpdateTranslation(dbTestFile, "Translation in English", createdTerm.ID, "gb")
  UpdateTranslation(dbTestFile, "Translation in Russian", createdTerm.ID, "ru")
  UpdateTranslation(dbTestFile, "Translation in Greece", createdTerm.ID, "gr")

  // Before start, assert there are nine translation (see schema.sql for details)
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

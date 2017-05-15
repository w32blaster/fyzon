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
    dbTestFile = "./testDb.sqlite3"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

/**
 * After and Before. Set up the database to test
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
  row, _ := db.Query("SELECT count(*) FROM translations")
  var count int
  for row.Next() {
    row.Scan(&count)
  }

  // assert that we inserted demo data to DB
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
  code := "test.term.key"
  descr := "Some description"
  projectId := 1

  // And:
  expectedId := 7 // because we have 6 testing terms in the database and next ID is 7th

  // When:
  createdTerm := AddNewTerm(dbTestFile, &code, &descr, &projectId)

  // Then:
  assert.NotNil(t, createdTerm)

  // and:
  assert.Equal(t, code, createdTerm.Code)
  assert.Equal(t, descr, createdTerm.Comment)
  assert.Equal(t, projectId, createdTerm.ProjectId)
  assert.Equal(t, expectedId, createdTerm.ID)
}

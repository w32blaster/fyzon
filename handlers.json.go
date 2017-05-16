package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"io"
	"log"
  "os"
	"bufio"
	"strings"
	"errors"
)

/**
 * Ping
 */
func Ping(c *gin.Context) {
	c.JSON(200, gin.H {
		"message": "pong",
	})
}

/**
 * Update one term
 */
func PostOneTerm(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	countryCode := c.Param("lang")
	value := c.PostForm("value")

	UpdateTranslation(dbFile, value, id, countryCode)

	c.JSON(http.StatusOK, gin.H {})
}

/**
 * Add new language
 */
func PostNewLanguage(c *gin.Context) {

		if projectId, err := strconv.Atoi(c.PostForm("projectId")); err == nil {

			countryCode := c.PostForm("countryCode")
	    AddNewLanguage(projectId, countryCode)
			c.JSON(http.StatusOK, gin.H {})

	} else {
      c.AbortWithStatus(http.StatusBadRequest)
	}
}

/**
 * Add new term to project
 */
func PostNewTerm(c *gin.Context) {

    // key and projectId are mandatory
    termKey := c.PostForm("termKey")
		projectId, err := strconv.Atoi(c.PostForm("projectId"));

		if len(termKey) == 0 || err != nil {
      c.AbortWithStatus(http.StatusBadRequest)
		} else {
			// description is optional
			termDescr := c.PostForm("termDescr")
			addedTerm := AddNewTerm(dbFile, termKey, termDescr, projectId)
			c.JSON(http.StatusOK, gin.H {
				"term": addedTerm,
			})
		}
}

/**
 * Upload a new file with existing translations, parse it
 * and import all the items to a database
 */
func PostNewFile(c *gin.Context) {

    file, header, err := c.Request.FormFile("upload")
		delimeter := c.PostForm("delimeter")
		country := c.PostForm("country")
		projectId, err :=  strconv.Atoi(c.Param("id"))


    filename := header.Filename

    out, err := os.Create("/tmp/" + filename)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        log.Fatal(err)
    }

		mapLines, _ := readLines("/tmp/" + filename, delimeter)

		// save it somehow
		SaveImportedTermsForProject(*mapLines, country, projectId)

		// delete temp file
		_ = os.Remove("/tmp/" + filename)

		// redirect
		c.Redirect(http.StatusMovedPermanently, "/project/" + strconv.Itoa(projectId) + "?imported")
}

/**
 * Delete one term and all its translations
 */
func DeleteOneTerm(c *gin.Context)  {
		if termId, err := strconv.Atoi(c.Param("id")); err == nil {

			result := DeleteTerm(dbFile, termId)
			if (result) {
					c.JSON(http.StatusOK, gin.H {})
			}else {
				c.JSON(http.StatusBadRequest, gin.H {})
			}

		} else {
				c.AbortWithStatus(http.StatusBadRequest)
		}

}


/**
 * read lines for the given file
 */
func readLines(path string, delimeter string) (*map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
	  return nil, err
	}
	defer file.Close()

  mapLines := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if key, value, err := parseLine(scanner.Text(), delimeter); err == nil {
			mapLines[key] = value
		}
	}
	return &mapLines, scanner.Err()
}

/**
 * Parse one line and return pair of values
 */
func parseLine(line string, delimeter string) (string, string, error) {
  parts := strings.Split(line, delimeter)
	if (len(parts) == 2) {
		  return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
	} else {
			return "", "", errors.New("Didn't found two parts in the incoming string")
	}
}

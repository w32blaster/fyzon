package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/w32blaster/fyzon/generator"
)

/**
 * Ping
 */
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
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

	c.JSON(http.StatusOK, gin.H{})
}

/**
 * Add new language
 */
func PostNewLanguage(c *gin.Context) {

	if projectId, err := strconv.Atoi(c.PostForm("projectId")); err == nil {

		countryCode := c.PostForm("countryCode")
		AddNewLanguage(dbFile, projectId, countryCode)
		c.JSON(http.StatusOK, gin.H{})

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

/**
 * Download one ready File with translations
 */
func DownloadFile(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("project"))
	countryCode := c.Param("lang")
	delimeter := c.DefaultQuery("delimeter", ":")
	fileType := c.Param("type")

	if err == nil && len(countryCode) == 2 && len(delimeter) == 1 {

		fileContent, errFile := GenerateFile(projectID, countryCode, delimeter, generator.GetGenerator(fileType))

		if nil == errFile {
			c.Header("Content-Type", "application/txt; charset=utf-8")
			c.Header("Transfer-Encoding", "chunked")
			c.Header("Content-Disposition", "inline; filename=\"messages_"+countryCode+".properties\"")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.String(http.StatusOK, fileContent)
		} else {
			// error in the file content generation
			c.AbortWithStatus(http.StatusBadRequest)
		}

	} else {
		// bad parameters
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

/**
 * Add new term to project
 */
func PostNewTerm(c *gin.Context) {

	// key and projectId are mandatory
	termKey := c.PostForm("termKey")
	projectID, err := strconv.Atoi(c.PostForm("projectId"))

	if len(termKey) == 0 || err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// description is optional
		termDescr := c.PostForm("termDescr")
		addedTerm := AddNewTerm(dbFile, termKey, termDescr, projectID)
		c.JSON(http.StatusOK, gin.H{
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
	projectID, err := strconv.Atoi(c.Param("id"))

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

	mapLines, _ := readLines("/tmp/"+filename, delimeter)

	// save it somehow
	SaveImportedTermsForProject(dbFile, *mapLines, country, projectID)

	// delete temp file
	_ = os.Remove("/tmp/" + filename)

	// redirect
	c.Redirect(http.StatusMovedPermanently, "/project/"+strconv.Itoa(projectID)+"?imported")
}

/**
 * Delete one term and all its translations
 */
func DeleteOneTerm(c *gin.Context) {
	if termId, err := strconv.Atoi(c.Param("id")); err == nil {

		result := DeleteTerm(dbFile, termId)
		if result {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

/**
 * Delete one project
 */
func DeleteOneProject(c *gin.Context) {
	if projectId, err := strconv.Atoi(c.Param("id")); err == nil {

		result := DeleteProject(dbFile, projectId)
		if result {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

/**
 * Read lines for the given file
 * Returns the map of translations:
 */
func readLines(path string, delimeter string) (*map[string]ImportTranslation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mapLines := make(map[string]ImportTranslation)
	scanner := bufio.NewScanner(file)
	previousLineComment := ""
	currentLine := ""
	for scanner.Scan() {
		currentLine = scanner.Text()
		if currentLine[0] == '#' {
			// this line is the comment, remember it and move to the next line
			previousLineComment = currentLine[1:]

		} else if key, value, err := parseLine(currentLine, delimeter); err == nil {
			// that is normal line containing the pair of values key=translation
			mapLines[key] = ImportTranslation{Translation: value, Comment: previousLineComment}
			previousLineComment = ""
		}
	}
	return &mapLines, scanner.Err()
}

/**
 * Parse one line and return pair of values
 */
func parseLine(line string, delimeter string) (string, string, error) {
	parts := strings.Split(line, delimeter)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
	} else {
		return "", "", errors.New("Didn't found two parts in the incoming string")
	}
}

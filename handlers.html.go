package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

/**
 * Main page
 */
func MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H {
		"title": "Pojects",
		"projects": GetProjects(),
	})
}

/**
 * Get one given project
 */
func GetOneProject(c *gin.Context) {
	 id, _ := strconv.Atoi(c.Param("id"))
	 c.HTML(http.StatusOK, "project.tmpl", gin.H {
		 "title": "Pojects",
		 "projectId": id,
		 "currentLang": "all",
		 "project": FindOneProject(id, nil),
	 })
}

/**
 * Get one project with only untranslated terms
 */
func GetOneProjectUntranslated(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	countryCode := c.Param("lang")
	c.HTML(http.StatusOK, "project.tmpl", gin.H {
		"title": "Untranslated from " + countryCode,
		"currentLang": countryCode,
		"projectId": id,
		"project": FindOneProject(id, &countryCode),
	})
}

/**
 * Get all the translations for given term
 */
func GetAllTranslations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	term := GetTerm(id)
	c.HTML(http.StatusOK, "term.tmpl", gin.H {
		"title": "All translations",
		"term": term,
		"termId": id,
	})
}

package main

import (
	"net/http"
	"strconv"
	"monsieur-traducteur/services"
	"github.com/gin-gonic/gin"
)

/**
 * Main page
 */
func MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H {
		"title": "Pojects",
		"projects": services.GetProjects(),
	})
}

/**
 * Get one given project
 */
func GetOneProject(c *gin.Context) {
	 id, _ := strconv.Atoi(c.Param("id"))
	 c.HTML(http.StatusOK, "project.tmpl", gin.H {
		 "title": "Pojects",
		 "project": services.GetOneProject(id),
	 })
}

/**
 * Get all the translations for given term
 */
func GetAllTranslations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	term := services.GetTerm(id)
	c.HTML(http.StatusOK, "term.tmpl", gin.H {
		"title": "Pojects",
		"term": term,
		"termId": id,
	})
}

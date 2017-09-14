package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * Main page
 */
func MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":             "Pojects",
		"projects":          GetProjects(dbFile),
		"wasProjectDeleted": len(c.Query("projectdeleted")) > 0,
		"is_logged_in":      IsLoggedIn(c),
	})
}

/**
 * "Add new project" page
 */
func AddNewProjectForm(c *gin.Context) {
	c.HTML(http.StatusOK, "add-new-project.tmpl", gin.H{})
}

/**
 * Create new Project
 */
func PostCreateNewProject(c *gin.Context) {
	name := c.PostForm("name")
	country := c.PostForm("country")
	project := CreateNewProject(dbFile, name, country)
	c.Redirect(http.StatusMovedPermanently, "/project/"+strconv.Itoa(project.ID))
}

/**
 * Get one given project
 */
func GetOneProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	c.HTML(http.StatusOK, "project.tmpl", gin.H{
		"title":          "Pojects",
		"projectId":      id,
		"currentLang":    "all",
		"project":        FindOneProject(dbFile, id, ""),
		"wasTermDeleted": len(c.Query("termDeleted")) > 0,
		"is_logged_in":   IsLoggedIn(c),
	})
}

/**
 * Get one project with only untranslated terms
 */
func GetOneProjectUntranslated(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	countryCode := c.Param("lang")
	c.HTML(http.StatusOK, "project.tmpl", gin.H{
		"title":          "Untranslated from " + countryCode,
		"currentLang":    countryCode,
		"projectId":      id,
		"project":        FindOneProject(dbFile, id, countryCode),
		"wasTermDeleted": len(c.Query("termDeleted")) > 0,
		"is_logged_in":   IsLoggedIn(c),
	})
}

/**
 * Get all the translations for given term
 */
func GetAllTranslations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	term := GetTerm(dbFile, id)
	c.HTML(http.StatusOK, "term.tmpl", gin.H{
		"title":  "All translations",
		"term":   term,
		"termId": id,
	})
}

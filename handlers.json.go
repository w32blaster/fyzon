package main

import (
	"net/http"
	"monsieur-traducteur/services"
	"github.com/gin-gonic/gin"
	"strconv"
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

	services.UpdateTranslation(value, id, countryCode)

	c.JSON(http.StatusOK, gin.H {})
}

/**
 * Add new language
 */
func PostNewLanguage(c *gin.Context) {

		if projectId, err := strconv.Atoi(c.PostForm("projectId")); err == nil {

			countryCode := c.PostForm("countryCode")
	    services.AddNewLanguage(&projectId, &countryCode)
			c.JSON(http.StatusOK, gin.H {})

	} else {
      c.AbortWithStatus(http.StatusBadRequest)
	}

}

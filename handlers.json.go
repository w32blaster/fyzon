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
	value := c.PostForm("value")

	services.UpdateTranslation(value, id)

	c.JSON(http.StatusOK, gin.H {})
}

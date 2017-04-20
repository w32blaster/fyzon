package main

import (
  "github.com/gin-gonic/gin"
)

/**
 * Set up all the routes
 */
func initializeRoutes(router *gin.Engine) {

  // static assets
  router.Static("/semantic", "./semantic")
  router.Static("/assets", "./assets")

  // HTML endpoints
  router.GET("/", MainPage)
  router.GET("/project/:id", GetOneProject)

  // JSON endpoints
  router.GET("/ping", Ping)

  // Group API routes together
  api := router.Group("/api")
  {
    // update one term
		api.POST("/terms/:id/:lang", PostOneTerm)

    // get one term
    api.GET("/terms/:id", GetAllTranslations)
  }
}

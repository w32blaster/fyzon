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
  router.GET("/project/:id/untranslated/:lang", GetOneProjectUntranslated)

  router.POST("/project/:id/import", PostNewFile)

  router.GET("/add/new/project", AddNewProjectForm)

  // JSON endpoints
  router.GET("/ping", Ping)

  // Group API routes together
  api := router.Group("/api")
  {
    // update one term
		api.POST("/terms/:id/:lang", PostOneTerm)

    // get one term
    api.GET("/terms/:id", GetAllTranslations)

    // add new language to project
    api.POST("/project/add/language", PostNewLanguage)

    // add new term to project
    api.POST("/project/add/term", PostNewTerm)

    api.POST("/project/add", PostCreateNewProject)

    // delete one term and all its translations
    api.DELETE("/terms/:id", DeleteOneTerm)

    // delete one project and a-a-a-a-a-l its terms and translations
    api.DELETE("/project/:id", DeleteOneProject)
  }
}

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
  router.GET("/project/:id", ensureLoggedIn(), GetOneProject)
  router.GET("/project/:id/untranslated/:lang", ensureLoggedIn(), GetOneProjectUntranslated)

  // Import (upload) new messages file
  router.POST("/project/:id/import", ensureLoggedIn(), PostNewFile)

  router.GET("/add/new/project", ensureLoggedIn(), AddNewProjectForm)

  // Group API routes together
  api := router.Group("/api")
  {
    // update one term
		api.POST("/terms/:id/:lang", ensureLoggedIn(), PostOneTerm)

    // get one term
    api.GET("/terms/:id", ensureLoggedIn(), GetAllTranslations)

    // add new language to project
    api.POST("/project/add/language", ensureLoggedIn(), PostNewLanguage)

    // add new term to project
    api.POST("/project/add/term", ensureLoggedIn(), PostNewTerm)

    api.POST("/project/add", ensureLoggedIn(), PostCreateNewProject)

    // delete one term and all its translations
    api.DELETE("/terms/:id", ensureLoggedIn(), DeleteOneTerm)

    // delete one project and a-a-a-a-a-l its terms and translations
    api.DELETE("/project/:id", ensureLoggedIn(), DeleteOneProject)
  }

  // Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", ensureLoggedIn(), logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
}
}

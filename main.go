package main

import (
  "github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router := gin.Default()

  router.LoadHTMLGlob("templates/*")

  // initialize all the routes
  initializeRoutes(router)

	router.Run() // listen and serve on 0.0.0.0:8080
}

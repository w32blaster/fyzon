package main

import (
  "github.com/gin-gonic/gin"
)

var router *gin.Engine

const (
    dbFile = "./trans.sqlite3"
)

func main() {
	router := gin.Default()

  // indicating whether the request was from an authenticated user or not
  router.Use(setUserStatus())
  
  router.LoadHTMLGlob("templates/*")

  // initialize all the routes
  initializeRoutes(router)

	router.Run() // listen and serve on 0.0.0.0:8080
}

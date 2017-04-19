package main

import (
  "github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

  // static assets
  router.Static("/semantic", "./semantic")
  router.Static("/assets", "./assets")

  // HTML endpoints
  router.GET("/", MainPage)
  router.GET("/project/:id", GetOneProject)
  router.GET("/api/terms/:id", GetAllTerms)

  // JSON endpoints
  router.GET("/ping", Ping)
  router.POST("/api/terms/:id", PostOneTerm)
}

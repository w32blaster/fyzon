package main

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "monsieur-traducteur/services"
)

func main() {
	r := gin.Default()
  r.LoadHTMLGlob("templates/*")

  r.GET("/", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "index.tmpl", gin.H {
  			"title": "Posts",
  		})
  })

  r.GET("/projects", func(c *gin.Context) {
      c.HTML(http.StatusOK, "projects.tmpl", gin.H {
  			"title": "Pojects",
  		})
  })

  // Ping test
	r.GET("/ping", func(c *gin.Context) {
    
    projectList := services.GetProjects()
		c.JSON(200, gin.H {
			"message": "pong",
      "projects": projectList })
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

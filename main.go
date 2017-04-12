package main

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "monsieur-traducteur/services"
  "strconv"
)

func main() {
	r := gin.Default()
  r.LoadHTMLGlob("templates/*")

  // static assets
  r.Static("/semantic", "./semantic")
  r.Static("/assets", "./assets")

  r.GET("/", func(c *gin.Context) {
     c.HTML(http.StatusOK, "index.tmpl", gin.H {
       "title": "Pojects",
       "projects": services.GetProjects(),
     })
  })

  r.GET("/project/:id", func(c *gin.Context) {
     id, _ := strconv.Atoi(c.Param("id"))
     c.HTML(http.StatusOK, "project.tmpl", gin.H {
       "title": "Pojects",
       "project": services.GetOneProject(id),
     })
  })

  // Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H {
			"message": "pong",
    })
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

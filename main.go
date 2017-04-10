package main

import (
  "database/sql"
  "github.com/gin-gonic/gin"
  _ "github.com/mattn/go-sqlite3"
  "net/http"
  "log"
  "bytes"
)

func main() {
	r := gin.Default()
  r.LoadHTMLGlob("templates/*")

  const dbFile = "./trans.sqlite3"

  r.GET("/", func(c *gin.Context) {
  		c.HTML(http.StatusOK, "index.tmpl", gin.H {
  			"title": "Posts",
  		})
  })

  r.GET("/projects", func(c *gin.Context) {
      db, err := sql.Open("sqlite3", dbFile)
      if err != nil {
      		log.Fatal(err)
      	}
      defer db.Close()

  		c.HTML(http.StatusOK, "projects.tmpl", gin.H {
  			"title": "Pojects",
  		})
  })

  // Ping test
	r.GET("/ping", func(c *gin.Context) {
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("select id, name from projects")
  	if err != nil {
  		log.Fatal(err)
  	}
    defer rows.Close()

    var buffer bytes.Buffer
    for rows.Next() {
  		var id int
  		var name string
  		err = rows.Scan(&id, &name)
  		if err != nil {
  			log.Fatal(err)
  		}
  		buffer.WriteString(name)
      buffer.WriteString(" ,")
    }

		c.JSON(200, gin.H {
			"message": "pong",
      "projects": buffer.String()})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

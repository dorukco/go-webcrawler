package main

import (
	"go-webcrawler/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.IndexHandler)
	r.POST("/submit", handlers.SubmitHandler)

	r.Run(":8080")
}

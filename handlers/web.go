package handlers

import (
	"fmt"
	"go-webcrawler/crawler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func SubmitHandler(c *gin.Context) {
	textInput := strings.TrimSpace(c.PostForm("text_input"))

	if textInput == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "Please enter a URL!",
		})
		return
	}

	if !crawler.IsValidURL(textInput) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error":       "Please enter a valid URL (must start with http:// or https://)",
			"input_value": textInput,
		})
		return
	}

	fmt.Printf("WebCrawler processing URL: %s\n", textInput)

	result := crawler.CrawlURL(textInput)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"result": result,
	})
}

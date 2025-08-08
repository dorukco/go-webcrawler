package crawler

import (
	"fmt"
	"go-webcrawler/models"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

const (
	RequestTimeout = 10 * time.Second
)

func CrawlURL(url string) models.CrawlResult {
	result := models.CrawlResult{
		URL: url,
	}

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		result.Error = fmt.Sprintf("Network error: %v", err)
		result.Success = false
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Status = resp.Status

	if resp.StatusCode != http.StatusOK {
		result.Error = GetStatusCodeDescription(resp.StatusCode)
		result.Success = false
		return result
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to parse HTML: %v", err)
		result.Success = false
		return result
	}

	// Extract title
	title := ExtractTitle(doc)
	if title == "" {
		result.Title = "No title found"
	} else {
		result.Title = title
	}

	result.Success = true
	return result
}

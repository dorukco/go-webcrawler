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
	normalizedURL := NormalizeURL(url)

	result := models.CrawlResult{
		URL: normalizedURL,
	}

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	req, err := http.NewRequest("GET", normalizedURL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Success = false
		return result
	}

	// Add headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
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

	// Extract HTML version and DOCTYPE
	result.HTMLVersion, result.DocType = ExtractHTMLVersion(doc)

	// Extract headings
	result.Headings = ExtractHeadings(doc)

	// Detect login form
	result.HasLoginForm = DetectLoginForm(doc)

	// Extract link information
	result.InternalLinks, result.ExternalLinks, result.InaccessibleLinks = ExtractLinks(doc, normalizedURL)

	result.Success = true
	return result
}

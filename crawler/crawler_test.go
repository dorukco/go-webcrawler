package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawlURL_InvalidURL(t *testing.T) {
	result := CrawlURL("invalid-url")

	if result.Success {
		t.Error("Expected CrawlURL to fail for invalid URL")
	}

	if result.Error == "" {
		t.Error("Expected error message for invalid URL")
	}
}

func TestCrawlURL_ValidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test Page</title>
				</head>
				<body>
					<h1>Main Heading</h1>
					<h2>Sub Heading</h2>
					<a href="/internal">Internal Link</a>
					<a href="https://external.com">External Link</a>
					<form id="login-form">
						<input type="text" name="username">
						<input type="password" name="password">
					</form>
				</body>
			</html>
		`)
	}))
	defer server.Close()

	result := CrawlURL(server.URL)

	if !result.Success {
		t.Errorf("Expected successful crawl, got error: %s", result.Error)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if result.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got %q", result.Title)
	}

	if result.HTMLVersion != "HTML5" {
		t.Errorf("Expected HTML5, got %q", result.HTMLVersion)
	}

	if !result.HasLoginForm {
		t.Error("Expected to detect login form")
	}

	if result.Headings["h1"] != 1 {
		t.Errorf("Expected 1 h1 heading, got %d", result.Headings["h1"])
	}

	if result.Headings["h2"] != 1 {
		t.Errorf("Expected 1 h2 heading, got %d", result.Headings["h2"])
	}

	if result.InternalLinks != 1 {
		t.Errorf("Expected 1 internal link, got %d", result.InternalLinks)
	}

	if result.ExternalLinks != 1 {
		t.Errorf("Expected 1 external link, got %d", result.ExternalLinks)
	}
}

func TestCrawlURL_404Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
	}))
	defer server.Close()

	result := CrawlURL(server.URL)

	if result.Success {
		t.Error("Expected crawl to fail for 404 response")
	}

	if result.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", result.StatusCode)
	}

	if result.Error == "" {
		t.Error("Expected error message for 404 response")
	}
}

func TestCrawlURL_NoTitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<html><body><h1>No Title Page</h1></body></html>`)
	}))
	defer server.Close()

	result := CrawlURL(server.URL)

	if !result.Success {
		t.Errorf("Expected successful crawl, got error: %s", result.Error)
	}

	if result.Title != "No title found" {
		t.Errorf("Expected 'No title found', got %q", result.Title)
	}
}

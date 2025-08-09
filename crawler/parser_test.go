package crawler

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		htmlStr  string
		expected string
	}{
		{
			"Basic title",
			"<html><head><title>Test Title</title></head></html>",
			"Test Title",
		},
		{
			"Title with whitespace",
			"<html><head><title>  Test Title  </title></head></html>",
			"Test Title",
		},
		{
			"No title",
			"<html><head></head></html>",
			"",
		},
		{
			"Empty title",
			"<html><head><title></title></head></html>",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tt.htmlStr))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			result := ExtractTitle(doc)
			if result != tt.expected {
				t.Errorf("ExtractTitle() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractHeadings(t *testing.T) {
	htmlStr := `
	<html>
		<body>
			<h1>Heading 1</h1>
			<h2>Heading 2</h2>
			<h2>Another H2</h2>
			<h3>Heading 3</h3>
			<div>Not a heading</div>
		</body>
	</html>`

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	result := ExtractHeadings(doc)

	expected := map[string]int{
		"h1": 1,
		"h2": 2,
		"h3": 1,
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d heading types, got %d", len(expected), len(result))
	}

	for tag, count := range expected {
		if result[tag] != count {
			t.Errorf("Expected %d %s tags, got %d", count, tag, result[tag])
		}
	}
}

func TestDetectLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		htmlStr  string
		expected bool
	}{
		{
			"Form with login ID",
			`<html><body><form id="login-form"><input type="text"></form></body></html>`,
			true,
		},
		{
			"Form with signin class",
			`<html><body><form class="signin-form"><input type="text"></form></body></html>`,
			true,
		},
		{
			"Form with auth class",
			`<html><body><form class="auth-modal"><input type="text"></form></body></html>`,
			true,
		},
		{
			"Regular form",
			`<html><body><form id="contact-form"><input type="text"></form></body></html>`,
			false,
		},
		{
			"No form",
			`<html><body><div>No form here</div></body></html>`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tt.htmlStr))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			result := DetectLoginForm(doc)
			if result != tt.expected {
				t.Errorf("DetectLoginForm() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"HTTPS with www", "https://www.doruk.com/path", "doruk.com"},
		{"HTTP without www", "http://doruk.com", "doruk.com"},
		{"With port", "https://doruk.com:8080/path", "doruk.com:8080"},
		{"Subdomain", "https://api.doruk.com", "api.doruk.com"},
		{"No protocol", "doruk.com/path", "doruk.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDomain(tt.url)
			if result != tt.expected {
				t.Errorf("extractDomain(%q) = %q; want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	htmlStr := `
	<html>
		<body>
			<a href="/">Internal root</a>
			<a href="/page">Internal page</a>
			<a href="https://doruk.com/internal">Internal absolute</a>
			<a href="https://doruk.com">External</a>
			<a href="mailto:test@doruk.com">Email</a>
			<a href="#section">Fragment</a>
			<a>No href</a>
		</body>
	</html>`

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	internal, external, inaccessible := ExtractLinks(doc, "https://doruk.com")

	if internal != 3 {
		t.Errorf("Expected 3 internal links, got %d", internal)
	}
	if external != 1 {
		t.Errorf("Expected 1 external link, got %d", external)
	}
	if inaccessible != 3 {
		t.Errorf("Expected 3 inaccessible links, got %d", inaccessible)
	}
}

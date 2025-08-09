package crawler

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"Valid HTTP URL", "http://doruk.com", true},
		{"Valid HTTPS URL", "https://doruk.com", true},
		{"Valid URL with www", "https://www.doruk.com", true},
		{"Valid URL without protocol", "doruk.com", true},
		{"Valid URL with path", "https://doruk.com/path", true},
		{"Valid URL with query", "https://doruk.com?query=1", true},
		{"Invalid URL - no domain", "http://", false},
		{"Invalid URL - malformed", "not-a-url", false},
		{"Empty string", "", false},
		{"Just protocol", "https://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Already HTTP", "http://doruk.com", "http://doruk.com"},
		{"Already HTTPS", "https://doruk.com", "https://doruk.com"},
		{"No protocol", "doruk.com", "https://doruk.com"},
		{"With spaces", "  doruk.com  ", "https://doruk.com"},
		{"Mixed case", "Doruk.COM", "https://Doruk.COM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeURL(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetStatusCodeDescription(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   string
	}{
		{"Bad Request", 400, "Bad Request - The server cannot process the request"},
		{"Unauthorized", 401, "Unauthorized - Authentication is required"},
		{"Forbidden", 403, "Forbidden - Access to this resource is denied"},
		{"Not Found", 404, "Not Found - The requested page does not exist"},
		{"Internal Server Error", 500, "Internal Server Error - The server encountered an error"},
		{"Other 4xx", 418, "Client Error - There's an issue with the request"},
		{"Other 5xx", 502, "Server Error - The server encountered an error"},
		{"Unexpected", 200, "Unexpected status code"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusCodeDescription(tt.statusCode)
			if result != tt.expected {
				t.Errorf("GetStatusCodeDescription(%d) = %q; want %q", tt.statusCode, result, tt.expected)
			}
		})
	}
}

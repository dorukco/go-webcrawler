package crawler

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return strings.TrimSpace(n.FirstChild.Data)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := ExtractTitle(c); title != "" {
			return title
		}
	}

	return ""
}

func ExtractHTMLVersion(n *html.Node) (string, string) {
	doctype := findDocType(n)

	if doctype == "" {
		return "Unknown", "No DOCTYPE found"
	}

	version := analyzeDocType(doctype)
	return version, doctype
}

func findDocType(n *html.Node) string {
	if n.Type == html.DoctypeNode {
		return strings.TrimSpace(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if doctype := findDocType(c); doctype != "" {
			return doctype
		}
	}

	return ""
}

func analyzeDocType(doctype string) string {
	doctype = strings.ToLower(strings.TrimSpace(doctype))

	switch {
	case doctype == "html":
		return "HTML5"
	case strings.Contains(doctype, "html 4.01 strict"):
		return "HTML 4.01 Strict"
	case strings.Contains(doctype, "html 4.01 transitional"):
		return "HTML 4.01 Transitional"
	case strings.Contains(doctype, "xhtml 1.0"):
		return "XHTML 1.0"
	case strings.Contains(doctype, "xhtml 1.1"):
		return "XHTML 1.1"
	default:
		return "Unknown/Custom"
	}
}

func ExtractHeadings(n *html.Node) map[string]int {
	headings := make(map[string]int)
	countHeadings(n, headings)
	return headings
}

func countHeadings(n *html.Node, headings map[string]int) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			headings[n.Data]++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countHeadings(c, headings)
	}
}

func DetectLoginForm(n *html.Node) bool {
	return hasLoginForm(n)
}

func hasLoginForm(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "form" {
		if isLoginForm(n) {
			return true
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hasLoginForm(c) {
			return true
		}
	}

	return false
}

func isLoginForm(formNode *html.Node) bool {
	for _, attr := range formNode.Attr {
		if attr.Key == "id" || attr.Key == "class" {
			value := strings.ToLower(attr.Val)
			if strings.Contains(value, "login") || strings.Contains(value, "signin") || strings.Contains(value, "auth") {
				return true
			}
		}
	}

	return false
}

func ExtractLinks(n *html.Node, baseURL string) (int, int, int) {
	internal := 0
	external := 0
	inaccessible := 0

	countLinks(n, baseURL, &internal, &external, &inaccessible)
	return internal, external, inaccessible
}

func countLinks(n *html.Node, baseURL string, internal, external, inaccessible *int) {
	if n.Type == html.ElementNode && n.Data == "a" {
		href := getHrefAttribute(n)
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "mailto:") {
			*inaccessible++
		} else if strings.HasPrefix(href, "http") {
			if strings.Contains(href, extractDomain(baseURL)) {
				*internal++
			} else {
				*external++
			}
		} else {
			*internal++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countLinks(c, baseURL, internal, external, inaccessible)
	}
}

func getHrefAttribute(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return strings.TrimSpace(attr.Val)
		}
	}
	return ""
}

func extractDomain(url string) string {
	url = strings.ToLower(url)
	if strings.HasPrefix(url, "https://") {
		url = url[8:]
	} else if strings.HasPrefix(url, "http://") {
		url = url[7:]
	}
	if strings.HasPrefix(url, "www.") {
		url = url[4:]
	}
	if idx := strings.IndexAny(url, "/?"); idx != -1 {
		url = url[:idx]
	}
	return url
}

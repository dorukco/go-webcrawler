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

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

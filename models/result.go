package models

type CrawlResult struct {
	URL               string
	StatusCode        int
	Status            string
	Title             string
	HTMLVersion       string
	DocType           string
	Headings          map[string]int
	HasLoginForm      bool
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	Error             string
	Success           bool
}

package models

type CrawlResult struct {
	URL        string
	StatusCode int
	Status     string
	Title      string
	Error      string
	Success    bool
}

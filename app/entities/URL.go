package entities

import "time"

// URL is a URL and its corresponding original URL
type URL struct {
	URL         string
	OriginalURL string
}

// VisitDetails is details of a single visits of a URL
type VisitDetails struct {
	IP        string
	Time      time.Time
	UserAgent string
	ShortURL  string
}

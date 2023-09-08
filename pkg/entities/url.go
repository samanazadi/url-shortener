package entities

import "time"

// URL is a URL and its corresponding original URL
type URL struct {
	ShortURL    string
	OriginalURL string
}

// VisitDetail is details of a single visits of a URL
type VisitDetail struct {
	IP        string
	Time      time.Time
	UserAgent string
	ShortURL  string
}

package JSON

import "time"

type VisitDetail struct {
	IP        string    `json:"ip"`
	Time      time.Time `json:"time"`
	UserAgent string    `json:"agent"`
}

type SuccessRetrieval struct {
	Message      string        `json:"message"`
	OriginalURL  string        `json:"original"`
	Total        int           `json:"total"`
	VisitDetails []VisitDetail `json:"visits"`
}

type UnsuccessfulRetrieval struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

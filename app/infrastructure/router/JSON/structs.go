package JSON

type SuccessRetrieval struct {
	Message     string
	OriginalURL string
}

type UnsuccessfulRetrieval struct {
	Message string
	Error   string
}

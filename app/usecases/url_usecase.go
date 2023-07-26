package usecases

import (
	"github.com/samanazadi/url-shortener/app/entities"
	"time"
)

// URLUsecase implements url fetching logic
type URLUsecase struct {
	URLRepository URLRepository
}

// OriginalURL returns original URL or an error
func (u URLUsecase) OriginalURL(url string) (string, error) {
	ue, err := u.URLRepository.FindURL(url)
	if err != nil {
		return "", err
	}
	return ue.OriginalURL, nil
}

func (u URLUsecase) SaveVisitDetails(vd VisitDetails) error {
	return u.URLRepository.SaveVisitDetails(vd)
}

// URLRepository defines abstract repository operations
type URLRepository interface {
	FindURL(string) (entities.URL, error)
	SaveVisitDetails(VisitDetails) error
}

// VisitDetails is details of a single visits of a URL
type VisitDetails struct {
	IP        string
	Time      time.Time
	UserAgent string
	ShortURL  string
}

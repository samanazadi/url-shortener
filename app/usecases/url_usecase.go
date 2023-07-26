package usecases

import (
	"github.com/samanazadi/url-shortener/app/entities"
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

func (u URLUsecase) SaveVisitDetails(vd entities.VisitDetails) error {
	return u.URLRepository.SaveVisitDetails(vd)
}

// URLRepository defines abstract repository operations
type URLRepository interface {
	FindURL(string) (entities.URL, error)
	SaveVisitDetails(entities.VisitDetails) error
}

package usecases

import (
	"github.com/samanazadi/url-shortener/pkg/base62"
	"github.com/samanazadi/url-shortener/pkg/entities"
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

func (u URLUsecase) SaveURL(url string, machineID uint16) (string, error) {
	shortURL := base62.GenerateID(machineID)
	if err := u.URLRepository.SaveShortURL(url, shortURL); err != nil {
		return "", nil
	}
	return shortURL, nil
}

func (u URLUsecase) SaveVisitDetail(vd entities.VisitDetail) error {
	return u.URLRepository.SaveVisitDetail(vd)
}

func (u URLUsecase) Visits(url string, offset int, limit int) ([]entities.VisitDetail, int, error) {
	total := u.URLRepository.TotalVisits(url)
	vds, err := u.URLRepository.FindVisits(url, offset, limit)
	return vds, total, err
}

// URLRepository defines abstract repository operations
type URLRepository interface {
	FindURL(string) (entities.URL, error)
	SaveVisitDetail(entities.VisitDetail) error
	SaveShortURL(u string, s string) error
	FindVisits(u string, offset int, limit int) ([]entities.VisitDetail, error)
	TotalVisits(u string) int
}

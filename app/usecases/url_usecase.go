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
	FindVisits(u string, offset int, limit int) ([]entities.VisitDetail, error)
	TotalVisits(u string) int
}

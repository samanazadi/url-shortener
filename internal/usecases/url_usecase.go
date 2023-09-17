package usecases

import (
	"context"
	"github.com/samanazadi/url-shortener/pkg/base62"
	"github.com/samanazadi/url-shortener/pkg/entities"
)

// URLUsecase implements url fetching logic
type URLUsecase struct {
	URLRepository URLRepository
}

// OriginalURL returns original URL or an error
func (u URLUsecase) OriginalURL(ctx context.Context, url string) (string, error) {
	ue, err := u.URLRepository.FindURL(ctx, url)
	if err != nil {
		return "", err
	}
	return ue.OriginalURL, nil
}

func (u URLUsecase) SaveURL(ctx context.Context, url string, machineID uint16) (string, error) {
	shortURL := base62.GenerateID(machineID)
	if err := u.URLRepository.SaveShortURL(ctx, url, shortURL); err != nil {
		return "", nil
	}
	return shortURL, nil
}

func (u URLUsecase) SaveVisitDetail(ctx context.Context, vd entities.VisitDetail) error {
	return u.URLRepository.SaveVisitDetail(ctx, vd)
}

func (u URLUsecase) Visits(ctx context.Context, url string, offset int, limit int) ([]entities.VisitDetail, int, error) {
	total, err1 := u.URLRepository.TotalVisits(ctx, url)
	vds, err2 := u.URLRepository.FindVisits(ctx, url, offset, limit)
	if err1 != nil {
		return vds, total, err1
	} else if err2 != nil {
		return vds, total, err2
	} else {
		return vds, total, nil
	}

}

// URLRepository defines abstract repository operations
type URLRepository interface {
	FindURL(context.Context, string) (entities.URL, error)
	SaveVisitDetail(context.Context, entities.VisitDetail) error
	SaveShortURL(ctx context.Context, u string, s string) error
	FindVisits(ctx context.Context, u string, offset int, limit int) ([]entities.VisitDetail, error)
	TotalVisits(ctx context.Context, u string) (int, error)
}

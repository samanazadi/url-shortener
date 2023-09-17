package controllers

import (
	"context"
	"github.com/samanazadi/url-shortener/internal/adapters/controllers/models"
	"github.com/samanazadi/url-shortener/internal/config"
	"github.com/samanazadi/url-shortener/internal/usecases"
	"github.com/samanazadi/url-shortener/pkg/entities"
	"github.com/samanazadi/url-shortener/pkg/logging"
	"gorm.io/gorm"
	"strconv"
)

// URLController is responsible for redirecting user
type URLController struct {
	urlUseCase usecases.URLUsecase
}

// NewURLController is a creator function for URLController
func NewURLController(db *gorm.DB) *URLController {
	return &URLController{
		urlUseCase: usecases.URLUsecase{
			URLRepository: URLControllerRepository{
				DB: db,
			},
		},
	}
}

// GetDetails retrieves the original input
func (u URLController) GetDetails(ctx context.Context, p URLControllerInputPort, cfg *config.Config) {
	shortURL := p.Param("id")
	offset := AtoIWithDefault(p.Param("offset"), 0)
	limit := AtoIWithDefault(p.Param("limit"), cfg.DefaultLimit)
	originalURL, err := u.urlUseCase.OriginalURL(ctx, shortURL)
	if err != nil {
		logging.Logger.Warn(err.Error(), "short_url", shortURL)
		p.OutputError(URLNotFound, err)
		return
	}
	vds, total, err := u.urlUseCase.Visits(ctx, shortURL, offset, limit)
	if err != nil {
		logging.Logger.Error(err.Error(), "short_url", shortURL)
	}

	logging.Logger.Debug("visit details returned successfully",
		"short_url", shortURL, "offset", offset, "limit", limit)
	p.OutputVisitDetails(originalURL, vds, total)
}

// CreateShortLink create a short link for the URL in request
func (u URLController) CreateShortLink(ctx context.Context, p URLControllerInputPort, cfg *config.Config) {
	originalURL, err := p.GetCreateShortURLRequest()
	if err != nil {
		logging.Logger.Warn(err.Error())
		p.OutputError(BadRequest, err)
		return
	}
	shortURL, err := u.urlUseCase.SaveURL(ctx, originalURL, uint16(cfg.MachineID))
	if err != nil {
		logging.Logger.Error(err.Error())
		p.OutputError(CannotCreateShortLink, err)
		return
	}
	logging.Logger.Debug("Short url created", "short_url", shortURL, "original_url", originalURL)
	p.OutputShortURL(shortURL)
}

// RedirectToOriginalURL redirects to original URL is exists and redirects to homepage otherwise
func (u URLController) RedirectToOriginalURL(ctx context.Context, p URLControllerInputPort) {
	shortURL := p.Param("id")
	originalURL, err := u.urlUseCase.OriginalURL(ctx, shortURL)
	if err != nil {
		logging.Logger.Warn("Short url not found", "short_url", shortURL)
		p.OutputError(RedirectToHomePage, err)
		return
	}
	vd := p.GetVisitDetail()
	vd.ShortURL = shortURL
	err = u.urlUseCase.SaveVisitDetail(ctx, vd)
	if err != nil {
		logging.Logger.Error("Cannot save visit details: "+err.Error(), "short_url", shortURL)
	}
	logging.Logger.Debug("Redirect successful", "short_url", shortURL, "original_url", originalURL)
	p.Redirect(originalURL)
}

const (
	// URLNotFound is an unsuccessful retrieval of a URL
	URLNotFound = iota
	// RedirectToHomePage unsuccessful and redirect to homepage
	RedirectToHomePage
	CannotCreateShortLink
	BadRequest
)

// URLControllerInputPort will be injected by infrastructure layer
type URLControllerInputPort interface {
	Param(string) string
	GetCreateShortURLRequest() (string, error)
	GetVisitDetail() entities.VisitDetail
	OutputShortURL(string)
	OutputVisitDetails(string, []entities.VisitDetail, int)
	OutputError(int, error)
	Redirect(u string)
}

// URLControllerRepository is an implementation of usecases.URLRepository
type URLControllerRepository struct {
	DB *gorm.DB
}

func (r URLControllerRepository) SaveShortURL(ctx context.Context, original string, short string) error {
	u := &models.URL{
		ShortURL:    short,
		OriginalURL: original,
	}
	result := r.DB.WithContext(ctx).Create(u)
	return result.Error
}

// FindURL queries the database for specified URL
func (r URLControllerRepository) FindURL(ctx context.Context, short string) (entities.URL, error) {
	var u models.URL
	result := r.DB.WithContext(ctx).Where("short_url", short).First(&u)
	if result.Error != nil {
		return entities.URL{}, result.Error
	}
	return entities.URL{
		ShortURL:    u.ShortURL,
		OriginalURL: u.OriginalURL,
	}, nil
}

func (r URLControllerRepository) SaveVisitDetail(ctx context.Context, vd entities.VisitDetail) error {
	v := &models.Visit{
		IP:        vd.IP,
		Time:      vd.Time,
		UserAgent: &vd.UserAgent,
		ShortURL:  vd.ShortURL,
	}
	result := r.DB.WithContext(ctx).Create(v)
	return result.Error
}

func (r URLControllerRepository) FindVisits(ctx context.Context, short string, offset int, limit int) ([]entities.VisitDetail, error) {
	var visits []models.Visit
	result := r.DB.WithContext(ctx).Model(&models.Visit{}).Where("short_url = ?", short).Limit(limit).Offset(offset).Find(&visits)
	if result.Error != nil {
		return nil, result.Error
	}
	vds := make([]entities.VisitDetail, 0)
	for _, v := range visits {
		vds = append(vds, entities.VisitDetail{
			IP:        v.IP,
			Time:      v.Time,
			UserAgent: *v.UserAgent,
			ShortURL:  v.ShortURL,
		})
	}
	return vds, nil
}

func (r URLControllerRepository) TotalVisits(ctx context.Context, u string) (int, error) {
	var count int64
	result := r.DB.WithContext(ctx).Model(&models.Visit{}).Where("short_url = ?", u).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

func AtoIWithDefault(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

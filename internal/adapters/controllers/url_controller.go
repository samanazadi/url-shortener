package controllers

import (
	"context"
	"github.com/samanazadi/url-shortener/internal/config"
	"github.com/samanazadi/url-shortener/internal/usecases"
	"github.com/samanazadi/url-shortener/pkg/entities"
	"github.com/samanazadi/url-shortener/pkg/logging"
	"strconv"
	"time"
)

// URLController is responsible for redirecting user
type URLController struct {
	urlUseCase usecases.URLUsecase
}

// NewURLController is a creator function for URLController
func NewURLController(h SQLHandler) *URLController {
	return &URLController{
		urlUseCase: usecases.URLUsecase{
			URLRepository: URLControllerRepository{
				SQLHandler: h,
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

// SQLHandler will be injected by infrastructure layer
type SQLHandler interface {
	ExecContext(context.Context, string, ...any) (Result, error)
	QueryRowContext(context.Context, string, ...any) Row
	QueryContext(context.Context, string, ...any) (Rows, error)
	Close() error
}

// Result is a SQL Exec result
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Row is one row in a SQL table
type Row interface {
	Scan(...any) error
}

type Rows interface {
	Scan(...any) error
	Next() bool
}

// URLControllerRepository is an implementation of usecases.URLRepository
type URLControllerRepository struct {
	SQLHandler SQLHandler
}

func (r URLControllerRepository) SaveShortURL(ctx context.Context, u string, s string) error {
	_, err := r.SQLHandler.ExecContext(ctx, "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)",
		s, u)
	return err
}

// FindURL queries the database for specified URL
func (r URLControllerRepository) FindURL(ctx context.Context, u string) (entities.URL, error) {
	row := r.SQLHandler.QueryRowContext(ctx, "SELECT * FROM urls WHERE short_url = $1", u)
	var shortURL, originalURL string
	err := row.Scan(&shortURL, &originalURL)
	if err != nil {
		return entities.URL{}, err
	}
	return entities.URL{ShortURL: shortURL, OriginalURL: originalURL}, nil
}

func (r URLControllerRepository) SaveVisitDetail(ctx context.Context, vd entities.VisitDetail) error {
	_, err := r.SQLHandler.ExecContext(ctx, "INSERT INTO visits (ip, time, user_agent, short_url) VALUES ($1, $2, $3, $4)",
		vd.IP, vd.Time, vd.UserAgent, vd.ShortURL)
	return err
}

func (r URLControllerRepository) FindVisits(ctx context.Context, u string, offset int, limit int) ([]entities.VisitDetail, error) {
	rows, err := r.SQLHandler.QueryContext(ctx, "SELECT ip, time, user_agent, short_url FROM visits WHERE short_url = $1 LIMIT $2 OFFSET $3",
		u, limit, offset)
	if err != nil {
		return nil, err
	}
	vds := make([]entities.VisitDetail, 0)
	var (
		ip        string
		t         time.Time
		userAgent string
		shortURL  string
	)
	for rows.Next() {
		err = rows.Scan(&ip, &t, &userAgent, &shortURL)
		if err != nil {
			return vds, err
		}
		vds = append(vds, entities.VisitDetail{IP: ip, Time: t, UserAgent: userAgent, ShortURL: shortURL})
	}
	return vds, nil
}

func (r URLControllerRepository) TotalVisits(ctx context.Context, u string) int {
	row := r.SQLHandler.QueryRowContext(ctx, "SELECT COUNT(*) FROM visits WHERE short_url = $1", u)
	var total int = 0
	if err := row.Scan(&total); err != nil {
		return 0
	}
	return total
}

func AtoIWithDefault(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

package controllers

import (
	"github.com/samanazadi/url-shortener/app/entities"
	"github.com/samanazadi/url-shortener/app/usecases"
	"log"
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
func (u URLController) GetDetails(p URLControllerInputPort) {
	url := p.Param("id")
	offset := AtoIWithDefault(p.Param("offset"), 0)
	limit := AtoIWithDefault(p.Param("limit"), 10)
	originalURL, err := u.urlUseCase.OriginalURL(url)
	if err != nil {
		p.OutputError(URLNotFound, err)
		return
	}
	vds, total, _ := u.urlUseCase.Visits(url, offset, limit)

	p.Output(originalURL, vds, total)
}

// RedirectToOriginalURL redirects to original URL is exists and redirects to homepage otherwise
func (u URLController) RedirectToOriginalURL(p URLControllerInputPort) {
	shortURL := p.Param("id")
	originalURL, err := u.urlUseCase.OriginalURL(shortURL)
	if err != nil {
		p.OutputError(RedirectToHomePage, err)
		return
	}
	vd := p.GetVisitDetail()
	vd.ShortURL = shortURL
	err = u.urlUseCase.SaveVisitDetail(vd)
	if err != nil {
		log.Printf("Cannot save visit: %s", err)
	}
	p.Redirect(originalURL)
}

const (
	// URLNotFound is an unsuccessful retrieval of a URL
	URLNotFound = iota
	// RedirectToHomePage unsuccessful and redirect to homepage
	RedirectToHomePage
)

// URLControllerInputPort will be injected by infrastructure layer
type URLControllerInputPort interface {
	Param(string) string
	GetVisitDetail() entities.VisitDetail
	Output(string, []entities.VisitDetail, int)
	Redirect(u string)
	OutputError(int, error)
}

// SQLHandler will be injected by infrastructure layer
type SQLHandler interface {
	Exec(string, ...any) (Result, error)
	QueryRow(string, ...any) Row
	Query(string, ...any) (Rows, error)
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

// FindURL queries the database for specified URL
func (r URLControllerRepository) FindURL(u string) (entities.URL, error) {
	row := r.SQLHandler.QueryRow("SELECT * FROM urls WHERE short_url = $1", u)
	var shortURL, originalURL string
	err := row.Scan(&shortURL, &originalURL)
	if err != nil {
		return entities.URL{}, err
	}
	return entities.URL{URL: shortURL, OriginalURL: originalURL}, nil
}

func (r URLControllerRepository) SaveVisitDetail(vd entities.VisitDetail) error {
	_, err := r.SQLHandler.Exec("INSERT INTO visits (ip, time, user_agent, short_url) VALUES ($1, $2, $3, $4)",
		vd.IP, vd.Time, vd.UserAgent, vd.ShortURL)
	return err
}

func (r URLControllerRepository) FindVisits(u string, offset int, limit int) ([]entities.VisitDetail, error) {
	rows, err := r.SQLHandler.Query("SELECT ip, time, user_agent, short_url FROM visits WHERE short_url = $1 LIMIT $2 OFFSET $3",
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

func (r URLControllerRepository) TotalVisits(u string) int {
	row := r.SQLHandler.QueryRow("SELECT COUNT(*) FROM visits WHERE short_url = $1", u)
	var total int = 0
	if err := row.Scan(&total); err != nil {
		return 0
	}
	return total
}

func AtoIWithDefault(s string, d int) int {
	if s == "" {
		return d
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

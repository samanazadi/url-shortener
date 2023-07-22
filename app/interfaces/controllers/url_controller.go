package controllers

import (
	"github.com/samanazadi/url-shortener/app/entities"
	"github.com/samanazadi/url-shortener/app/usecases"
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

// GetOriginalURL retrieves the original input
func (u URLController) GetOriginalURL(p URLControllerInputPort) {
	url := p.Param("id")
	originalURL, err := u.urlUseCase.OriginalURL(url)
	if err != nil {
		p.OutputError(URLNotFound, err)
		return
	}
	p.Output(Show, originalURL)
}

const (
	// Show id a successful retrieval of an URL
	Show = iota
)

const (
	// URLNotFound is an unsuccessful retrieval of an URL
	URLNotFound = iota
)

// URLControllerInputPort will be injected by infrastructure layer
type URLControllerInputPort interface {
	Param(string) string
	Output(int, any)
	OutputError(int, error)
}

// SQLHandler will be injected by infrastructure layer
type SQLHandler interface {
	Exec(string, ...any) (Result, error)
	QueryRow(string, ...any) Row
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

// URLControllerRepository is an implementation of usecases.URLRepository
type URLControllerRepository struct {
	SQLHandler SQLHandler
}

// FindURL queries the database for specified URL
func (r URLControllerRepository) FindURL(u string) (entities.URL, error) {
	row := r.SQLHandler.QueryRow("SELECT * FROM urls WHERE url = $1", u)
	var url, originalURL string
	err := row.Scan(&url, &originalURL)
	if err != nil {
		return entities.URL{}, err
	}
	return entities.URL{URL: url, OriginalURL: originalURL}, nil
}

package router

import (
	"github.com/samanazadi/url-shortener/internal/config"
	"github.com/samanazadi/url-shortener/internal/infrastructure/database"
	"github.com/samanazadi/url-shortener/pkg/entities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samanazadi/url-shortener/internal/adapters/controllers"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router/json"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gin.Engine, *gorm.DB, error) {
	router := gin.Default()
	db, err := database.NewGormDB(cfg)
	if err != nil {
		return nil, nil, err
	}
	urlController := controllers.NewURLController(db)

	router.GET("/u/:id", func(c *gin.Context) {
		urlController.GetDetails(c, WebURLControllerInputPort{c: c}, cfg)
	})
	router.POST("/u", func(c *gin.Context) {
		urlController.CreateShortLink(c, WebURLControllerInputPort{c: c}, cfg)
	})
	router.GET("/:id", func(c *gin.Context) {
		urlController.RedirectToOriginalURL(c, WebURLControllerInputPort{c: c})
	})

	return router, db, nil
}

// WebURLControllerInputPort implements controllers.URLControllerInputPort
type WebURLControllerInputPort struct {
	c *gin.Context
}

// Param retrieves URL parameter p
func (w WebURLControllerInputPort) Param(param string) string {
	if p := w.c.Param(param); p != "" {
		return p
	}
	return w.c.Query(param)

}

// OutputVisitDetails returns result JSON to client
func (w WebURLControllerInputPort) OutputVisitDetails(u string, v []entities.VisitDetail, total int) {
	vds := make([]json.VisitDetail, 0)
	for _, vd := range v {
		vds = append(vds, json.VisitDetail{
			IP:        vd.IP,
			Time:      vd.Time,
			UserAgent: vd.UserAgent,
		})
	}

	res := json.SuccessRetrieval{
		Message:      "Successful",
		OriginalURL:  u,
		Total:        total,
		VisitDetails: vds,
	}
	w.c.IndentedJSON(http.StatusOK, res)
}

func (w WebURLControllerInputPort) OutputShortURL(shortURL string) {
	res := json.SuccessShortURLCreated{
		Message:  "Successful",
		ShortURL: shortURL,
	}
	w.c.IndentedJSON(http.StatusCreated, res)
}

func (w WebURLControllerInputPort) Redirect(u string) {
	w.c.Redirect(http.StatusFound, u)
}

// OutputError returns error JSON to client
func (w WebURLControllerInputPort) OutputError(op int, err error) {
	switch op {
	case controllers.URLNotFound:
		j := json.MessageError{
			Message: "Unsuccessful",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusNotFound, j)
	case controllers.RedirectToHomePage:
		w.c.Redirect(http.StatusFound, "/")
	case controllers.CannotCreateShortLink:
		j := json.MessageError{
			Message: "Unsuccessful",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusInternalServerError, j)
	case controllers.BadRequest:
		j := json.MessageError{
			Message: "Bad Request",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusBadRequest, j)
	}
}

func (w WebURLControllerInputPort) GetVisitDetail() entities.VisitDetail {
	vd := entities.VisitDetail{}
	vd.IP = w.c.ClientIP()
	vd.Time = time.Now()
	vd.UserAgent = w.c.GetHeader("User-Agent")
	return vd
}

func (w WebURLControllerInputPort) GetCreateShortURLRequest() (string, error) {
	var reqBody json.CreateShortURLRequestBody
	err := w.c.BindJSON(&reqBody)
	return reqBody.URL, err
}

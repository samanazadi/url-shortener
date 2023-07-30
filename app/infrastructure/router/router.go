package router

import (
	"github.com/samanazadi/url-shortener/app/entities"
	"github.com/samanazadi/url-shortener/app/infrastructure/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samanazadi/url-shortener/app/adapters/controllers"
	"github.com/samanazadi/url-shortener/app/infrastructure/database/postgres"
	"github.com/samanazadi/url-shortener/app/infrastructure/router/JSON"
)

// Router is main gin router
var Router *gin.Engine

func init() {
	Router = gin.Default()
	urlController := controllers.NewURLController(postgres.NewSQLHandler())

	Router.GET("/u/:id", func(c *gin.Context) {
		urlController.GetDetails(WebURLControllerInputPort{c: c})
	})
	Router.POST("/u", func(c *gin.Context) {
		urlController.CreateShortLink(WebURLControllerInputPort{c: c})
	})
	Router.GET("/:id", func(c *gin.Context) {
		urlController.RedirectToOriginalURL(WebURLControllerInputPort{c: c})
	})
}

// WebURLControllerInputPort implements controllers.URLControllerInputPort
type WebURLControllerInputPort struct {
	c *gin.Context
}

func (w WebURLControllerInputPort) GetMachineID() uint16 {
	return config.GetUint16("machineid")
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
	vds := make([]JSON.VisitDetail, 0)
	for _, vd := range v {
		vds = append(vds, JSON.VisitDetail{
			IP:        vd.IP,
			Time:      vd.Time,
			UserAgent: vd.UserAgent,
		})
	}

	res := JSON.SuccessRetrieval{
		Message:      "Successful",
		OriginalURL:  u,
		Total:        total,
		VisitDetails: vds,
	}
	w.c.IndentedJSON(http.StatusOK, res)
}

func (w WebURLControllerInputPort) OutputShortURL(shortURL string) {
	res := JSON.SuccessShortURLCreated{
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
		j := JSON.MessageError{
			Message: "Unsuccessful",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusNotFound, j)
	case controllers.RedirectToHomePage:
		w.c.Redirect(http.StatusFound, "/")
	case controllers.CannotCreateShortLink:
		j := JSON.MessageError{
			Message: "Unsuccessful",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusInternalServerError, j)
	case controllers.BadRequest:
		j := JSON.MessageError{
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
	var reqBody JSON.CreateShortURLRequestBody
	err := w.c.BindJSON(&reqBody)
	return reqBody.URL, err
}

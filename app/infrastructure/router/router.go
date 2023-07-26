package router

import (
	"github.com/samanazadi/url-shortener/app/entities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samanazadi/url-shortener/app/infrastructure/database/postgres"
	"github.com/samanazadi/url-shortener/app/infrastructure/router/JSON"
	"github.com/samanazadi/url-shortener/app/interfaces/controllers"
)

// Router is main gin router
var Router *gin.Engine

func init() {
	Router = gin.Default()
	urlController := controllers.NewURLController(postgres.NewSQLHandler())

	Router.GET("/u/:id", func(c *gin.Context) {
		urlController.GetOriginalURL(WebURLControllerInputPort{c: c})
	})
	Router.GET("/:id", func(c *gin.Context) {
		urlController.RedirectToOriginalURL(WebURLControllerInputPort{c: c})
	})
}

// WebURLControllerInputPort implements controllers.URLControllerInputPort
type WebURLControllerInputPort struct {
	c *gin.Context
}

// Param retrieves URL parameter p
func (w WebURLControllerInputPort) Param(p string) string {
	return w.c.Param(p)
}

// Output returns result JSON to client
func (w WebURLControllerInputPort) Output(op int, res any) {
	switch op {
	case controllers.Show:
		j := JSON.SuccessRetrieval{
			Message:     "Successful",
			OriginalURL: res.(string)}
		w.c.IndentedJSON(http.StatusOK, j)
	case controllers.Redirect:
		w.c.Redirect(http.StatusFound, res.(string))
	}
}

// OutputError returns error JSON to client
func (w WebURLControllerInputPort) OutputError(op int, err error) {
	switch op {
	case controllers.URLNotFound:
		j := JSON.UnsuccessfulRetrieval{
			Message: "Unsuccessful",
			Error:   err.Error(),
		}
		w.c.IndentedJSON(http.StatusNotFound, j)
	case controllers.RedirectToHomePage:
		w.c.Redirect(http.StatusFound, "/")
	}
}

func (w WebURLControllerInputPort) GetVisitDetail() entities.VisitDetail {
	vd := entities.VisitDetail{}
	vd.IP = w.c.ClientIP()
	vd.Time = time.Now()
	vd.UserAgent = w.c.GetHeader("User-Agent")
	return vd
}

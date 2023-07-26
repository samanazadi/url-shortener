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
		urlController.GetDetails(WebURLControllerInputPort{c: c})
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
func (w WebURLControllerInputPort) Param(param string) string {
	if p := w.c.Param(param); p != "" {
		return p
	}
	return w.c.Query(param)

}

// Output returns result JSON to client
func (w WebURLControllerInputPort) Output(u string, v []entities.VisitDetail, total int) {
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

func (w WebURLControllerInputPort) Redirect(u string) {
	w.c.Redirect(http.StatusFound, u)
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

package router

import "github.com/gin-gonic/gin"

// Router is main gin router
var Router *gin.Engine

func init() {
	Router := gin.Default()

	Router.GET("/get", func(c *gin.Context) {})
}

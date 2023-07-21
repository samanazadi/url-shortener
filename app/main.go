package main

import (
	"github.com/samanazadi/url-shortener/app/infrastructure/config"
	"github.com/samanazadi/url-shortener/app/infrastructure/router"
)

func main() {
	router.Router.Run(config.GetString("server") + ":" + config.GetString("port"))
}

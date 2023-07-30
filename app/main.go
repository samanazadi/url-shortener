package main

import (
	"github.com/samanazadi/url-shortener/app/infrastructure"
	"github.com/samanazadi/url-shortener/app/infrastructure/router"
	"github.com/samanazadi/url-shortener/app/utilities"
)

func main() {
	defer utilities.Logger.Sync()
	err := router.Router.Run(
		infrastructure.Config.GetString("server") + ":" + infrastructure.Config.GetString("port"))
	if err != nil {
		utilities.Logger.Panic(err.Error())
	}
}

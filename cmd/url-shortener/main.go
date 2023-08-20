package main

import (
	"github.com/samanazadi/url-shortener/app/infrastructure/router"
	"github.com/samanazadi/url-shortener/configs"
	"github.com/samanazadi/url-shortener/internal/utilities"
)

func main() {
	defer utilities.Logger.Sync()
	err := router.Router.Run(
		configs.Config.GetString("server") + ":" + configs.Config.GetString("port"))
	if err != nil {
		utilities.Logger.Panic(err.Error())
	}
}

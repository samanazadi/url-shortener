package main

import (
	"github.com/samanazadi/url-shortener/configs"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router"
	"github.com/samanazadi/url-shortener/internal/utilities/logging"
)

func main() {
	// logging
	if err := logging.Init(); err != nil {
		panic(err)
	}
	defer func() {
		if err := logging.Logger.Sync(); err != nil {
			panic(err)
		}
	}()

	// router
	if err := router.Router.Run(configs.Config.GetString("server") + ":" + configs.Config.GetString("port")); err != nil {
		logging.Logger.Panic(err.Error())
	}
}

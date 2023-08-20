package main

import (
	"github.com/samanazadi/url-shortener/configs"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router"
	"github.com/samanazadi/url-shortener/internal/logging"
	"github.com/samanazadi/url-shortener/internal/usecases/base62"
)

func main() {
	// config
	configs.Init()

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
	router.Init()
	if err := router.Router.Run(configs.Config.GetString("server") + ":" + configs.Config.GetString("port")); err != nil {
		logging.Logger.Panic(err.Error())
	}

	// base62
	base62.Init()
}

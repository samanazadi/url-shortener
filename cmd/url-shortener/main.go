package main

import (
	"flag"
	"github.com/samanazadi/url-shortener/configs"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router"
	"github.com/samanazadi/url-shortener/internal/logging"
	"github.com/samanazadi/url-shortener/internal/usecases/base62"
)

func main() {
	// command line flags
	var cfgPath string //config path
	flag.StringVar(&cfgPath, "c", ".env", "config path")
	flag.Parse()

	// config
	if err := configs.Init(cfgPath); err != nil {
		panic(err)
	}

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
	if err := router.Init(); err != nil {
		logging.Logger.Panic(err.Error())
	}
	if err := router.Router.Run(configs.Config.GetString("server") + ":" + configs.Config.GetString("port")); err != nil {
		logging.Logger.Panic(err.Error())
	}

	// base62
	base62.Init()
}

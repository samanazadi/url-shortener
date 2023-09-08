package main

import (
	"fmt"
	"github.com/samanazadi/url-shortener/internal/config"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router"
	"github.com/samanazadi/url-shortener/pkg/base62"
	"github.com/samanazadi/url-shortener/pkg/logging"
	"github.com/spf13/pflag"
)

func main() {
	// command line flags
	cfgPath := pflag.StringP("config", "c", ".env", "config file path")
	pflag.Parse()
	fmt.Printf("config file: %s", *cfgPath)

	// config
	cfg, err := config.New(*cfgPath)
	if err != nil {
		panic(err)
	}

	// logging
	if err := logging.Init(cfg); err != nil {
		panic(err)
	}
	defer func() {
		if err := logging.Logger.Sync(); err != nil {
			panic(err)
		}
	}()
	logging.Logger.Info("logger started")

	// router
	if err := router.Init(cfg); err != nil {
		logging.Logger.Panic(err.Error())
	}
	logging.Logger.Info("router started")

	// base62
	base62.Init()
	logging.Logger.Info("algorithms initialized")
}

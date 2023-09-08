package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/samanazadi/url-shortener/internal/config"
	"github.com/samanazadi/url-shortener/internal/infrastructure/router"
	"github.com/samanazadi/url-shortener/pkg/base62"
	"github.com/samanazadi/url-shortener/pkg/logging"
	"github.com/spf13/pflag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// command line flags
	cfgPath := pflag.StringP("config", "c", ".env", "config file path")
	pflag.Parse()
	fmt.Printf("config file: %s\n", *cfgPath)

	// config
	cfg, err := config.New(*cfgPath)
	if err != nil {
		panic(err)
	}

	// logging
	if err := logging.Init(cfg.Development); err != nil {
		panic(err)
	}
	defer logging.Logger.Sync()
	logging.Logger.Info("logger started")

	// base62
	base62.Init()
	logging.Logger.Info("algorithms initialized")

	// router
	rtr, err := router.New(cfg)
	if err != nil {
		logging.Logger.Panic(err.Error())
	}
	logging.Logger.Info("router created")

	// server
	server := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: rtr,
	}

	go func() {
		logging.Logger.Info("server started at port: " + cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Logger.Panic("cannot start server", "error", err)
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logging.Logger.Info("starting server shutdown ...")

	if err := server.Shutdown(context.Background()); err != nil {
		logging.Logger.Panic("server Shutdown", "error", err)
	}
	logging.Logger.Info("server stopped")
}

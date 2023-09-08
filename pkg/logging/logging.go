package logging

import (
	"github.com/samanazadi/url-shortener/internal/config"
	"go.uber.org/zap"
)

var Logger logger

func Init(cfg *config.Config) error {
	var (
		l   *zap.Logger
		err error
	)
	if cfg.Development {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	Logger = zapLogger{
		logger: l.Sugar(),
	}
	return nil
}

type logger interface {
	Debug(msg string, p ...any)
	Info(msg string, p ...any)
	Warn(msg string, p ...any)
	Error(msg string, p ...any)
	DPanic(msg string, p ...any) // panics in development mode
	Panic(msg string, p ...any)  // logs then panics
	Fatal(msg string, p ...any)  // logs then calls os.Exit(1)
	Sync() error
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func (l zapLogger) Debug(msg string, p ...any) {
	l.logger.Debugw(msg, p...)
}

func (l zapLogger) Info(msg string, p ...any) {
	l.logger.Infow(msg, p...)
}

func (l zapLogger) Warn(msg string, p ...any) {
	l.logger.Warnw(msg, p...)
}

func (l zapLogger) Error(msg string, p ...any) {
	l.logger.Errorw(msg, p...)
}

func (l zapLogger) DPanic(msg string, p ...any) {
	l.logger.DPanicw(msg, p...)
}

func (l zapLogger) Panic(msg string, p ...any) {
	l.logger.Panicw(msg, p...)
}

func (l zapLogger) Fatal(msg string, p ...any) {
	l.logger.Fatalw(msg, p...)
}

func (l zapLogger) Sync() error {
	return l.logger.Sync()
}

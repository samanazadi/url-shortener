package logging

import (
	"go.uber.org/zap"
)

type Options struct {
	Development bool
}

var Logger Log

// Init initializes the package
func Init(opt Options) error {
	var (
		l   *zap.Logger
		err error
	)
	if opt.Development {
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

// Log is the logging interface which must be implemented by loggers.
type Log interface {
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

// Debug prints a debug level log
func (l zapLogger) Debug(msg string, p ...any) {
	l.logger.Debugw(msg, p...)
}

// Info prints a info level log
func (l zapLogger) Info(msg string, p ...any) {
	l.logger.Infow(msg, p...)
}

// Warn prints a warn level log
func (l zapLogger) Warn(msg string, p ...any) {
	l.logger.Warnw(msg, p...)
}

// Error prints a error level log
func (l zapLogger) Error(msg string, p ...any) {
	l.logger.Errorw(msg, p...)
}

// DPanic prints a message and then only in development panics.
func (l zapLogger) DPanic(msg string, p ...any) {
	l.logger.DPanicw(msg, p...)
}

// Panic prints a message and then panics.
func (l zapLogger) Panic(msg string, p ...any) {
	l.logger.Panicw(msg, p...)
}

// Fatal prints a message and then calls os.Exit.
func (l zapLogger) Fatal(msg string, p ...any) {
	l.logger.Fatalw(msg, p...)
}

// Sync flushes any buffered log entries.
func (l zapLogger) Sync() error {
	return l.logger.Sync()
}

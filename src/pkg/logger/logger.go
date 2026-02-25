package logger

import (
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/myproject/api/config"
)

type ILogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	log *slog.Logger
}

func NewLogger(cfg config.Config) *httplog.Logger {
	log := httplog.NewLogger("myapp-api", httplog.Options{
		JSON:             cfg.LogJSON,
		LogLevel:         cfg.LogLevel,
		Concise:          false,
		RequestHeaders:   false,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})
	return log
}

func (l *logger) Debug(message string, args ...any) {
	l.log.Debug(message, args...)
}

func (l *logger) Info(message string, args ...any) {
	l.log.Info(message, args...)
}

func (l *logger) Warn(message string, args ...any) {
	l.log.Warn(message, args...)
}

func (l *logger) Error(message string, args ...any) {
	l.log.Error(message, args...)
}

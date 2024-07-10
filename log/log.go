package log

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

var logger *slog.Logger

func init() {
	logger = slog.New(devslog.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

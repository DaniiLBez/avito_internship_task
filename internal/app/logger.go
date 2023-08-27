package app

import (
	"log/slog"
	"os"
)

const (
	InfoLevel    = "INFO"
	DebugLevel   = "DEBUG"
	WarningLevel = "WARN"
	ErrorLevel   = "ERROR"
)

func SetLogger(level string) (logger *slog.Logger) {
	switch level {
	case InfoLevel:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case DebugLevel:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case WarningLevel:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case ErrorLevel:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}
	return
}

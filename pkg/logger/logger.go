package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var once sync.Once

func InitLogger(logLevel, logFormat string) error {
	var err error

	once.Do(func() {
		var handler slog.Handler

		level := stringToLogLevel(logLevel)

		switch logFormat {
		case "json":
			handler = newJSONLogHandler(os.Stdout, level)
		case "text":
			handler = newColorLogHandler(os.Stdout, level)
		default:
			err = ErrInvalidLogFormat

			return
		}

		logger := slog.New(handler)
		slog.SetDefault(logger)
	})

	return err
}

func stringToLogLevel(level string) slog.Level {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

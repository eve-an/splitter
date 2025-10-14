package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func parseLogLevel(logLevel string) (slog.Level, error) {
	switch logLevel {
	case "debug":
		return slog.LevelDebug, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "info":
		return slog.LevelInfo, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", logLevel)
	}
}

func NewLogger(logLevel string) (*slog.Logger, error) {
	slogLevel, err := parseLogLevel(logLevel)
	if err != nil {
		return slog.Default(), err
	}

	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slogLevel,
			TimeFormat: time.Kitchen,
		}),
	)

	slog.SetDefault(logger)

	return logger, nil
}

package logger

import (
	"log/slog"
	"os"
	"strings"
)

// SetLogger initializes the global structured logger.
func SetLogger() {
	level := getLogLevel()

	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug, // Show source file/line only in debug mode
	}

	handler := slog.NewJSONHandler(os.Stderr, options)
	logger := slog.New(handler)

	slog.SetDefault(logger)
	slog.Info("logger initialized", "level", level.String())
}

// getLogLevel parses the LOG_LEVEL environment variable.
func getLogLevel() slog.Level {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to INFO
	}
}

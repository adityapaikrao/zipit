package logger

import (
	"log/slog"
	"os"
)

func SetLogger() {
	options := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stderr, options)
	logger := slog.New(handler)

	slog.SetDefault(logger)

}

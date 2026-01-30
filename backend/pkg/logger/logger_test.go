package logger

import (
	"log/slog"
	"os"
	"testing"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected slog.Level
	}{
		{"Default should be INFO", "", slog.LevelInfo},
		{"Case Insensitive Debug", "debug", slog.LevelDebug},
		{"valid WARN", "WARN", slog.LevelWarn},
		{"Invalid values default o INFO", "INVALID_VALUE", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.envValue)

			logLevel := getLogLevel()
			if tt.expected != logLevel {
				t.Errorf("%v Log Level failed. Expected %v, got %v", tt.name, tt.expected, logLevel)
			}
		})
	}
}

func TestSetLogger(t *testing.T) {
	SetLogger()
}

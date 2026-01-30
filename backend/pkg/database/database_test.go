package database

import (
	"context"
	"os"
	"testing"
	"time"
	"zipit/pkg/config"
)

var keyMap = [][]string{
	{"DB_PORT", "5432"},
	{"DB_USER", "test_user"},
	{"DB_PASSWORD", "test_password"},
	{"DB_NAME", "test_db"},
	{"DB_HOST", "localhost"},
}

func TestDefaultPoolConfig(t *testing.T) {
	cfg := DefaultPoolConfig()

	if cfg.MaxOpenConns != 25 {
		t.Errorf("expected MaxOpenConns 25, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 25 {
		t.Errorf("expected MaxIdleConns 25, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("expected ConnMaxLifetime 5m, got %v", cfg.ConnMaxLifetime)
	}
}

func TestDatabaseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database integration tests in short mode")
	}

	// Set test config keys
	os.Clearenv()
	for _, configPairs := range keyMap {
		os.Setenv(configPairs[0], configPairs[1])
	}

	testConfig, err := config.NewDBConfig()
	if err != nil {
		t.Fatalf("Error with package config: %v", err)
	}

	db, err := NewDatabase(testConfig)
	if err != nil {
		t.Fatalf("Failed setting up new database connection: %v", err)
	}
	defer db.Close()

	// Test Health check
	t.Run("Health", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.Health(ctx); err != nil {
			t.Errorf("Health() check failed: %v", err)
		}
	})

	// Test Stats
	t.Run("Stats", func(t *testing.T) {
		stats := db.Stats()
		if stats.MaxOpenConnections != 25 {
			t.Errorf("expected 25 max open connections in stats, got %d", stats.MaxOpenConnections)
		}
	})
}

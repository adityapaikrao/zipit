package repository

import (
	"context"
	"os"
	"testing"
	"time"
	"zipit/pkg/config"
	"zipit/pkg/database"
)

var keyMap = [][]string{
	{"DB_PORT", "5432"},
	{"DB_USER", "test_user"},
	{"DB_PASSWORD", "test_password"},
	{"DB_NAME", "test_db"},
	{"DB_HOST", "localhost"},
}

func TestPostgresRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Postgres test in short mode")
	}

	// setup test env variables
	os.Clearenv()
	for _, envPair := range keyMap {
		os.Setenv(envPair[0], envPair[1])
	}
	config, err := config.NewDBConfig()
	if err != nil {
		t.Fatalf("failed to setup test env %v", err)
	}
	db, err := database.NewDatabase(config)
	if err != nil {
		t.Fatalf("failed to setup test db %v", err)
	}
	defer db.Close()

	// Ensure table exists before running tests
	setupSchema(t, db)

	repo := NewPostgresRepository(db)

	// Sub-test for a complete shorten and resolve cycle
	t.Run("Full Cycle", func(t *testing.T) {
		cleanup(t, db) // Ensure we start with a clean table
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		longURL := "https://example.com/some-very-long-link"

		// 1. Create entry
		id, err := repo.CreateURL(ctx, longURL)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// 2. Mock a short code and save it
		code := "abc123"
		err = repo.SetShortCode(ctx, id, code)
		if err != nil {
			t.Fatalf("SetShortCode failed: %v", err)
		}

		// 3. Verify it's there via URLExists
		exists, existingID, err := repo.URLExists(ctx, longURL)
		if err != nil || !exists || existingID != id {
			t.Errorf("URLExists failed: exists=%v, expectedID=%d, gotID=%d, err=%v", exists, id, existingID, err)
		}

		// 4. Resolve it back from short code to long URL
		resURL, err := repo.GetURLByShortCode(ctx, code)
		if err != nil || resURL != longURL {
			t.Errorf("GetURLByShortCode failed: expected %s, got %s, err=%v", longURL, resURL, err)
		}
	})
}

// setupSchema ensures the database table exists.
// In a larger project, this would be handled by migration files.
func setupSchema(t *testing.T, db *database.Database) {
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id BIGSERIAL PRIMARY KEY,
		long_url TEXT NOT NULL,
		short_code VARCHAR(12) UNIQUE,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Conn.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
}

// cleanup clears the urls table and resets the auto-increment ID counter.
func cleanup(t *testing.T, db *database.Database) {
	_, err := db.Conn.Exec("TRUNCATE TABLE urls RESTART IDENTITY")
	if err != nil {
		t.Fatalf("failed to cleanup database: %v", err)
	}
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"zipit/pkg/database"
)

// postgresRepository implements the URLRepository interface using a PostgreSQL database.
type postgresRepository struct {
	db *database.Database
}

// CreateURL inserts a new long URL into the database and returns its auto-generated ID.
// We do this first because we need the ID to generate a unique short code.
func (pgRepo *postgresRepository) CreateURL(ctx context.Context, longURL string) (int64, error) {
	var id int64
	// PostgreSQL "RETURNING id" saves us an extra query to find out what ID was assigned.
	query := `INSERT INTO urls (long_url) VALUES ($1) RETURNING id`

	// QueryRowContext is for queries that return exactly one row.
	err := pgRepo.db.Conn.QueryRowContext(ctx, query, longURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert long URL: %w", err)
	}

	return id, nil
}

// GetURLByShortCode retrieves the original long URL for a given short code alias.
func (pgRepo *postgresRepository) GetURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	var longURL string
	query := "SELECT long_url FROM urls WHERE short_code = $1"

	err := pgRepo.db.Conn.QueryRowContext(ctx, query, shortCode).Scan(&longURL)
	if err != nil {
		// sql.ErrNoRows means the query was valid but no matching record exists.
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short code %q does not exist: %w", shortCode, sql.ErrNoRows)
		}
		return "", fmt.Errorf("failed to retrieve URL: %w", err)
	}
	return longURL, nil
}

// SetShortCode updates an existing URL record with its calculated short code.
func (pgRepo *postgresRepository) SetShortCode(ctx context.Context, id int64, shortCode string) error {
	query := "UPDATE urls SET short_code = $1 WHERE id = $2"

	// ExecContext is for statements that modify data but don't return rows.
	results, err := pgRepo.db.Conn.ExecContext(ctx, query, shortCode, id)
	if err != nil {
		return fmt.Errorf("failed to update short code: %w", err)
	}

	// Verify that the update actually happened (ID might not exist).
	rows, _ := results.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no URL found with ID %d", id)
	}
	return nil
}

// URLExists checks if a long URL has already been shortened to avoid duplicates.
func (pgRepo *postgresRepository) URLExists(ctx context.Context, longURL string) (bool, int64, error) {
	var id int64
	query := "SELECT id FROM urls WHERE long_url = $1"

	err := pgRepo.db.Conn.QueryRowContext(ctx, query, longURL).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return 0 as the ID when not found (standard Go idiom for missing values).
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("failed to check url existence: %w", err)
	}
	return true, id, nil
}

// NewPostgresRepository constructor for creating a new repository instance.
func NewPostgresRepository(db *database.Database) URLRepository {
	return &postgresRepository{
		db: db,
	}
}

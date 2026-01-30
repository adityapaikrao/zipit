package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"zipit/pkg/config"

	_ "github.com/lib/pq"
)

// Database wraps the sql.DB connection.
type Database struct {
	Conn *sql.DB
}

// PoolConfig holds database connection pool settings.
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultPoolConfig returns sensible default connection pool settings.
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// NewDatabase creates a new database connection with default pool settings.
func NewDatabase(cfg *config.DBConfig) (*Database, error) {
	return NewDatabaseWithPool(cfg, DefaultPoolConfig())
}

// NewDatabaseWithPool creates a new database connection with custom pool settings.
func NewDatabaseWithPool(cfg *config.DBConfig, pool PoolConfig) (*Database, error) {
	db, err := sql.Open(cfg.Driver, cfg.ConnectionURL)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(pool.MaxOpenConns)
	db.SetMaxIdleConns(pool.MaxIdleConns)
	db.SetConnMaxLifetime(pool.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("no response received from database: %w", err)
	}

	slog.Info("Database connection established successfully")
	return &Database{
		Conn: db,
	}, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	slog.Info("Closing database connection")
	return d.Conn.Close()
}

// Health checks if the database is reachable.
func (d *Database) Health(ctx context.Context) error {
	return d.Conn.PingContext(ctx)
}

// Stats returns database connection pool statistics.
func (d *Database) Stats() sql.DBStats {
	return d.Conn.Stats()
}

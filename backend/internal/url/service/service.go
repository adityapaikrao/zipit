package service

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("url not found")
	ErrInvalidURL    = errors.New("invalid url")
	ErrDatabaseRead  = errors.New("error reading from database")
	ErrDatabaseWrite = errors.New("error writing to database")
)

// URLService defines the interface for URL shortening operations.
// It provides methods to create shortened URLs and retrieve original URLs
// from their short code representations.
//
// ShortenURL takes a long URL and returns a shortened version of it.
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - longURL: The original URL to be shortened
//
// Returns:
//   - string: The generated short code or shortened URL
//   - error: An error if the URL is invalid or the operation fails
//
// GetLongURL retrieves the original URL associated with a given short code.
// Parameters:
//   - ctx: Context for request cancellation and timeouts
//   - shortCode: The short code identifier for the URL
//
// Returns:
//   - string: The original long URL
//   - error: An error if the short code is not found or the operation fails
type URLService interface {
	ShortenURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
}

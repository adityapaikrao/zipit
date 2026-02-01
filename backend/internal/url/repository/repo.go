package repository

import "context"

// URLRepository provides an abstraction for persisting and retrieving URL
// records used by the URL-shortening service. Implementations handle storage
// details (e.g., database access), transactional concerns, and concurrency.
//
// CreateURL creates a new record for the provided longURL and returns the
// generated record ID, or an error if creation fails.
//
// GetURLByShortCode returns the original long URL associated with the given
// shortCode. If no record is found, the implementation should return a
// non-nil error describing the condition.
//
// URLExists checks whether the given longURL already exists in storage.
// It returns a boolean indicating existence, the existing record ID (0 if not
// found), and an error if the existence check could not be performed.
//
// SetShortCode associates the provided shortCode with an existing record
// identified by id. It returns an error if the update fails or if the id does
// not correspond to an existing record.
type URLRepository interface {
	CreateURL(ctx context.Context, longURL string) (int64, error)
	GetURLByShortCode(ctx context.Context, shortCode string) (string, error)
	URLExists(ctx context.Context, longURL string) (bool, int64, error)
	SetShortCode(ctx context.Context, id int64, shortCode string) error
}

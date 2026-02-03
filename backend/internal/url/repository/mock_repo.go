package repository

import (
	"context"
	"fmt"
)

// MockRepo for testing service logic
type MockRepo struct {
	URLExistsFunc         func(ctx context.Context, longURL string) (bool, int64, error)
	CreateURLFunc         func(ctx context.Context, longURL string) (int64, error)
	GetURLByShortCodeFunc func(ctx context.Context, shortCode string) (string, error)
	SetShortCodeFunc      func(ctx context.Context, id int64, shortCode string) error
}

// CreateURL implements [URLRepository].
func (m *MockRepo) CreateURL(ctx context.Context, longURL string) (int64, error) {
	if m.CreateURLFunc != nil {
		return m.CreateURLFunc(ctx, longURL)
	}
	return 0, fmt.Errorf("some error creating long url")
}

// GetURLByShortCode implements [URLRepository].
func (m *MockRepo) GetURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	if m.GetURLByShortCodeFunc != nil {
		return m.GetURLByShortCodeFunc(ctx, shortCode)
	}
	return "", fmt.Errorf("some error fetching long url by short code") // database error
}

// SetShortCode implements [URLRepository].
func (m *MockRepo) SetShortCode(ctx context.Context, id int64, shortCode string) error {
	if m.SetShortCodeFunc != nil {
		return m.SetShortCodeFunc(ctx, id, shortCode)
	}
	return fmt.Errorf("some error setting short code") // database error
}

// URLExists implements [URLRepository].
func (m *MockRepo) URLExists(ctx context.Context, longURL string) (bool, int64, error) {
	if m.URLExistsFunc != nil {
		return m.URLExistsFunc(ctx, longURL)
	}
	return false, 0, fmt.Errorf("some error checking long url in database") // database error
}

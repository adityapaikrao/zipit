package repository

import "context"

type URLRepository interface {
	CreateURL(ctx context.Context, longURL string) (int64, error)
	GetURLByShortCode(ctx context.Context, shortCode string) (string, error)
	URLExists(ctx context.Context, longURL string) (bool, int64, error)
	SetShortCode(ctx context.Context, id int64, shortCode string) error
}

package service

import "context"

type URLService interface {
	ShortenURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
}

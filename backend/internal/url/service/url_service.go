package service

import (
	"context"
	"zipit/internal/url/repository"
	"zipit/pkg/shortener"
)

type urlSvc struct {
	repo      repository.URLRepository
	shortener shortener.Shortener
}

// GetLongURL implements [URLService].
func (u *urlSvc) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	panic("unimplemented")
}

// ShortenURL implements [URLService].
func (u *urlSvc) ShortenURL(ctx context.Context, longURL string) (string, error) {
	panic("unimplemented")
}

func NewUrlSvc(repo repository.URLRepository, shortnener shortener.Shortener) URLService {
	return &urlSvc{
		repo:      repo,
		shortener: shortnener,
	}
}

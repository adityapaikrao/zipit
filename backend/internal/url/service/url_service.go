package service

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"zipit/internal/url/repository"
	"zipit/pkg/shortener"
)

type urlSvc struct {
	repo      repository.URLRepository
	shortener shortener.Shortener
}

// GetLongURL implements [URLService].
func (svc *urlSvc) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	longURL, err := svc.repo.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", ErrDatabaseRead
	}
	return longURL, nil
}

// ShortenURL implements [URLService].
func (svc *urlSvc) ShortenURL(ctx context.Context, longURL string) (string, error) {
	if !isValidURL(longURL) {
		return "", ErrInvalidURL
	}

	urlExists, id, err := svc.repo.URLExists(ctx, longURL)
	if err != nil {
		return "", ErrDatabaseRead
	}

	// create new entry
	if urlExists {
		return svc.shortener.Encode(id), nil
	}

	id, err = svc.repo.CreateURL(ctx, longURL)
	if err != nil {
		return "", ErrDatabaseWrite
	}
	shortCode := svc.shortener.Encode(id)
	err = svc.repo.SetShortCode(ctx, id, shortCode)
	if err != nil {
		return "", ErrDatabaseWrite
	}

	return shortCode, nil
}

func NewUrlSvc(repo repository.URLRepository, shortener shortener.Shortener) URLService {
	return &urlSvc{
		repo:      repo,
		shortener: shortener,
	}
}

func isValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

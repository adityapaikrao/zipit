package service

import (
	"context"
	"database/sql"
	"errors"
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
	// TODO: implement URL validation
	// keeping it simple for now...
	if longURL == "" {
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

func NewUrlSvc(repo repository.URLRepository, shortnener shortener.Shortener) URLService {
	return &urlSvc{
		repo:      repo,
		shortener: shortnener,
	}
}

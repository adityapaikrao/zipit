package service

import (
	"context"
	"fmt"
	"zipit/internal/url/repository"
	"zipit/pkg/shortener"
)

type urlSvc struct {
	repo      repository.URLRepository
	shortener shortener.Shortener
}

// GetLongURL implements [URLService].
func (svc *urlSvc) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	return svc.repo.GetURLByShortCode(ctx, shortCode)
}

// ShortenURL implements [URLService].
func (svc *urlSvc) ShortenURL(ctx context.Context, longURL string) (string, error) {
	// TODO: implement URL validation
	// keeping it simple for now...
	if longURL == "" {
		return "", fmt.Errorf("Invalid URL: %v", longURL)
	}

	urlExists, id, err := svc.repo.URLExists(ctx, longURL)
	if err != nil {
		return "", fmt.Errorf("error fetching url from database: %w", err)
	}

	// create new entry
	if urlExists {
		return svc.shortener.Encode(id), nil
	}

	id, err = svc.repo.CreateURL(ctx, longURL)
	if err != nil {
		return "", fmt.Errorf("could not create new url: %w", err)
	}
	shortCode := svc.shortener.Encode(id)
	err = svc.repo.SetShortCode(ctx, id, shortCode)
	if err != nil {
		return "", fmt.Errorf("error writing shortCode for new url %w", err)
	}

	return shortCode, nil
}

func NewUrlSvc(repo repository.URLRepository, shortnener shortener.Shortener) URLService {
	return &urlSvc{
		repo:      repo,
		shortener: shortnener,
	}
}

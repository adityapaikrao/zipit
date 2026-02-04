package service

import (
	"context"
	"testing"
	"zipit/internal/url/repository"
	"zipit/pkg/shortener"
)

func TestUrlSvc_ShortenURL_Success_NewURL(t *testing.T) {
	// 1. Setup mocks
	mockRepo := &repository.MockRepo{
		URLExistsFunc: func(ctx context.Context, longURL string) (bool, int64, error) {
			return false, 0, nil // URL doesn't exist
		},
		CreateURLFunc: func(ctx context.Context, longURL string) (int64, error) {
			return 12345, nil // Return a fake database ID
		},
		SetShortCodeFunc: func(ctx context.Context, id int64, shortCode string) error {
			return nil // Simulate successful update
		},
	}

	// We'll use the real shortener since it's a pure function (no side effects)
	realShortener := shortener.NewBase62Shortener()
	svc := NewUrlSvc(mockRepo, realShortener)

	// 2. Define input
	longURL := "https://example.com"
	expectedCode := realShortener.Encode(12345) // Based on our mock ID

	// 3. Execute
	ctx := context.Background()
	code, err := svc.ShortenURL(ctx, longURL)

	// 4. Verification
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if code != expectedCode {
		t.Errorf("Expected code %s, got %s", expectedCode, code)
	}
}

func TestUrlSvc_ShortenURL_AlreadyExists(t *testing.T) {
	mockRepo := &repository.MockRepo{
		URLExistsFunc: func(ctx context.Context, longURL string) (bool, int64, error) {
			return true, 999, nil
		},
	}
	realShortener := shortener.NewBase62Shortener()
	svc := NewUrlSvc(mockRepo, realShortener)

	expectedCode := realShortener.Encode(999)
	code, err := svc.ShortenURL(context.Background(), "https://existing.com")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if code != expectedCode {
		t.Errorf("Expected code %s, got %s", expectedCode, code)
	}
}

func TestUrlSvc_ShortenURL_EmptyURL(t *testing.T) {
	svc := NewUrlSvc(&repository.MockRepo{}, shortener.NewBase62Shortener())
	_, err := svc.ShortenURL(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty URL, got nil")
	}
}

func TestUrlSvc_ShortenURL_URLExistsError(t *testing.T) {
	mockRepo := &repository.MockRepo{
		URLExistsFunc: func(ctx context.Context, longURL string) (bool, int64, error) {
			return false, 0, context.DeadlineExceeded
		},
	}
	svc := NewUrlSvc(mockRepo, shortener.NewBase62Shortener())
	_, err := svc.ShortenURL(context.Background(), "https://example.com")
	if err == nil {
		t.Error("Expected error from URLExists, got nil")
	}
}

func TestUrlSvc_ShortenURL_CreateURLError(t *testing.T) {
	mockRepo := &repository.MockRepo{
		URLExistsFunc: func(ctx context.Context, longURL string) (bool, int64, error) {
			return false, 0, nil
		},
		CreateURLFunc: func(ctx context.Context, longURL string) (int64, error) {
			return 0, context.DeadlineExceeded
		},
	}
	svc := NewUrlSvc(mockRepo, shortener.NewBase62Shortener())
	_, err := svc.ShortenURL(context.Background(), "https://example.com")
	if err == nil {
		t.Error("Expected error from CreateURL, got nil")
	}
}

func TestUrlSvc_ShortenURL_SetShortCodeError(t *testing.T) {
	mockRepo := &repository.MockRepo{
		URLExistsFunc: func(ctx context.Context, longURL string) (bool, int64, error) {
			return false, 0, nil
		},
		CreateURLFunc: func(ctx context.Context, longURL string) (int64, error) {
			return 123, nil
		},
		SetShortCodeFunc: func(ctx context.Context, id int64, shortCode string) error {
			return context.DeadlineExceeded
		},
	}
	svc := NewUrlSvc(mockRepo, shortener.NewBase62Shortener())
	_, err := svc.ShortenURL(context.Background(), "https://example.com")
	if err == nil {
		t.Error("Expected error from SetShortCode, got nil")
	}
}

func TestUrlSvc_GetLongURL_Success(t *testing.T) {
	expectedURL := "https://example.com"
	mockRepo := &repository.MockRepo{
		GetURLByShortCodeFunc: func(ctx context.Context, shortCode string) (string, error) {
			return expectedURL, nil
		},
	}
	svc := NewUrlSvc(mockRepo, shortener.NewBase62Shortener())
	url, err := svc.GetLongURL(context.Background(), "abc")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if url != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, url)
	}
}

func TestUrlSvc_GetLongURL_Error(t *testing.T) {
	mockRepo := &repository.MockRepo{
		GetURLByShortCodeFunc: func(ctx context.Context, shortCode string) (string, error) {
			return "", context.DeadlineExceeded
		},
	}
	svc := NewUrlSvc(mockRepo, shortener.NewBase62Shortener())
	_, err := svc.GetLongURL(context.Background(), "abc")
	if err == nil {
		t.Error("Expected error from GetURLByShortCode, got nil")
	}
}

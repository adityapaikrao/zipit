package grpc

import (
	"context"
	"errors"
	"testing"

	pb "zipit/gen/url"
	"zipit/internal/url/service"
)

type mockURLService struct {
	shortenURLFunc func(ctx context.Context, longURL string) (string, error)
	getLongURLFunc func(ctx context.Context, shortCode string) (string, error)
}

func (m *mockURLService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	return m.shortenURLFunc(ctx, longURL)
}

func (m *mockURLService) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	return m.getLongURLFunc(ctx, shortCode)
}

func TestPostURL(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.LongURL
		mockCode    string
		mockErr     error
		wantAlias   string
		wantErrCode string
	}{
		{
			name:      "Success",
			req:       &pb.LongURL{Url: "https://example.com"},
			mockCode:  "abcde",
			wantAlias: "abcde",
		},
		{
			name:        "Nil Request",
			req:         nil,
			wantErrCode: "InvalidArgument",
		},
		{
			name:        "Empty URL",
			req:         &pb.LongURL{Url: ""},
			wantErrCode: "InvalidArgument",
		},
		{
			name:        "Service Returns ErrInvalidURL",
			req:         &pb.LongURL{Url: "not-a-url"},
			mockErr:     service.ErrInvalidURL,
			wantErrCode: "InvalidArgument",
		},
		{
			name:        "Service Returns Internal Error",
			req:         &pb.LongURL{Url: "https://example.com"},
			mockErr:     errors.New("db down"),
			wantErrCode: "Internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockURLService{
				shortenURLFunc: func(ctx context.Context, longURL string) (string, error) {
					return tt.mockCode, tt.mockErr
				},
			}
			h := NewURLHandler(mock)
			resp, err := h.PostURL(context.Background(), tt.req)

			if tt.wantErrCode != "" {
				if err == nil {
					t.Fatalf("expected error with code %s, got nil", tt.wantErrCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.Alias != tt.wantAlias {
				t.Errorf("expected alias %s, got %s", tt.wantAlias, resp.Alias)
			}
		})
	}
}

func TestGetLongURL(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.ShortURL
		mockURL     string
		mockErr     error
		wantURL     string
		wantErrCode string
	}{
		{
			name:    "Success",
			req:     &pb.ShortURL{Alias: "abcde"},
			mockURL: "https://example.com",
			wantURL: "https://example.com",
		},
		{
			name:        "Nil Request",
			req:         nil,
			wantErrCode: "InvalidArgument",
		},
		{
			name:        "Empty Alias",
			req:         &pb.ShortURL{Alias: ""},
			wantErrCode: "InvalidArgument",
		},
		{
			name:        "Not Found",
			req:         &pb.ShortURL{Alias: "miss"},
			mockErr:     service.ErrNotFound,
			wantErrCode: "NotFound",
		},
		{
			name:        "Internal Error",
			req:         &pb.ShortURL{Alias: "abcde"},
			mockErr:     errors.New("db down"),
			wantErrCode: "Internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockURLService{
				getLongURLFunc: func(ctx context.Context, shortCode string) (string, error) {
					return tt.mockURL, tt.mockErr
				},
			}
			h := NewURLHandler(mock)
			resp, err := h.GetLongURL(context.Background(), tt.req)

			if tt.wantErrCode != "" {
				if err == nil {
					t.Fatalf("expected error with code %s, got nil", tt.wantErrCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.Url != tt.wantURL {
				t.Errorf("expected url %s, got %s", tt.wantURL, resp.Url)
			}
		})
	}
}

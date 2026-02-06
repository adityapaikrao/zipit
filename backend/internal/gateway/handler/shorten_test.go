package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "zipit/gen/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockURLServiceClient struct {
	pb.URLServiceClient
	postURLFunc    func(ctx context.Context, in *pb.LongURL, opts ...grpc.CallOption) (*pb.ShortURL, error)
	getLongURLFunc func(ctx context.Context, in *pb.ShortURL, opts ...grpc.CallOption) (*pb.LongURL, error)
}

func (m *mockURLServiceClient) PostURL(ctx context.Context, in *pb.LongURL, opts ...grpc.CallOption) (*pb.ShortURL, error) {
	return m.postURLFunc(ctx, in, opts...)
}

func (m *mockURLServiceClient) GetLongURL(ctx context.Context, in *pb.ShortURL, opts ...grpc.CallOption) (*pb.LongURL, error) {
	return m.getLongURLFunc(ctx, in, opts...)
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		mockResp       *pb.ShortURL
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			payload:        `{"long_url": "https://example.com"}`,
			mockResp:       &pb.ShortURL{Alias: "abcde"},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"short_code":"abcde"}`,
		},
		{
			name:           "Unknown Fields",
			payload:        `{"long_url": "https://example.com", "extra": "field"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON payload"}`,
		},
		{
			name:           "Invalid JSON",
			payload:        `{"long_url": "https://example.com",}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON payload"}`,
		},
		{
			name:           "Missing URL",
			payload:        `{"long_url": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"long_url is required"}`,
		},
		{
			name:           "Invalid URL (gRPC error)",
			payload:        `{"long_url": "invalid-url"}`,
			mockErr:        status.Error(codes.InvalidArgument, "invalid url"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid url"}`,
		},
		{
			name:           "Internal gRPC Error",
			payload:        `{"long_url": "https://example.com"}`,
			mockErr:        errors.New("some grpc error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to shorten url"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockURLServiceClient{
				postURLFunc: func(ctx context.Context, in *pb.LongURL, opts ...grpc.CallOption) (*pb.ShortURL, error) {
					return tt.mockResp, tt.mockErr
				},
			}
			h := NewGatewayHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.payload))
			rr := httptest.NewRecorder()

			h.ShortenURL(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

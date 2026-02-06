package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "zipit/gen/url"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		mockResp       *pb.LongURL
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			code:           "abcde",
			mockResp:       &pb.LongURL{Url: "https://example.com"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"long_url":"https://example.com"}`,
		},
		{
			name:           "Not Found",
			code:           "miss",
			mockErr:        status.Error(codes.NotFound, "not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"short url not found"}`,
		},
		{
			name:           "Internal gRPC Error",
			code:           "err",
			mockErr:        errors.New("some grpc error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to resolve url"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockURLServiceClient{
				getLongURLFunc: func(ctx context.Context, in *pb.ShortURL, opts ...grpc.CallOption) (*pb.LongURL, error) {
					return tt.mockResp, tt.mockErr
				},
			}
			h := NewGatewayHandler(mockSvc)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/%s", tt.code), nil)

			// Chi context logic to simulate {code} parameter
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("code", tt.code)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			h.ResolveURL(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

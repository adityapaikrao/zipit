package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "zipit/gen/url"
	"zipit/internal/gateway/handler"

	"google.golang.org/grpc"
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

func TestHealthEndpoint(t *testing.T) {
	mockSvc := &mockURLServiceClient{}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	expected := `{"status":"ok"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestRouteNotFound(t *testing.T) {
	mockSvc := &mockURLServiceClient{}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	mockSvc := &mockURLServiceClient{}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	req := httptest.NewRequest(http.MethodOptions, "/api/shorten", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 for CORS preflight, got %d", rr.Code)
	}

	if allow := rr.Header().Get("Access-Control-Allow-Origin"); allow != "*" {
		t.Errorf("expected Access-Control-Allow-Origin *, got %s", allow)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mockSvc := &mockURLServiceClient{}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	// DELETE is not allowed on /api/shorten
	req := httptest.NewRequest(http.MethodDelete, "/api/shorten", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestShortenRouteExists(t *testing.T) {
	mockSvc := &mockURLServiceClient{
		postURLFunc: func(ctx context.Context, in *pb.LongURL, opts ...grpc.CallOption) (*pb.ShortURL, error) {
			return &pb.ShortURL{Alias: "abc"}, nil
		},
	}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"long_url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}
}

func TestResolveRouteExists(t *testing.T) {
	mockSvc := &mockURLServiceClient{
		getLongURLFunc: func(ctx context.Context, in *pb.ShortURL, opts ...grpc.CallOption) (*pb.LongURL, error) {
			return &pb.LongURL{Url: "https://example.com"}, nil
		},
	}
	h := handler.NewGatewayHandler(mockSvc)
	r := New(h)

	req := httptest.NewRequest(http.MethodGet, "/api/abc", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

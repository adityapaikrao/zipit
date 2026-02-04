package grpc

import (
	"context"
	"errors"
	pb "zipit/gen/url"
	"zipit/internal/url/service"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type URLHandler struct {
	pb.UnimplementedURLServiceServer
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) PostURL(ctx context.Context, req *pb.LongURL) (*pb.ShortURL, error) {
	if req == nil || req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	shortCode, err := h.svc.ShortenURL(ctx, req.Url)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			return nil, status.Error(codes.InvalidArgument, "invalid URL")
		}
		return nil, status.Error(codes.Internal, "failed to shorten URL")
	}
	return &pb.ShortURL{Alias: shortCode}, nil
}

func (h *URLHandler) GetLongURL(ctx context.Context, req *pb.ShortURL) (*pb.LongURL, error) {
	if req == nil || req.Alias == "" {
		return nil, status.Error(codes.InvalidArgument, "alias is required")
	}
	longURL, err := h.svc.GetLongURL(ctx, req.Alias)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "url not found")
		}
		return nil, status.Error(codes.Internal, "failed to fetch url")
	}
	return &pb.LongURL{Url: longURL}, nil

}

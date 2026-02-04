package grpc

import (
	"context"
	"fmt"
	pb "zipit/gen/url"
	"zipit/internal/url/service"
)

type URLHandler struct {
	pb.UnimplementedURLServiceServer
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) PostURL(ctx context.Context, req *pb.LongURL) (*pb.ShortURL, error) {
	shortCode, err := h.svc.ShortenURL(ctx, req.Url)
	if err != nil {
		return nil, fmt.Errorf("probably need to return 501 internal error?")
	}
	return &pb.ShortURL{Alias: shortCode}, nil
}

func (h *URLHandler) GetLongURL(ctx context.Context, req *pb.ShortURL) (*pb.LongURL, error) {
	longURL, err := h.svc.GetLongURL(ctx, req.Alias)
	if err != nil {
		return nil, fmt.Errorf("same as before, return 501?")
	}
	return &pb.LongURL{Url: longURL}, nil

}

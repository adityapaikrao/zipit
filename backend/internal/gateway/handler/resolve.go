package handler

import (
	"net/http"

	pb "zipit/gen/url"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResolveURL handles GET /api/{code}
func (h *GatewayHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeJSONError(w, http.StatusBadRequest, "short code is required")
		return
	}

	resp, err := h.urlSvc.GetLongURL(r.Context(), &pb.ShortURL{Alias: code})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.NotFound {
			writeJSONError(w, http.StatusNotFound, "short url not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to resolve url")
		return
	}

	writeJSON(w, http.StatusOK, ResolveResponse{LongURL: resp.GetUrl()})
}

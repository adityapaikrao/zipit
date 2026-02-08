package handler

import (
	"net/http"
	"regexp"

	pb "zipit/gen/url"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validShortCode = regexp.MustCompile(`^[0-9a-zA-Z]{1,12}$`)

// ResolveURL handles GET /api/{code}
func (h *GatewayHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeJSONError(w, http.StatusBadRequest, "short code is required")
		return
	}
	if !validShortCode.MatchString(code) {
		writeJSONError(w, http.StatusBadRequest, "invalid short code format")
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

	http.Redirect(w, r, resp.GetUrl(), http.StatusFound)
}

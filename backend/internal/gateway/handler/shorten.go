package handler

import (
	"encoding/json"
	"net/http"

	pb "zipit/gen/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenURL handles POST /api/shorten
func (h *GatewayHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req PostURLRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	if req.LongURL == "" {
		writeJSONError(w, http.StatusBadRequest, "long_url is required")
		return
	}

	resp, err := h.urlSvc.PostURL(r.Context(), &pb.LongURL{Url: req.LongURL})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.InvalidArgument {
			writeJSONError(w, http.StatusBadRequest, "invalid url")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to shorten url")
		return
	}

	writeJSON(w, http.StatusCreated, ShortenResponse{ShortCode: resp.Alias})
}

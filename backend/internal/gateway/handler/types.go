package handler

type PostURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

type ResolveResponse struct {
	LongURL string `json:"long_url"`
}

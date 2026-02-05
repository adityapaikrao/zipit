package handler

type GetLongURLRequest struct {
	ShortCode string `json:"short_code"`
}

type PostURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

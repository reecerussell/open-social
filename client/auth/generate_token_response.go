package auth

// GenerateTokenResponse is the response body of the request.
type GenerateTokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

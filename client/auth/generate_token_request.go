package auth

// GenerateTokenRequest is the response body of the request.
type GenerateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

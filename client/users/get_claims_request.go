package users

// GetClaimsRequest represents the request body of a get claims resquest.
type GetClaimsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

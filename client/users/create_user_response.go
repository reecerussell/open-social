package users

// CreateUserResponse represents the response body of a create user request.
type CreateUserResponse struct {
	ReferenceID string `json:"referenceId"`
	Username    string `json:"username"`
}

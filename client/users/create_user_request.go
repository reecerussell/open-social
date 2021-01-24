package users

// CreateUserRequest represents the request body of a create user resquest.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

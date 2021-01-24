package posts

// CreateRequest is the body of the request.
type CreateRequest struct {
	UserReferenceID string `json:"userReferenceId"`
	Caption         string `json:"caption"`
}

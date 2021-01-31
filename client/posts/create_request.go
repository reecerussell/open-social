package posts

// CreateRequest is the body of the request.
type CreateRequest struct {
	UserReferenceID string `json:"userReferenceId"`
	MediaID         *int   `json:"mediaId"`
	Caption         string `json:"caption"`
}

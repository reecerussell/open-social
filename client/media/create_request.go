package media

// CreateRequest is the body of the request.
type CreateRequest struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

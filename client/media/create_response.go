package media

// CreateResponse is the body of the response.
type CreateResponse struct {
	ID          int    `json:"id"`
	ReferenceID string `json:"referenceId"`
}

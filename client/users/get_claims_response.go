package users

// GetClaimsResponse represents the response body of a get claims request.
type GetClaimsResponse struct {
	Claims map[string]interface{} `json:"claims"`
}

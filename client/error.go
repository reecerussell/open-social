package client

// Error is a custom error type which has a status code.
type Error struct {
	StatusCode int
	Message    string
}

// NewError returns a new instance of Error
func NewError(status int, message string) error {
	return &Error{
		StatusCode: status,
		Message:    message,
	}
}

// Error returns the error message.
func (err *Error) Error() string {
	return err.Message
}

// ErrorResponse represents the error response of a request.
type ErrorResponse struct {
	Message string `json:"message"`
}

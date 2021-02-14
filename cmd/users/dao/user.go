package dao

// User is a data access object for the user domain model.
type User struct {
	ID           int
	ReferenceID  string
	Username     string
	PasswordHash string

	// IsFollowing indicates wether the requesting user
	// is following this user or not.
	IsFollowing bool
}

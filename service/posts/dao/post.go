package dao

import (
	"time"
)

// Post is a data access object for the post domain.
type Post struct {
	ID          int
	ReferenceID string
	UserID      int
	Posted      time.Time
	Caption     string
}

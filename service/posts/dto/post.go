package dto

import "time"

// Post is data transfer object used to get a post's data.
type Post struct {
	ID       string    `json:"id"`
	MediaID  *string   `json:"mediaId"`
	Posted   time.Time `json:"posted"`
	Username string    `json:"username"`
	Caption  string    `json:"caption"`
	Likes    int       `json:"likes"`
	HasLiked bool      `json:"hasLiked"`
}

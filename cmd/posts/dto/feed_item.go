package dto

import "time"

// FeedItem represents a post in a feed.
type FeedItem struct {
	ID       string    `json:"id"`
	MediaID  *string   `json:"mediaId"`
	Caption  string    `json:"caption"`
	Posted   time.Time `json:"posted"`
	Username string    `json:"username"`
	Likes    int       `json:"likes"`
	HasLiked bool      `json:"hasLiked"`
	IsAuthor bool      `json:"isAuthor"`
}

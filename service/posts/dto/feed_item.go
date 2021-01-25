package dto

import "time"

// FeedItem represents a post in a feed.
type FeedItem struct {
	ID           string    `json:"id"`
	Caption      string    `json:"caption"`
	Posted       time.Time `json:"posted"`
	Username     string    `json:"username"`
	Likes        int       `json:"likes"`
	HasUserLiked bool      `json:"hasUserLiked"`
}

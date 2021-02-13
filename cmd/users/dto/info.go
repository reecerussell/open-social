package dto

// Info is a data transfer object used to provide basic information on a user.
type Info struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	MediaID       *string `json:"mediaId"`
	FollowerCount int     `json:"followerCount"`
}

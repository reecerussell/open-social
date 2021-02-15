package users

// Profile contains the data returns from getting a user's profile.
type Profile struct {
	UserID        string  `json:"userId"`
	Username      string  `json:"username"`
	MediaID       *string `json:"mediaId"`
	Bio           *string `json:"bio"`
	FollowerCount int     `json:"followerCount"`
	IsFollowing   bool    `json:"isFollowing"`
	IsOwner       bool    `json:"isOwner"`
	PostCount     int     `json:"postCount"`
}

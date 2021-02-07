package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/reecerussell/open-social/service/posts/dao"
)

const (
	maxCaptionLength = 255
)

// Post is a domain model for the posts.
type Post struct {
	id          int
	referenceID string
	userID      int
	mediaID     *int
	posted      time.Time
	caption     string

	likeCount    int
	hasUserLiked bool
}

// NewPost returns a new instance of the Post domain model,
// providing the given data is valid. This function assumes
// userID is a valid user id.
func NewPost(userID int, mediaID *int, caption string) (*Post, error) {
	p := &Post{
		userID:  userID,
		mediaID: mediaID,
		posted:  time.Now().UTC(),
	}

	err := p.updateCaption(caption)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ID returns the post's id.
func (p *Post) ID() int {
	return p.id
}

// ReferenceID returns the post's reference id.
func (p *Post) ReferenceID() string {
	return p.referenceID
}

func (p *Post) updateCaption(caption string) error {
	caption = strings.TrimSpace(caption)

	if caption == "" {
		return errors.New("caption cannot be empty")
	}

	if len(caption) > maxCaptionLength {
		return fmt.Errorf("caption cannot be greater than %d characters long", maxCaptionLength)
	}

	p.caption = caption

	return nil
}

// SetID sets the id of the post.
func (p *Post) SetID(id int) {
	p.id = id
}

// SetReferenceID sets the reference id of the post.
func (p *Post) SetReferenceID(referenceID string) {
	p.referenceID = referenceID
}

// Dao returns a data access object populated with the post's data.
func (p *Post) Dao() *dao.Post {
	return &dao.Post{
		ID:           p.id,
		ReferenceID:  p.referenceID,
		MediaID:      p.mediaID,
		UserID:       p.userID,
		Posted:       p.posted,
		Caption:      p.caption,
		LikeCount:    p.likeCount,
		HasUserLiked: p.hasUserLiked,
	}
}

// PostFromDao returns a new instance of Post, populated with the
// data from the data access object. This should only be used
// by the PostRepository, to instantiate new domain models.
func PostFromDao(d *dao.Post) *Post {
	return &Post{
		id:           d.ID,
		referenceID:  d.ReferenceID,
		mediaID:      d.MediaID,
		userID:       d.UserID,
		posted:       d.Posted,
		caption:      d.Caption,
		likeCount:    d.LikeCount,
		hasUserLiked: d.HasUserLiked,
	}
}

// CanLike determines if a user can like this post or not. An error is returned,
// if the user cannot like it.
func (p *Post) CanLike() error {
	if p.hasUserLiked {
		return errors.New("user has already liked this post")
	}

	return nil
}

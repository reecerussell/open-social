package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/reecerussell/open-social/cmd/posts/dao"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	testMediaID := 321
	p, err := NewPost(123, &testMediaID, "My first post  ")
	assert.NoError(t, err)
	assert.Equal(t, 123, p.userID)
	assert.Equal(t, testMediaID, *p.mediaID)
	assert.Equal(t, "My first post", p.caption)
}

func TestPost_UpdateCaption_ReturnsError(t *testing.T) {
	t.Run("Empty Caption", func(t *testing.T) {
		p, err := NewPost(123, nil, "")
		assert.Nil(t, p)
		assert.Equal(t, "caption cannot be empty", err.Error())
	})

	t.Run("Long Caption", func(t *testing.T) {
		str := make([]rune, maxCaptionLength+1)
		p, err := NewPost(123, nil, string(str))
		assert.Nil(t, p)

		exp := fmt.Sprintf("caption cannot be greater than %d characters long", maxCaptionLength)
		assert.Equal(t, exp, err.Error())
	})
}

func TestPost_SetID(t *testing.T) {
	const testID = 123

	var post Post
	post.SetID(testID)

	assert.Equal(t, testID, post.id)
}

func TestPost_SetReferenceID(t *testing.T) {
	const testReferenceID = "273021"

	var post Post
	post.SetReferenceID(testReferenceID)

	assert.Equal(t, testReferenceID, post.referenceID)
}

func TestPost_Dao(t *testing.T) {
	const (
		testPostID       = 1
		testReferenceID  = "3y294"
		testUserID       = 23
		testCaption      = "Hello World"
		testLikeCount    = 12
		testHasUserLiked = true
	)
	testPostedDate := time.Now().UTC()
	testMediaID := 321

	post := &Post{
		id:           testPostID,
		referenceID:  testReferenceID,
		mediaID:      &testMediaID,
		userID:       testUserID,
		posted:       testPostedDate,
		caption:      testCaption,
		likeCount:    testLikeCount,
		hasUserLiked: testHasUserLiked,
	}

	d := post.Dao()

	assert.Equal(t, testPostID, d.ID)
	assert.Equal(t, testReferenceID, d.ReferenceID)
	assert.Equal(t, testUserID, d.UserID)
	assert.Equal(t, testMediaID, *d.MediaID)
	assert.Equal(t, testPostedDate, d.Posted)
	assert.Equal(t, testCaption, d.Caption)
	assert.Equal(t, testLikeCount, d.LikeCount)
	assert.Equal(t, testHasUserLiked, d.HasUserLiked)
}

func TestPostFromDao(t *testing.T) {
	const (
		testPostID       = 1
		testReferenceID  = "3y294"
		testUserID       = 23
		testCaption      = "Hello World"
		testLikeCount    = 12
		testHasUserLiked = true
	)
	testPostedDate := time.Now().UTC()
	testMediaID := 10

	d := &dao.Post{
		ID:           testPostID,
		ReferenceID:  testReferenceID,
		MediaID:      &testMediaID,
		UserID:       testUserID,
		Posted:       testPostedDate,
		Caption:      testCaption,
		LikeCount:    testLikeCount,
		HasUserLiked: testHasUserLiked,
	}

	post := PostFromDao(d)

	assert.Equal(t, testPostID, post.id)
	assert.Equal(t, testReferenceID, post.referenceID)
	assert.Equal(t, testMediaID, *post.mediaID)
	assert.Equal(t, testUserID, post.userID)
	assert.Equal(t, testPostedDate, post.posted)
	assert.Equal(t, testCaption, post.caption)
	assert.Equal(t, testLikeCount, post.likeCount)
	assert.Equal(t, testHasUserLiked, post.hasUserLiked)
}

func TestPost_ReferenceID(t *testing.T) {
	const testReferenceID = "38924yhwd"

	post := &Post{referenceID: testReferenceID}

	assert.Equal(t, testReferenceID, post.ReferenceID())
}

func TestPost_ID(t *testing.T) {
	const testID = 123

	post := &Post{id: testID}

	assert.Equal(t, testID, post.ID())
}

func TestPost_CanLike(t *testing.T) {
	post := &Post{
		hasUserLiked: false,
	}

	err := post.CanLike()
	assert.NoError(t, err)
}

func TestPost_CanLike_ReturnsError(t *testing.T) {
	t.Run("User Already Liked", func(t *testing.T) {
		post := &Post{
			hasUserLiked: true,
		}

		err := post.CanLike()
		assert.Equal(t, "user has already liked this post", err.Error())
	})
}

func TestPost_CanUnike(t *testing.T) {
	post := &Post{
		hasUserLiked: true,
	}

	err := post.CanUnlike()
	assert.NoError(t, err)
}

func TestPost_CanUnlike_ReturnsError(t *testing.T) {
	t.Run("User Not Liked", func(t *testing.T) {
		post := &Post{
			hasUserLiked: false,
		}

		err := post.CanUnlike()
		assert.Equal(t, "user has not liked this post", err.Error())
	})
}

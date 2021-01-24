package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	p, err := NewPost(123, "My first post  ")
	assert.NoError(t, err)
	assert.Equal(t, 123, p.userID)
	assert.Equal(t, "My first post", p.caption)
}

func TestPost_UpdateCaption_ReturnsError(t *testing.T) {
	t.Run("Empty Caption", func(t *testing.T) {
		p, err := NewPost(123, "")
		assert.Nil(t, p)
		assert.Equal(t, "caption cannot be empty", err.Error())
	})

	t.Run("Long Caption", func(t *testing.T) {
		str := make([]rune, maxCaptionLength+1)
		p, err := NewPost(123, string(str))
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
		testPostID      = 1
		testReferenceID = "3y294"
		testUserID      = 23
		testCaption     = "Hello World"
	)
	testPostedDate := time.Now().UTC()

	post := &Post{
		id:          testPostID,
		referenceID: testReferenceID,
		userID:      testUserID,
		posted:      testPostedDate,
		caption:     testCaption,
	}

	d := post.Dao()

	assert.Equal(t, testPostID, d.ID)
	assert.Equal(t, testReferenceID, d.ReferenceID)
	assert.Equal(t, testUserID, d.UserID)
	assert.Equal(t, testPostedDate, d.Posted)
	assert.Equal(t, testCaption, d.Caption)
}

func TestPost_ReferenceID(t *testing.T) {
	const testReferenceID = "38924yhwd"

	post := &Post{referenceID: testReferenceID}

	assert.Equal(t, testReferenceID, post.ReferenceID())
}

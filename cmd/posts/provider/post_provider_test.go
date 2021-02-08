package provider

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/reecerussell/open-social/mock/database"
)

// post found
func TestPostProvider_Get_ReturnsPostSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "12037021"
	testUserReferenceID := "07213042"
	testMediaID := "2379470324"
	testPosted := time.Now().UTC()
	testUsername := "test"
	testCaption := "Hello World"
	testLikes := 12
	testHasUserLiked := true
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		*(dest[0].(*string)) = testPostReferenceID
		*(dest[1].(**string)) = &testMediaID
		*(dest[2].(*time.Time)) = testPosted
		*(dest[3].(*string)) = testUsername
		*(dest[4].(*string)) = testCaption
		*(dest[5].(*int)) = testLikes
		*(dest[6].(*bool)) = testHasUserLiked

		return nil
	})

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	p := NewPostProvider(mockDatabase)
	post, err := p.Get(testCtx, testPostReferenceID, testUserReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, testPostReferenceID, post.ID)
	assert.Equal(t, &testMediaID, post.MediaID)
	assert.Equal(t, testPosted, post.Posted)
	assert.Equal(t, testUsername, post.Username)
	assert.Equal(t, testCaption, post.Caption)
	assert.Equal(t, testLikes, post.Likes)
	assert.Equal(t, testHasUserLiked, post.HasLiked)
}

// post not found
func TestPostProvider_GetNonExistantPost_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "12037021"
	testUserReferenceID := "07213042"
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	p := NewPostProvider(mockDatabase)
	post, err := p.Get(testCtx, testPostReferenceID, testUserReferenceID)
	assert.Nil(t, post)
	assert.Equal(t, ErrPostNotFound, err)
}

// db err
func TestPostProvider_GetScanFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "12037021"
	testUserReferenceID := "07213042"
	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).Return(testError)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	p := NewPostProvider(mockDatabase)
	post, err := p.Get(testCtx, testPostReferenceID, testUserReferenceID)
	assert.Nil(t, post)
	assert.Equal(t, testError, err)
}

func TestPostProvider_GetQueryFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "12037021"
	testUserReferenceID := "07213042"
	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(nil, testError)

	p := NewPostProvider(mockDatabase)
	post, err := p.Get(testCtx, testPostReferenceID, testUserReferenceID)
	assert.Nil(t, post)
	assert.Equal(t, testError, err)
}

func TestPostProvider_GetProfileFeed_ReturnsProfileFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3740423"
	testPostID := "2349734"
	testMediaID := "3204703"
	testCaption := "Hello World"
	testPosted := time.Now().UTC()
	testUsername := "test"
	testLikes := 12
	testHasUserLiked := true
	testIsAuthor := false
	testCtx := context.Background()

	readCount := 0

	mockRows := mock.NewMockRows(ctrl)
	mockRows.EXPECT().Next().DoAndReturn(func() bool {
		if readCount > 0 {
			return false
		}

		readCount++
		return true
	}).Times(2)
	mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		*(dest[0].(*string)) = testPostID
		*(dest[1].(**string)) = &testMediaID
		*(dest[2].(*string)) = testCaption
		*(dest[3].(*time.Time)) = testPosted
		*(dest[4].(*string)) = testUsername
		*(dest[5].(*int)) = testLikes
		*(dest[6].(*bool)) = testHasUserLiked
		*(dest[7].(*bool)) = testIsAuthor

		return nil
	})
	mockRows.EXPECT().Err().Return(nil)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Multiple(testCtx, gomock.Any(), gomock.Any()).Return(mockRows, nil)

	provider := NewPostProvider(mockDatabase)
	feedItems, err := provider.GetProfileFeed(testCtx, testUsername, testUserReferenceID)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(feedItems))
	assert.Equal(t, testPostID, feedItems[0].ID)
	assert.Equal(t, &testMediaID, feedItems[0].MediaID)
	assert.Equal(t, testCaption, feedItems[0].Caption)
	assert.Equal(t, testPosted, feedItems[0].Posted)
	assert.Equal(t, testUsername, feedItems[0].Username)
	assert.Equal(t, testLikes, feedItems[0].Likes)
	assert.Equal(t, testHasUserLiked, feedItems[0].HasUserLiked)
	assert.Equal(t, testIsAuthor, feedItems[0].IsAuthor)
}

func TestPostProvider_GetProfileFeedQueryFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3740423"
	testUsername := "test"
	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Multiple(testCtx, gomock.Any(), gomock.Any()).Return(nil, testError)

	provider := NewPostProvider(mockDatabase)
	feedItems, err := provider.GetProfileFeed(testCtx, testUsername, testUserReferenceID)
	assert.Nil(t, feedItems)
	assert.Equal(t, testError, err)
}

func TestPostProvider_GetProfileFeedScanFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3740423"
	testUsername := "test"
	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockRows := mock.NewMockRows(ctrl)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(gomock.Any()).Return(testError)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Multiple(testCtx, gomock.Any(), gomock.Any()).Return(mockRows, nil)

	provider := NewPostProvider(mockDatabase)
	feedItems, err := provider.GetProfileFeed(testCtx, testUsername, testUserReferenceID)
	assert.Nil(t, feedItems)
	assert.Equal(t, testError, err)
}

func TestPostProvider_GetProfileFeedRowsErrors_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3740423"
	testUsername := "test"
	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockRows := mock.NewMockRows(ctrl)
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Err().Return(testError)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Multiple(testCtx, gomock.Any(), gomock.Any()).Return(mockRows, nil)

	provider := NewPostProvider(mockDatabase)
	feedItems, err := provider.GetProfileFeed(testCtx, testUsername, testUserReferenceID)
	assert.Nil(t, feedItems)
	assert.Equal(t, testError, err)
}

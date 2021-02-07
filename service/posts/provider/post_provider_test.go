package provider

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/reecerussell/open-social/core/mock/database"
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

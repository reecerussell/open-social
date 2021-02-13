package provider

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/reecerussell/open-social/mock/database"
)

func TestUserProvider_GetProfile_ReturnsProfileSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "320434"
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*string)) = "test"
			*(dest[1].(**string)) = nil
			*(dest[2].(**string)) = nil
			*(dest[3].(*int)) = 10
			*(dest[4].(*bool)) = true
			*(dest[5].(*bool)) = false
			*(dest[6].(*int)) = 5

			return nil
		})

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	provider := NewUserProvider(mockDatabase)
	profile, err := provider.GetProfile(testCtx, testUsername, testUserReferenceID)
	assert.NoError(t, err)

	assert.Equal(t, testUsername, profile.Username)
	assert.Nil(t, profile.MediaID)
	assert.Nil(t, profile.Bio)
	assert.Equal(t, 10, profile.FollowerCount)
	assert.True(t, profile.IsFollowing)
	assert.False(t, profile.IsOwner)
	assert.Equal(t, 5, profile.PostCount)
}

func TestUserProvider_GetProfileQueryFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "320434"
	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(nil, testError)

	provider := NewUserProvider(mockDatabase)
	profile, err := provider.GetProfile(testCtx, testUsername, testUserReferenceID)
	assert.Nil(t, profile)
	assert.Equal(t, testError, err)
}

func TestUserProvider_GetProfileDoesNotExist_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "320434"
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	provider := NewUserProvider(mockDatabase)
	profile, err := provider.GetProfile(testCtx, testUsername, testUserReferenceID)
	assert.Nil(t, profile)
	assert.Equal(t, ErrProfileNotFound, err)
}

func TestUserProvider_GetInfo_ReturnsInfoSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "320434"
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			*(dest[0].(*string)) = testUserReferenceID
			*(dest[1].(*string)) = "test"
			*(dest[2].(**string)) = nil
			*(dest[3].(*int)) = 10

			return nil
		})

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	provider := NewUserProvider(mockDatabase)
	info, err := provider.GetInfo(testCtx, testUserReferenceID)
	assert.NoError(t, err)

	assert.Equal(t, testUserReferenceID, info.ID)
	assert.Equal(t, "test", info.Username)
	assert.Nil(t, info.MediaID)
	assert.Equal(t, 10, info.FollowerCount)
}

func TestUserProvider_GetInfoQueryFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "320434"
	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(nil, testError)

	provider := NewUserProvider(mockDatabase)
	info, err := provider.GetInfo(testCtx, testUserReferenceID)
	assert.Nil(t, info)
	assert.Equal(t, testError, err)
}

func TestUserProvider_GetInfoDoesNotExist_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "320434"
	testCtx := context.Background()

	mockRow := mock.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Single(testCtx, gomock.Any(), gomock.Any()).Return(mockRow, nil)

	provider := NewUserProvider(mockDatabase)
	info, err := provider.GetInfo(testCtx, testUserReferenceID)
	assert.Nil(t, info)
	assert.Equal(t, ErrProfileNotFound, err)
}

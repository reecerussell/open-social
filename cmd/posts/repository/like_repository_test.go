package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/reecerussell/open-social/mock/database"
)

func TestLikeRepository_Create_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostID := 1
	testUserReferenceID := "37947230"
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(1), nil)

	repo := NewLikeRepository(mockDatabase)
	err := repo.Create(testCtx, testPostID, testUserReferenceID)
	assert.NoError(t, err)
}

func TestLikeRepository_CreateExecuteFails_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostID := 1
	testUserReferenceID := "37947230"
	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(-1), testError)

	repo := NewLikeRepository(mockDatabase)
	err := repo.Create(testCtx, testPostID, testUserReferenceID)
	assert.Equal(t, testError, err)
}

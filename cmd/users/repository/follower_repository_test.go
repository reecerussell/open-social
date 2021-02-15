package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/reecerussell/open-social/mock/database"
)

func TestFollowerRepository_Create_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(1), nil)

	repo := NewFollowerRepository(mockDatabase)

	err := repo.Create(testCtx, 10, "390274jlw")
	assert.NoError(t, err)
}

func TestFollowerRepository_CreateExecutionFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(-1), testError)

	repo := NewFollowerRepository(mockDatabase)

	err := repo.Create(testCtx, 10, "390274jlw")
	assert.Equal(t, testError, err)
}

func TestFollowerRepository_Delete_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(1), nil)

	repo := NewFollowerRepository(mockDatabase)

	err := repo.Delete(testCtx, 10, "390274jlw")
	assert.NoError(t, err)
}

func TestFollowerRepository_DeleteExecutionFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testError := errors.New("an error occured")
	testCtx := context.Background()

	mockDatabase := mock.NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Execute(testCtx, gomock.Any(), gomock.Any()).Return(int64(-1), testError)

	repo := NewFollowerRepository(mockDatabase)

	err := repo.Delete(testCtx, 10, "390274jlw")
	assert.Equal(t, testError, err)
}

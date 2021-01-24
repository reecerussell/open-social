package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	hashermock "github.com/reecerussell/adaptive-password-hasher/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/service/users/mock"
	"github.com/reecerussell/open-social/service/users/mock/repository"
	"github.com/reecerussell/open-social/service/users/model"
)

func TestCreateUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUsername    = "test-create"
		testPassword    = "Password1"
		testID          = 1
		testReferenceID = "68130159-39f4-42f6-be68-4379d1ac9b11"
	)

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte(testPassword))

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().DoesUsernameExist(gomock.Any(), testUsername, nil).Return(false, nil)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *model.User) error {
			u.SetID(testID)
			u.SetReferenceID(testReferenceID)

			return nil
		})

	handler := NewCreateUserHandler(mockValidator, mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"referenceId\":\"%s\",\"username\":\"%s\"}\n", testReferenceID, testUsername)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreateUserHandler_GivenInvalidUserData_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUsername     = "test-create"
		testPassword     = "Password1"
		testErrorMessage = "an error occured"
	)

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(errors.New(testErrorMessage))

	handler := NewCreateUserHandler(mockValidator, nil, nil)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreateUserHandler_DoesUsernameExistError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUsername     = "test-create"
		testPassword     = "Password1"
		testID           = 1
		testReferenceID  = "68130159-39f4-42f6-be68-4379d1ac9b11"
		testErrorMessage = "an error occured"
	)

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte(testPassword))

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().DoesUsernameExist(gomock.Any(), testUsername, nil).Return(false, errors.New(testErrorMessage))

	handler := NewCreateUserHandler(mockValidator, mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreateUserHandler_UsernameAlreadyExists_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUsername = "test-create"
		testPassword = "Password1"
	)

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte(testPassword))

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().DoesUsernameExist(gomock.Any(), testUsername, nil).Return(true, nil)

	handler := NewCreateUserHandler(mockValidator, mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"the username '%s' is taken\"}\n", testUsername)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreateUserHandler_CreateError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUsername     = "test-create"
		testPassword     = "Password1"
		testID           = 1
		testReferenceID  = "68130159-39f4-42f6-be68-4379d1ac9b11"
		testErrorMessage = "an error occured"
	)

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte(testPassword))

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().DoesUsernameExist(gomock.Any(), testUsername, nil).Return(false, nil)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New(testErrorMessage))

	handler := NewCreateUserHandler(mockValidator, mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

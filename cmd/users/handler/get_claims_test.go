package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	hashermock "github.com/reecerussell/adaptive-password-hasher/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/users/dao"
	"github.com/reecerussell/open-social/cmd/users/mock/repository"
	"github.com/reecerussell/open-social/cmd/users/model"
	repo "github.com/reecerussell/open-social/cmd/users/repository"
)

func getMockUser(referenceID, username string) *model.User {
	return model.NewUserFromDao(&dao.User{
		ID:           1,
		ReferenceID:  referenceID,
		Username:     username,
		PasswordHash: "t6832ihlw",
	})
}

func TestGetClaimsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUsername = "testing"
	const testReferenceID = "23470324"
	const testPassword = "Password_123"
	testUser := getMockUser(testReferenceID, testUsername)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(true)

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(testUser, nil)

	handler := NewGetClaimsHandler(mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"claims\":{\"uid\":\"%s\",\"username\":\"%s\"}}\n", testReferenceID, testUsername)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestGetClaimsHandler_UserDoesNotExist_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUsername = "testing"
	const testPassword = "Password_123"

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, repo.ErrUserNotFound)

	handler := NewGetClaimsHandler(nil, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%v\"}\n", repo.ErrUserNotFound)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestGetClaimsHandler_FailedToGetUser_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUsername = "testing"
	const testPassword = "Password_123"
	const errorMessage = "an error occured"

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(nil, errors.New(errorMessage))

	handler := NewGetClaimsHandler(nil, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%v\"}\n", errorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestGetClaimsHandler_GivenInvalidPassword_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUsername = "testing"
	const testReferenceID = "23470324"
	const testPassword = "Password_123"
	testUser := getMockUser(testReferenceID, testUsername)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(false)

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetUserByUsername(gomock.Any(), testUsername).Return(testUser, nil)

	handler := NewGetClaimsHandler(mockHasher, mockRepo)

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, testUsername, testPassword)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%v\"}\n", model.ErrInvalidPassword)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

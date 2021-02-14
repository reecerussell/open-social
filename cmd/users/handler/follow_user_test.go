package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/users/dao"
	mock "github.com/reecerussell/open-social/cmd/users/mock/repository"
	"github.com/reecerussell/open-social/cmd/users/model"
	"github.com/reecerussell/open-social/cmd/users/repository"
)

func TestFollowUser_GivenValidData_ReturnsOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "32047023"
	const testUserID = 1203
	const testFollowerReferenceID = "01379nldsd"

	mockUsers := mock.NewMockUserRepository(ctrl)
	mockUsers.EXPECT().GetUserByReference(gomock.Any(), testUserReferenceID, testFollowerReferenceID).
		Return(model.NewUserFromDao(&dao.User{ID: testUserID, IsFollowing: false}), nil)

	mockFollowers := mock.NewMockFollowerRepository(ctrl)
	mockFollowers.EXPECT().Create(gomock.Any(), testUserID, testFollowerReferenceID).
		Return(nil)

	handler := NewFollowUserHandler(mockUsers, mockFollowers)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}/{followerReferenceId}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testUserReferenceID, testFollowerReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFollowUser_UserDoesNotExist_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "32047023"
	const testUserID = 1203
	const testFollowerReferenceID = "01379nldsd"

	mockUsers := mock.NewMockUserRepository(ctrl)
	mockUsers.EXPECT().GetUserByReference(gomock.Any(), testUserReferenceID, testFollowerReferenceID).
		Return(nil, repository.ErrUserNotFound)

	handler := NewFollowUserHandler(mockUsers, nil)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}/{followerReferenceId}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testUserReferenceID, testFollowerReferenceID), nil)
	router.ServeHTTP(rr, req)

	var data map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&data)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, repository.ErrUserNotFound.Error(), data["message"])
}

func TestFollowUser_UserAlreadyFollowing_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "32047023"
	const testUserID = 1203
	const testFollowerReferenceID = "01379nldsd"

	mockUsers := mock.NewMockUserRepository(ctrl)
	mockUsers.EXPECT().GetUserByReference(gomock.Any(), testUserReferenceID, testFollowerReferenceID).
		Return(model.NewUserFromDao(&dao.User{ID: testUserID, IsFollowing: true}), nil)

	handler := NewFollowUserHandler(mockUsers, nil)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}/{followerReferenceId}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testUserReferenceID, testFollowerReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestFollowUser_FollowerNotFound_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "32047023"
	const testUserID = 1203
	const testFollowerReferenceID = "01379nldsd"

	mockUsers := mock.NewMockUserRepository(ctrl)
	mockUsers.EXPECT().GetUserByReference(gomock.Any(), testUserReferenceID, testFollowerReferenceID).
		Return(model.NewUserFromDao(&dao.User{ID: testUserID, IsFollowing: false}), nil)

	mockFollowers := mock.NewMockFollowerRepository(ctrl)
	mockFollowers.EXPECT().Create(gomock.Any(), testUserID, testFollowerReferenceID).
		Return(repository.ErrFollowerNotFound)

	handler := NewFollowUserHandler(mockUsers, mockFollowers)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}/{followerReferenceId}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", testUserReferenceID, testFollowerReferenceID), nil)
	router.ServeHTTP(rr, req)

	var data map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&data)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, repository.ErrFollowerNotFound.Error(), data["message"])
}

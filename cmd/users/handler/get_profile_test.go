package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/users/dto"
	mock "github.com/reecerussell/open-social/cmd/users/mock/provider"
	"github.com/reecerussell/open-social/cmd/users/provider"
)

func TestGetProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "3274032"
	testMediaID := "30743555"
	testBio := "Hello World"
	testFollowerCount := 10
	testIsFollowing := true
	testIsOwner := false
	testPostCount := 15

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetProfile(gomock.Any(), testUsername, testUserReferenceID).
		Return(&dto.Profile{
			Username:      testUsername,
			MediaID:       &testMediaID,
			Bio:           &testBio,
			FollowerCount: testFollowerCount,
			IsFollowing:   testIsFollowing,
			IsOwner:       testIsOwner,
			PostCount:     testPostCount,
		}, nil)

	handler := NewGetProfileHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{username}/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testUsername, testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)

	assert.Equal(t, testUsername, data["username"])
	assert.Equal(t, testMediaID, data["mediaId"])
	assert.Equal(t, testBio, data["bio"])
	assert.Equal(t, float64(testFollowerCount), data["followerCount"])
	assert.Equal(t, testIsFollowing, data["isFollowing"])
	assert.Equal(t, testIsOwner, data["isOwner"])
	assert.Equal(t, float64(testPostCount), data["postCount"])
}

func TestGetProfileHandler_ProfileDoesNotExist_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "3274032"

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetProfile(gomock.Any(), testUsername, testUserReferenceID).Return(nil, provider.ErrProfileNotFound)

	handler := NewGetProfileHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{username}/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testUsername, testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, provider.ErrProfileNotFound.Error(), data["message"])
}

func TestGetProfileHandler_ProviderReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "3274032"
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetProfile(gomock.Any(), testUsername, testUserReferenceID).Return(nil, testError)

	handler := NewGetProfileHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{username}/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testUsername, testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, testError.Error(), data["message"])
}

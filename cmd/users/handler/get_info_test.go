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

func TestGetInfoHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "3274032"
	testMediaID := "30743555"
	testFollowerCount := 10

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetInfo(gomock.Any(), testUserReferenceID).
		Return(&dto.Info{
			ID:            testUserReferenceID,
			Username:      testUsername,
			MediaID:       &testMediaID,
			FollowerCount: testFollowerCount,
		}, nil)

	handler := NewGetInfoHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)

	assert.Equal(t, testUserReferenceID, data["id"])
	assert.Equal(t, testUsername, data["username"])
	assert.Equal(t, testMediaID, data["mediaId"])
	assert.Equal(t, float64(testFollowerCount), data["followerCount"])
}

func TestGetInfoHandler_ProfileDoesNotExist_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3274032"

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetInfo(gomock.Any(), testUserReferenceID).Return(nil, provider.ErrProfileNotFound)

	handler := NewGetInfoHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, provider.ErrProfileNotFound.Error(), data["message"])
}

func TestGetInfoHandler_ProviderReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "3274032"
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockUserProvider(ctrl)
	mockProvider.EXPECT().GetInfo(gomock.Any(), testUserReferenceID).Return(nil, testError)

	handler := NewGetInfoHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Equal(t, testError.Error(), data["message"])
}

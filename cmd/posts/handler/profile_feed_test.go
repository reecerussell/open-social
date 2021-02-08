package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/posts/dto"
	mock "github.com/reecerussell/open-social/cmd/posts/mock/provider"
)

func TestProfileFeedHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "2398yhlwd"
	testPostedDate := time.Now()

	mockProvider := mock.NewMockPostProvider(ctrl)
	mockProvider.EXPECT().GetProfileFeed(gomock.Any(), testUsername, testUserReferenceID).Return([]*dto.FeedItem{
		{
			ID:           "23123",
			Caption:      "Hello World",
			Posted:       testPostedDate,
			Username:     testUsername,
			Likes:        1,
			HasUserLiked: false,
			IsAuthor:     true,
		},
	}, nil)

	handler := NewProfileFeedHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{username}/{userReferenceId}", handler).Methods("GET")

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testUsername, testUserReferenceID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))

	var data []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(data))

	item := data[0]
	expPostedDate, _ := testPostedDate.MarshalText()

	assert.Equal(t, "23123", item["id"])
	assert.Equal(t, "Hello World", item["caption"])
	assert.Equal(t, string(expPostedDate), item["posted"])
	assert.Equal(t, testUsername, item["username"])
	assert.Equal(t, float64(1), item["likes"])
	assert.Equal(t, false, item["hasUserLiked"])
	assert.Equal(t, true, item["isAuthor"])
}

func TestProfileFeedHandler_ProviderReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testUserReferenceID := "2398yhlwd"
	testErrorMessage := "an error occured"

	mockProvider := mock.NewMockPostProvider(ctrl)
	mockProvider.EXPECT().GetProfileFeed(gomock.Any(), testUsername, testUserReferenceID).Return(nil, errors.New(testErrorMessage))

	handler := NewProfileFeedHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{username}/{userReferenceId}", handler).Methods("GET")

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testUsername, testUserReferenceID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

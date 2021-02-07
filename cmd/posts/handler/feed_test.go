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
	repository "github.com/reecerussell/open-social/cmd/posts/mock/repository"
)

func TestFeedHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2398yhlwd"
	testPostedDate := time.Now()

	mockRepo := repository.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().GetFeed(gomock.Any(), testUserReferenceID).Return([]*dto.FeedItem{
		{
			ID:           "23123",
			Caption:      "Hello World",
			Posted:       testPostedDate,
			Username:     "User123",
			Likes:        1,
			HasUserLiked: false,
		},
	}, nil)

	handler := NewFeedHandler(mockRepo)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}", handler).Methods("GET")

	req, _ := http.NewRequest(http.MethodGet, "/"+testUserReferenceID, nil)
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
	assert.Equal(t, "User123", item["username"])
	assert.Equal(t, float64(1), item["likes"])
	assert.Equal(t, false, item["hasUserLiked"])
}

func TestFeedHandler_RepoReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2398yhlwd"
	testErrorMessage := "an error occured"

	mockRepo := repository.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().GetFeed(gomock.Any(), testUserReferenceID).Return(nil, errors.New(testErrorMessage))

	handler := NewFeedHandler(mockRepo)
	router := mux.NewRouter()
	router.Handle("/{userReferenceId}", handler).Methods("GET")

	req, _ := http.NewRequest(http.MethodGet, "/"+testUserReferenceID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

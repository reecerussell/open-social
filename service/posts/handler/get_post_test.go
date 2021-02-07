package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/service/posts/dto"
	mock "github.com/reecerussell/open-social/service/posts/mock/provider"
	"github.com/reecerussell/open-social/service/posts/provider"
)

func TestGetPostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"
	testPost := dto.Post{
		ID:       testPostReferenceID,
		MediaID:  nil,
		Posted:   time.Now().UTC(),
		Username: "test",
		Caption:  "Hello World",
		Likes:    1,
		HasLiked: true,
	}

	mockProvider := mock.NewMockPostProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(&testPost, nil)

	handler := NewGetPostHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{postReferenceID}/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testPostReferenceID, testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, testPost.ID, data["id"])
	assert.Nil(t, data["mediaId"])

	expPostedDate, _ := testPost.Posted.MarshalText()
	assert.Equal(t, string(expPostedDate), data["posted"])
	assert.Equal(t, testPost.Username, data["username"])
	assert.Equal(t, testPost.Caption, data["caption"])
	assert.Equal(t, float64(testPost.Likes), data["likes"])
	assert.True(t, data["hasLiked"].(bool))
}

func TestGetPostHandler_WithNonExistantPost_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"

	mockProvider := mock.NewMockPostProvider(ctrl)
	mockProvider.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(nil, provider.ErrPostNotFound)

	handler := NewGetPostHandler(mockProvider)
	router := mux.NewRouter()
	router.Handle("/{postReferenceID}/{userReferenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", testPostReferenceID, testUserReferenceID), nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, provider.ErrPostNotFound.Error(), data["message"])
}

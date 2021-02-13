package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/posts/dao"
	mock "github.com/reecerussell/open-social/cmd/posts/mock/repository"
	"github.com/reecerussell/open-social/cmd/posts/model"
	"github.com/reecerussell/open-social/cmd/posts/repository"
)

func TestLikePostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostID := 12
	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"
	testPost := model.PostFromDao(&dao.Post{
		ID:       testPostID,
		HasLiked: false,
	})

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(testPost, nil)

	mockLikes := mock.NewMockLikeRepository(ctrl)
	mockLikes.EXPECT().Create(gomock.Any(), testPostID, testUserReferenceID).Return(nil)

	handler := NewLikePostHandler(mockRepo, mockLikes)
	rr := httptest.NewRecorder()

	body := fmt.Sprintf("{\"postReferenceId\":\"%s\",\"userReferenceId\":\"%s\"}", testPostReferenceID, testUserReferenceID)
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLikePostHandler_GivenNonExistantPost_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(nil, repository.ErrPostNotFound)

	handler := NewLikePostHandler(mockRepo, nil)
	rr := httptest.NewRecorder()

	body := fmt.Sprintf("{\"postReferenceId\":\"%s\",\"userReferenceId\":\"%s\"}", testPostReferenceID, testUserReferenceID)
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, repository.ErrPostNotFound.Error(), data["message"])
}

func TestLikePostHandler_WhereUserHasAlreadyLikedPost_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"
	testPost := model.PostFromDao(&dao.Post{
		HasLiked: true,
	})

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(testPost, nil)

	handler := NewLikePostHandler(mockRepo, nil)
	rr := httptest.NewRecorder()

	body := fmt.Sprintf("{\"postReferenceId\":\"%s\",\"userReferenceId\":\"%s\"}", testPostReferenceID, testUserReferenceID)
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestLikePostHandler_CreateLikeFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPostID := 12
	testPostReferenceID := "5234934"
	testUserReferenceID := "1740398"
	testError := errors.New("an error occured")
	testPost := model.PostFromDao(&dao.Post{
		ID:       testPostID,
		HasLiked: false,
	})

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any(), testPostReferenceID, testUserReferenceID).Return(testPost, nil)

	mockLikes := mock.NewMockLikeRepository(ctrl)
	mockLikes.EXPECT().Create(gomock.Any(), testPostID, testUserReferenceID).Return(testError)

	handler := NewLikePostHandler(mockRepo, mockLikes)
	rr := httptest.NewRecorder()

	body := fmt.Sprintf("{\"postReferenceId\":\"%s\",\"userReferenceId\":\"%s\"}", testPostReferenceID, testUserReferenceID)
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, testError.Error(), data["message"])
}

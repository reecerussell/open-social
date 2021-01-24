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
	"github.com/stretchr/testify/assert"

	clientMock "github.com/reecerussell/open-social/client/mock/users"
	repoMock "github.com/reecerussell/open-social/service/posts/mock/repository"
	"github.com/reecerussell/open-social/service/posts/model"
)

func TestCreatePostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testPostID          = 123
		testReferenceID     = "21932"
		testUserReferenceID = "2392"
		testCaption         = "Hello World"
	)
	testUserID := 12

	mockClient := clientMock.NewMockClient(ctrl)
	mockClient.EXPECT().GetIDByReference(testUserReferenceID).Return(&testUserID, nil)

	mockRepo := repoMock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, p *model.Post) error {
			p.SetID(testPostID)
			p.SetReferenceID(testReferenceID)

			return nil
		})

	handler := NewCreatePostHandler(mockRepo, mockClient)

	body := fmt.Sprintf(`{"userReferenceId": "%s", "caption": "%s"}`, testUserReferenceID, testCaption)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"referenceId\":\"%s\"}\n", testReferenceID)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreatePostHandler_GivenInvalidUser_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "2893ks"
	const testErrorMessage = "an error occured"

	mockClient := clientMock.NewMockClient(ctrl)
	mockClient.EXPECT().GetIDByReference(testUserReferenceID).Return(nil, errors.New(testErrorMessage))

	mockRepo := repoMock.NewMockPostRepository(ctrl)

	handler := NewCreatePostHandler(mockRepo, mockClient)

	body := fmt.Sprintf(`{"userReferenceId": "%s", "caption": "Hello World"}`, testUserReferenceID)
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

func TestCreatePostHandler_GivenInvalidPostData_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		testUserReferenceID = "2392"
		testCaption         = ""
	)
	testUserID := 12

	mockClient := clientMock.NewMockClient(ctrl)
	mockClient.EXPECT().GetIDByReference(testUserReferenceID).Return(&testUserID, nil)

	mockRepo := repoMock.NewMockPostRepository(ctrl)
	handler := NewCreatePostHandler(mockRepo, mockClient)

	body := fmt.Sprintf(`{"userReferenceId": "%s", "caption": "%s"}`, testUserReferenceID, testCaption)
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"caption cannot be empty\"}\n")
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestCreatePostHandler_RepoReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUserReferenceID = "2893ks"
	const testErrorMessage = "an error occured"
	testUserID := 12

	mockClient := clientMock.NewMockClient(ctrl)
	mockClient.EXPECT().GetIDByReference(testUserReferenceID).Return(&testUserID, nil)

	mockRepo := repoMock.NewMockPostRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New(testErrorMessage))

	handler := NewCreatePostHandler(mockRepo, mockClient)

	body := fmt.Sprintf(`{"userReferenceId": "%s", "caption": "Hello World"}`, testUserReferenceID)
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

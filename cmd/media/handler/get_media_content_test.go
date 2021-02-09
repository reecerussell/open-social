package handler

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	repoMock "github.com/reecerussell/open-social/cmd/media/mock/repository"
	"github.com/reecerussell/open-social/cmd/media/repository"
	"github.com/reecerussell/open-social/mock/media"
)

func TestGetMediaContentHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"
	const testContentType = "text/plain"
	const testContent = "SGVsbG8gV29ybGQ="

	mockRepo := repoMock.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().GetContentType(gomock.Any(), testReferenceID).Return(testContentType, nil)

	mockDownloader := media.NewMockService(ctrl)
	mockDownloader.EXPECT().Download(gomock.Any(), testReferenceID).
		DoAndReturn(func(ctx context.Context, key string) ([]byte, error) {
			return base64.StdEncoding.DecodeString(testContent)
		})

	handler := NewGetMediaContentHandler(mockRepo, mockDownloader)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"contentType\":\"%s\",\"content\":\"%s\"}\n", testContentType, testContent)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestGetMediaContentHandler_GivenInvalidReferenceID_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"

	mockRepo := repoMock.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().GetContentType(gomock.Any(), testReferenceID).Return("", repository.ErrMediaNotFound)

	handler := NewGetMediaContentHandler(mockRepo, nil)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", repository.ErrMediaNotFound)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestGetMediaContentHandler_RepoReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"
	const testErrorMessage = "an error occured"

	mockRepo := repoMock.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().GetContentType(gomock.Any(), testReferenceID).Return("", errors.New(testErrorMessage))

	handler := NewGetMediaContentHandler(mockRepo, nil)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestGetMediaContentHandler_DownloaderReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"
	const testContentType = "text/plain"
	const testErrorMessage = "an error occured"

	mockRepo := repoMock.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().GetContentType(gomock.Any(), testReferenceID).Return(testContentType, nil)

	mockDownloader := media.NewMockService(ctrl)
	mockDownloader.EXPECT().Download(gomock.Any(), testReferenceID).Return(nil, errors.New(testErrorMessage))

	handler := NewGetMediaContentHandler(mockRepo, mockDownloader)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

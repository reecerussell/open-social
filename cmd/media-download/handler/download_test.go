package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	media "github.com/reecerussell/open-social/client/mock/media"
)

func TestDownloadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"
	const testContentType = "text/plain"
	const testContent = "SGVsbG8gV29ybGQ="

	mockClient := media.NewMockClient(ctrl)
	mockClient.EXPECT().GetContent(testReferenceID).
		DoAndReturn(func(key string) (string, []byte, error) {
			bytes, _ := base64.StdEncoding.DecodeString(testContent)
			return testContentType, bytes, nil
		})

	handler := NewDownloadHandler(mockClient)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	expData, _ := base64.StdEncoding.DecodeString(testContent)
	assert.Equal(t, expData, data)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, testContentType, rr.Header().Get("Content-Type"))
	assert.Equal(t, "private, max-age=3600", rr.Header().Get("Cache-Control"))
}

func TestDownloadHandler_FailedGetContent_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "23984yks"
	const testErrorMessage = "an error occured"

	mockClient := media.NewMockClient(ctrl)
	mockClient.EXPECT().GetContent(testReferenceID).Return("", nil, errors.New(testErrorMessage))

	handler := NewDownloadHandler(mockClient)
	router := mux.NewRouter()
	router.Handle("/{referenceID}", handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

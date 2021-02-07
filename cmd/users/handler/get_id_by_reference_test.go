package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/users/mock/repository"
	repo "github.com/reecerussell/open-social/cmd/users/repository"
)

func TestGetIDByReferenceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "923643"
	testUserID := 1023

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetIDByReference(gomock.Any(), testReferenceID).Return(&testUserID, nil)

	handler := NewGetIDByReferenceHandler(mockRepo)

	router := mux.NewRouter()
	router.Handle("/{referenceId}", handler).Methods(http.MethodGet)

	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"id\":%d}\n", testUserID)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestGetIDByReference_RepoReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "923643"
	const testErrorMessage = "an error occured"

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetIDByReference(gomock.Any(), testReferenceID).Return(nil, errors.New(testErrorMessage))

	handler := NewGetIDByReferenceHandler(mockRepo)

	router := mux.NewRouter()
	router.Handle("/{referenceId}", handler).Methods(http.MethodGet)

	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

func TestGetIDByReference_UserNotFound_ReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testReferenceID = "923643"
	const testErrorMessage = "an error occured"

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetIDByReference(gomock.Any(), testReferenceID).Return(nil, repo.ErrUserNotFound)

	handler := NewGetIDByReferenceHandler(mockRepo)

	router := mux.NewRouter()
	router.Handle("/{referenceId}", handler).Methods(http.MethodGet)

	req, _ := http.NewRequest(http.MethodGet, "/"+testReferenceID, nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%v\"}\n", repo.ErrUserNotFound)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "application/json", rr.HeaderMap.Get("Content-Type"))
}

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

	"github.com/reecerussell/open-social/cmd/media/mock/repository"
	"github.com/reecerussell/open-social/cmd/media/model"
	"github.com/reecerussell/open-social/mock/media"
)

func TestCreateMediaHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testID = 123
	const testReferenceID = "247"

	mockRepo := repository.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, m *model.Media) (func(bool), error) {
			m.SetID(testID)
			m.SetReferenceID(testReferenceID)

			return func(ok bool) {}, nil
		})

	mockUploader := media.NewMockService(ctrl)
	mockUploader.EXPECT().Upload(gomock.Any(), testReferenceID, gomock.Any()).Return(nil)

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg"}`
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo, mockUploader)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"id\":%d,\"referenceId\":\"%s\"}\n", testID, testReferenceID)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestCreateMediaHandler_GivenInvalidData_ReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockMediaRepository(ctrl)

	rr := httptest.NewRecorder()

	body := `` // invalid body
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo, nil)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"contentType is a required field\"}\n")
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestCreateMediaHandler_CreateReturnsError_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testErrorMessage = "an error occured"

	mockRepo := repository.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return(func(ok bool) {}, errors.New(testErrorMessage))

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg"}`
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo, nil)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestCreateMediaHandler_MediaUploadFails_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testID = 123
	const testReferenceID = "247"
	const testErrorMessage = "an error occured"

	mockRepo := repository.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, m *model.Media) (func(bool), error) {
			m.SetID(testID)
			m.SetReferenceID(testReferenceID)

			return func(ok bool) {
				assert.False(t, ok)
			}, nil
		})

	mockUploader := media.NewMockService(ctrl)
	mockUploader.EXPECT().Upload(gomock.Any(), testReferenceID, gomock.Any()).
		Return(errors.New(testErrorMessage))

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg"}`
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo, mockUploader)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestCreateMediaHandler_MediaContentIsInvalid_ReturnsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testID = 123
	const testReferenceID = "247"
	const testErrorMessage = "an error occured"

	mockRepo := repository.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, m *model.Media) (func(bool), error) {
			m.SetID(testID)
			m.SetReferenceID(testReferenceID)

			return func(ok bool) {
				assert.False(t, ok)
			}, nil
		})

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg","content":"320+ 3jflsd"}` // invalid base64 content
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo, nil)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

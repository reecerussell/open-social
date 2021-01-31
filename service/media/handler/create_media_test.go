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

	"github.com/reecerussell/open-social/service/media/mock/repository"
	"github.com/reecerussell/open-social/service/media/model"
)

func TestCreateMediaHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testID = 123

	mockRepo := repository.NewMockMediaRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, m *model.Media) error {
			m.SetID(testID)

			return nil
		})

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg"}`
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"id\":%d}\n", testID)
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

	handler := NewCreateMediaHandler(mockRepo)
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
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New(testErrorMessage))

	rr := httptest.NewRecorder()

	body := `{"contentType":"image/jpeg"}`
	req, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(body))

	handler := NewCreateMediaHandler(mockRepo)
	handler.ServeHTTP(rr, req)

	data := make([]byte, rr.Body.Len())
	rr.Body.Read(data)

	exp := fmt.Sprintf("{\"message\":\"%s\"}\n", testErrorMessage)
	assert.Equal(t, exp, string(data))
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

}

package media

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/client/mock"
)

func TestNew(t *testing.T) {
	c := New("http://test.io")
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*mediaClient).base)
}

func TestCreate_GivenValidData_ReturnsSuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateRequest{
		ContentType: "text/plain",
		Content:     "Hello World",
	}

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/media", testInput, gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			resp := respDest.(*CreateResponse)
			resp.ID = 123
			resp.ReferenceID = "923640234"

			return nil
		})

	c := &mediaClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.NoError(t, err)
	assert.Equal(t, 123, resp.ID)
	assert.Equal(t, "923640234", resp.ReferenceID)
}

func TestCreate_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateRequest{
		ContentType: "text/plain",
		Content:     "Hello World",
	}
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/media", testInput, gomock.Any()).Return(testError)

	c := &mediaClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.Nil(t, resp)
	assert.Equal(t, testError, err)
}

func TestGetContent_GivenValidReference_ReturnsContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "19263"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/media/content/"+testReferenceID, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := map[string]interface{}{
				"contentType": "text/plain",
				"content":     "SGVsbG8gV29ybGQ=",
			}
			*(respDest.(*map[string]interface{})) = resp

			return nil
		})

	c := &mediaClient{base: mockHTTP}

	contentType, content, err := c.GetContent(testReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
	assert.Equal(t, "Hello World", string(content))
}

func TestGetContent_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "19263"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/media/content/"+testReferenceID, gomock.Any()).Return(testError)

	c := &mediaClient{base: mockHTTP}

	contentType, content, err := c.GetContent(testReferenceID)
	assert.Empty(t, contentType)
	assert.Nil(t, content)
	assert.Equal(t, testError, err)
}

func TestGetContent_ReturnsInvalidBase64_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "19263"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/media/content/"+testReferenceID, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := map[string]interface{}{
				"contentType": "text/plain",
				"content":     "Hello World",
			}
			*(respDest.(*map[string]interface{})) = resp

			return nil
		})

	c := &mediaClient{base: mockHTTP}

	contentType, content, err := c.GetContent(testReferenceID)
	assert.Empty(t, contentType)
	assert.Nil(t, content)
	assert.Equal(t, "media: server responed with invalid content", err.Error())
}

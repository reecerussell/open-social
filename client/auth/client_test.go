package auth

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
	assert.NotNil(t, c.(*authClient).base)
}

func TestGenerateToken_GivenValidDate_ReturnsSuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &GenerateTokenRequest{
		Username: "test",
		Password: "password123",
	}

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/token", testInput, gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			resp := respDest.(*GenerateTokenResponse)
			resp.Token = "<access token>"
			resp.Expires = 12324

			return nil
		})

	c := &authClient{base: mockHTTP}

	resp, err := c.GenerateToken(testInput)
	assert.NoError(t, err)
	assert.Equal(t, "<access token>", resp.Token)
	assert.Equal(t, int64(12324), resp.Expires)
}

func TestGenerateToken_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &GenerateTokenRequest{
		Username: "test",
		Password: "password123",
	}
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/token", testInput, gomock.Any()).Return(testError)

	c := &authClient{base: mockHTTP}

	resp, err := c.GenerateToken(testInput)
	assert.Nil(t, resp)
	assert.Equal(t, testError, err)
}

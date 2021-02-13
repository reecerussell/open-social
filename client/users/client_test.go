package users

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/client/mock"
)

func TestNew(t *testing.T) {
	c := New("http://test.io")
	assert.NotNil(t, c)
	assert.NotNil(t, c.(*usersClient).base)
}

func TestCreate_GivenValidData_ReturnsSuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateUserRequest{
		Username: "test",
		Password: "password123",
	}

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/users", testInput, gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			resp := respDest.(*CreateUserResponse)
			resp.Username = "test"
			resp.ReferenceID = "923640234"

			return nil
		})

	c := &usersClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.NoError(t, err)
	assert.Equal(t, "test", resp.Username)
	assert.Equal(t, "923640234", resp.ReferenceID)
}

func TestCreate_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateUserRequest{
		Username: "test",
		Password: "password123",
	}
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/users", testInput, gomock.Any()).Return(testError)

	c := &usersClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.Nil(t, resp)
	assert.Equal(t, testError, err)
}

func TestGetClaims_GivenValidData_ReturnsClaimsSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &GetClaimsRequest{
		Username: "test",
		Password: "password123",
	}

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/claims", testInput, gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			resp := respDest.(*GetClaimsResponse)
			resp.Claims = map[string]interface{}{
				"foo": "bar",
			}

			return nil
		})

	c := &usersClient{base: mockHTTP}

	resp, err := c.GetClaims(testInput)
	assert.NoError(t, err)
	assert.Equal(t, "bar", resp.Claims["foo"])
}

func TestGetClaims_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &GetClaimsRequest{
		Username: "test",
		Password: "password123",
	}
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/claims", testInput, gomock.Any()).Return(testError)

	c := &usersClient{base: mockHTTP}

	resp, err := c.GetClaims(testInput)
	assert.Nil(t, resp)
	assert.Equal(t, testError, err)
}

func TestGetIDByReference_GivenValidReferenceID_ReturnsUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "304324"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/users/id/"+testReferenceID, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := respDest.(*GetIDByReferenceResponse)
			resp.ID = 27

			return nil
		})

	c := &usersClient{base: mockHTTP}

	id, err := c.GetIDByReference(testReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, 27, *id)
}

func TestGetIDByReference_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "304324"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/users/id/"+testReferenceID, gomock.Any()).Return(testError)

	c := &usersClient{base: mockHTTP}

	id, err := c.GetIDByReference(testReferenceID)
	assert.Nil(t, id)
	assert.Equal(t, testError, err)
}

func TestGetProfile_GivenValidData_ReturnsProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testReferenceID := "304324"

	expectedURL := fmt.Sprintf("/profile/%s/%s", testUsername, testReferenceID)
	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := respDest.(*Profile)
			resp.Username = testUsername

			return nil
		})

	c := &usersClient{base: mockHTTP}

	profile, err := c.GetProfile(testUsername, testReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, testUsername, profile.Username)
}

func TestGetProfile_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUsername := "test"
	testReferenceID := "304324"
	testError := errors.New("an error occured")

	expectedURL := fmt.Sprintf("/profile/%s/%s", testUsername, testReferenceID)
	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).Return(testError)

	c := &usersClient{base: mockHTTP}

	profile, err := c.GetProfile(testUsername, testReferenceID)
	assert.Nil(t, profile)
	assert.Equal(t, testError, err)
}

func TestGetInfo_GivenValidData_ReturnsInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "304324"

	expectedURL := fmt.Sprintf("/info/%s", testReferenceID)
	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := respDest.(*Info)
			resp.Username = "test"

			return nil
		})

	c := &usersClient{base: mockHTTP}

	info, err := c.GetInfo(testReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, "test", info.Username)
}

func TestGetInfo_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testReferenceID := "304324"
	testError := errors.New("an error occured")

	expectedURL := fmt.Sprintf("/info/%s", testReferenceID)
	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).Return(testError)

	c := &usersClient{base: mockHTTP}

	info, err := c.GetInfo(testReferenceID)
	assert.Nil(t, info)
	assert.Equal(t, testError, err)
}

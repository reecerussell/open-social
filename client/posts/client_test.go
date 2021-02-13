package posts

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
	assert.NotNil(t, c.(*postsClient).base)
}

func TestCreate_GivenValidData_ReturnsSuccessfulResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateRequest{
		UserReferenceID: "12343543",
		Caption:         "Hello World",
	}

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts", testInput, gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			resp := respDest.(*CreateResponse)
			resp.ReferenceID = "2304734"

			return nil
		})

	c := &postsClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.NoError(t, err)
	assert.Equal(t, "2304734", resp.ReferenceID)
}

func TestGenerateToken_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testInput := &CreateRequest{
		UserReferenceID: "12343543",
		Caption:         "Hello World",
	}
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts", testInput, gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}

	resp, err := c.Create(testInput)
	assert.Nil(t, resp)
	assert.Equal(t, testError, err)
}

func TestGetFeed_GivenValidUserReference_ReturnsFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/feed/"+testUserReferenceID, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := (respDest.(*[]*FeedItem))
			*resp = append(*resp, &FeedItem{
				Caption: "Hello World",
			})

			return nil
		})

	c := &postsClient{base: mockHTTP}

	feedItems, err := c.GetFeed(testUserReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(feedItems))
	assert.Equal(t, "Hello World", feedItems[0].Caption)
}

func TestGetFeed_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Get("/feed/"+testUserReferenceID, gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}

	feedItems, err := c.GetFeed(testUserReferenceID)
	assert.Nil(t, feedItems)
	assert.Equal(t, testError, err)
}

func TestGetProfileFeed_GivenValidUserReference_ReturnsFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testUsername := "test"

	mockHTTP := mock.NewMockHTTP(ctrl)
	expectedURL := fmt.Sprintf("/profile/feed/%s/%s", testUsername, testUserReferenceID)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := (respDest.(*[]*FeedItem))
			*resp = append(*resp, &FeedItem{
				Caption: "Hello World",
			})

			return nil
		})

	c := &postsClient{base: mockHTTP}

	feedItems, err := c.GetProfileFeed(testUsername, testUserReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(feedItems))
	assert.Equal(t, "Hello World", feedItems[0].Caption)
}

func TestGetProfileFeed_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testUsername := "test"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	expectedURL := fmt.Sprintf("/profile/feed/%s/%s", testUsername, testUserReferenceID)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}

	feedItems, err := c.GetProfileFeed(testUsername, testUserReferenceID)
	assert.Nil(t, feedItems)
	assert.Equal(t, testError, err)
}

func TestLikePost_GivenCorrectReferences_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts/like", gomock.Any(), gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			payload := body.(map[string]string)
			assert.Equal(t, testUserReferenceID, payload["userReferenceId"])
			assert.Equal(t, testPostReferenceID, payload["postReferenceId"])

			return nil
		})

	c := &postsClient{base: mockHTTP}
	err := c.LikePost(testPostReferenceID, testUserReferenceID)
	assert.NoError(t, err)
}

func TestLikePost_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts/like", gomock.Any(), gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}
	err := c.LikePost(testPostReferenceID, testUserReferenceID)
	assert.Equal(t, testError, err)
}

func TestUnlikePost_GivenCorrectReferences_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts/unlike", gomock.Any(), gomock.Any()).
		DoAndReturn(func(url string, body, respDest interface{}) error {
			payload := body.(map[string]string)
			assert.Equal(t, testUserReferenceID, payload["userReferenceId"])
			assert.Equal(t, testPostReferenceID, payload["postReferenceId"])

			return nil
		})

	c := &postsClient{base: mockHTTP}
	err := c.UnlikePost(testPostReferenceID, testUserReferenceID)
	assert.NoError(t, err)
}

func TestUnlikePost_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	mockHTTP.EXPECT().Post("/posts/unlike", gomock.Any(), gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}
	err := c.UnlikePost(testPostReferenceID, testUserReferenceID)
	assert.Equal(t, testError, err)
}

func TestGet_GivenValidReferences_ReturnsPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"

	mockHTTP := mock.NewMockHTTP(ctrl)
	expectedURL := fmt.Sprintf("/posts/%s/%s", testPostReferenceID, testUserReferenceID)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).
		DoAndReturn(func(url string, respDest interface{}) error {
			resp := respDest.(*Post)
			resp.Caption = "Hello World"

			return nil
		})

	c := &postsClient{base: mockHTTP}
	post, err := c.Get(testPostReferenceID, testUserReferenceID)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", post.Caption)
}

func TestGet_RequestFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserReferenceID := "2340703470324"
	testPostReferenceID := "3294849323233"
	testError := errors.New("an error occured")

	mockHTTP := mock.NewMockHTTP(ctrl)
	expectedURL := fmt.Sprintf("/posts/%s/%s", testPostReferenceID, testUserReferenceID)
	mockHTTP.EXPECT().Get(expectedURL, gomock.Any()).Return(testError)

	c := &postsClient{base: mockHTTP}
	post, err := c.Get(testPostReferenceID, testUserReferenceID)
	assert.Nil(t, post)
	assert.Equal(t, testError, err)
}

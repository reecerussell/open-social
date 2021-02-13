package posts

import (
	"fmt"

	"github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the posts API.
type Client interface {
	Create(in *CreateRequest) (*CreateResponse, error)
	GetFeed(userReferenceID string) ([]*FeedItem, error)
	GetProfileFeed(username, userReferenceID string) ([]*FeedItem, error)
	LikePost(postReferenceID, userReferenceID string) error
	UnlikePost(postReferenceID, userReferenceID string) error
	Get(postReferenceID, userReferenceID string) (*Post, error)
}

type postsClient struct {
	base client.HTTP
}

// New returns a new instance of the post client.
func New(url string) Client {
	return &postsClient{
		base: client.NewHTTP(url),
	}
}

func (c *postsClient) Create(in *CreateRequest) (*CreateResponse, error) {
	var resp CreateResponse
	err := c.base.Post("/posts", in, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *postsClient) GetFeed(userReferenceID string) ([]*FeedItem, error) {
	var items []*FeedItem
	err := c.base.Get("/feed/"+userReferenceID, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *postsClient) GetProfileFeed(username, userReferenceID string) ([]*FeedItem, error) {
	var items []*FeedItem
	url := fmt.Sprintf("/profile/feed/%s/%s", username, userReferenceID)
	err := c.base.Get(url, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *postsClient) LikePost(postReferenceID, userReferenceID string) error {
	payload := map[string]string{
		"postReferenceId": postReferenceID,
		"userReferenceId": userReferenceID,
	}

	err := c.base.Post("/posts/like", payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *postsClient) UnlikePost(postReferenceID, userReferenceID string) error {
	payload := map[string]string{
		"postReferenceId": postReferenceID,
		"userReferenceId": userReferenceID,
	}

	err := c.base.Post("/posts/unlike", payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *postsClient) Get(postReferenceID, userReferenceID string) (*Post, error) {
	var post Post
	url := fmt.Sprintf("/posts/%s/%s", postReferenceID, userReferenceID)
	err := c.base.Get(url, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

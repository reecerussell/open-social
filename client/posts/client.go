package posts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	clientpkg "github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the posts API.
type Client interface {
	Create(in *CreateRequest) (*CreateResponse, error)
	GetFeed(userReferenceID string) ([]*FeedItem, error)
	LikePost(postReferenceID, userReferenceID string) error
	UnlikePost(postReferenceID, userReferenceID string) error
	Get(postReferenceID, userReferenceID string) (*Post, error)
}

type client struct {
	url  string
	http *http.Client
}

// New returns a new instance of the post client.
func New(url string) Client {
	return &client{
		url: url,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *client) Create(in *CreateRequest) (*CreateResponse, error) {
	jsonBytes, _ := json.Marshal(in)
	body := bytes.NewBuffer(jsonBytes)

	url := c.url + "/posts"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data CreateResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return &data, nil
	}

	var data clientpkg.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return nil, clientpkg.NewError(resp.StatusCode, data.Message)
}

func (c *client) GetFeed(userReferenceID string) ([]*FeedItem, error) {
	url := c.url + "/feed/" + userReferenceID
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data []*FeedItem
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	var data clientpkg.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return nil, clientpkg.NewError(resp.StatusCode, data.Message)
}

func (c *client) LikePost(postReferenceID, userReferenceID string) error {
	payload := map[string]string{
		"postReferenceId": postReferenceID,
		"userReferenceId": userReferenceID,
	}
	jsonBytes, _ := json.Marshal(payload)
	body := bytes.NewReader(jsonBytes)

	url := c.url + "/posts/like"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var data clientpkg.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}

		return clientpkg.NewError(resp.StatusCode, data.Message)
	}

	return nil
}

func (c *client) UnlikePost(postReferenceID, userReferenceID string) error {
	payload := map[string]string{
		"postReferenceId": postReferenceID,
		"userReferenceId": userReferenceID,
	}
	jsonBytes, _ := json.Marshal(payload)
	body := bytes.NewReader(jsonBytes)

	url := c.url + "/posts/unlike"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var data clientpkg.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}

		return clientpkg.NewError(resp.StatusCode, data.Message)
	}

	return nil
}

func (c *client) Get(postReferenceID, userReferenceID string) (*Post, error) {
	url := fmt.Sprintf("%s/posts/%s/%s", c.url, postReferenceID, userReferenceID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var data clientpkg.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return nil, clientpkg.NewError(resp.StatusCode, data.Message)
	}

	var data Post
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

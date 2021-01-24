package posts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	clientpkg "github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the posts API.
type Client interface {
	Create(in *CreateRequest) (*CreateResponse, error)
}

type client struct {
	url  string
	http *http.Client
}

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

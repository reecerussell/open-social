package media

import (
	"encoding/base64"
	"errors"

	"github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the media API.
type Client interface {
	Create(in *CreateRequest) (*CreateResponse, error)
	GetContent(referenceID string) (string, []byte, error)
}

type mediaClient struct {
	base client.HTTP
}

// New returns a new instance of Client.
func New(url string) Client {
	return &mediaClient{
		base: client.NewHTTP(url),
	}
}

func (c *mediaClient) Create(in *CreateRequest) (*CreateResponse, error) {
	var resp CreateResponse
	err := c.base.Post("/media", in, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *mediaClient) GetContent(referenceID string) (string, []byte, error) {
	var data map[string]interface{}
	err := c.base.Get("/media/content/"+referenceID, &data)
	if err != nil {
		return "", nil, err
	}

	content, err := base64.StdEncoding.DecodeString(data["content"].(string))
	if err != nil {
		return "", nil, errors.New("media: server responed with invalid content")
	}

	return data["contentType"].(string), content, nil
}

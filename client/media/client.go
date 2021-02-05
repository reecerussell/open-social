package media

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	clientpkg "github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the media API.
type Client interface {
	Create(in *CreateRequest) (*CreateResponse, error)
	GetContent(referenceID string) (string, []byte, error)
}

type client struct {
	url  string
	http *http.Client
}

// New returns a new instance of Client.
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

	url := c.url + "/media"
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

func (c *client) GetContent(referenceID string) (string, []byte, error) {
	url := c.url + "/media/content/" + referenceID
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var data clientpkg.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", nil, err
		}

		return "", nil, clientpkg.NewError(resp.StatusCode, data.Message)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", nil, err
	}

	content, err := base64.StdEncoding.DecodeString(data["content"].(string))
	if err != nil {
		return "", nil, err
	}

	return data["contentType"].(string), content, nil
}

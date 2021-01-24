package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	clientpkg "github.com/reecerussell/open-social/client"
)

// Client is a interface to the auth API.
type Client interface {
	GenerateToken(in *GenerateTokenRequest) (*GenerateTokenResponse, error)
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
			Timeout: time.Second * 10,
		},
	}
}

func (c *client) GenerateToken(in *GenerateTokenRequest) (*GenerateTokenResponse, error) {
	jsonBytes, _ := json.Marshal(in)
	body := bytes.NewBuffer(jsonBytes)

	url := c.url + "/token"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data GenerateTokenResponse
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

package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Client is an interface used to interact with the users API.
type Client interface {
	Create(in *CreateUserRequest) (*CreateUserResponse, error)
	GetClaims(in *GetClaimsRequest) (*GetClaimsResponse, error)
	GetIDByReference(referenceID string) (*int, error)
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

type client struct {
	url  string
	http *http.Client
}

func (c *client) Create(in *CreateUserRequest) (*CreateUserResponse, error) {
	jsonBytes, _ := json.Marshal(in)
	body := bytes.NewBuffer(jsonBytes)

	url := c.url + "/users"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data CreateUserResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return &data, nil
	}

	var data ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return nil, NewError(resp.StatusCode, data.Message)
}

func (c *client) GetClaims(in *GetClaimsRequest) (*GetClaimsResponse, error) {
	jsonBytes, _ := json.Marshal(in)
	body := bytes.NewBuffer(jsonBytes)

	url := c.url + "/claims"
	req, _ := http.NewRequest(http.MethodPost, url, body)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data GetClaimsResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return &data, nil
	}

	var data ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return nil, NewError(resp.StatusCode, data.Message)
}

func (c *client) GetIDByReference(referenceID string) (*int, error) {
	url := c.url + "/users/id/" + referenceID
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data GetIDByReferenceResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		return &data.ID, nil
	}

	var data ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return nil, NewError(resp.StatusCode, data.Message)
}

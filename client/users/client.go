package users

import (
	"fmt"

	"github.com/reecerussell/open-social/client"
)

// Client is an interface used to interact with the users API.
type Client interface {
	Create(in *CreateUserRequest) (*CreateUserResponse, error)
	GetClaims(in *GetClaimsRequest) (*GetClaimsResponse, error)
	GetIDByReference(referenceID string) (*int, error)
	GetProfile(username, userReferenceID string) (*Profile, error)
	GetInfo(userReferenceID string) (*Info, error)
}

// New returns a new instance of Client.
func New(url string) Client {
	return &usersClient{
		base: client.NewHTTP(url),
	}
}

type usersClient struct {
	base client.HTTP
}

func (c *usersClient) Create(in *CreateUserRequest) (*CreateUserResponse, error) {
	var resp CreateUserResponse
	err := c.base.Post("/users", in, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *usersClient) GetClaims(in *GetClaimsRequest) (*GetClaimsResponse, error) {
	var resp GetClaimsResponse
	err := c.base.Post("/claims", in, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *usersClient) GetIDByReference(referenceID string) (*int, error) {
	var resp GetIDByReferenceResponse
	err := c.base.Get("/users/id/"+referenceID, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.ID, nil
}

func (c *usersClient) GetProfile(username, userReferenceID string) (*Profile, error) {
	var profile Profile
	url := fmt.Sprintf("/profile/%s/%s", username, userReferenceID)
	err := c.base.Get(url, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (c *usersClient) GetInfo(userReferenceID string) (*Info, error) {
	var info Info
	url := fmt.Sprintf("/info/%s", userReferenceID)
	err := c.base.Get(url, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

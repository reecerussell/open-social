package auth

import (
	"github.com/reecerussell/open-social/client"
)

// Client is a interface to the auth API.
type Client interface {
	GenerateToken(in *GenerateTokenRequest) (*GenerateTokenResponse, error)
}

type authClient struct {
	base client.HTTP
}

// New returns a new instance of Client.
func New(url string) Client {
	return &authClient{
		base: client.NewHTTP(url),
	}
}

func (c *authClient) GenerateToken(in *GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var resp GenerateTokenResponse
	err := c.base.Post("/token", in, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

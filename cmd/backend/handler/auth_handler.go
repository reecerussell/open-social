package handler

import (
	"encoding/json"
	"log"
	"net/http"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client"
	"github.com/reecerussell/open-social/client/auth"
	"github.com/reecerussell/open-social/client/users"
)

// AuthHandler provides functions to handle authentication of
// the platform, including registrations and logging in.
type AuthHandler struct {
	core.Handler
	users users.Client
	auth  auth.Client
}

// NewAuthHandler returns a new instance of AuthHandler, taking
// the following parameters as dependencies.
func NewAuthHandler(users users.Client, auth auth.Client) *AuthHandler {
	return &AuthHandler{
		users: users,
		auth:  auth,
	}
}

// RegisterUserResponse is the response body of the register request.
type RegisterUserResponse struct {
	ReferenceID string             `json:"referenceId"`
	Username    string             `json:"username"`
	AccessToken *RegisterUserToken `json:"accessToken"`
}

// RegisterUserToken represents a user's access token in the register user response body.
type RegisterUserToken struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// Register handles requests to register a user.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var data users.CreateUserRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	user, err := h.users.Create(&data)
	if err != nil {
		switch e := err.(type) {
		case *client.Error:
			h.RespondError(w, e, e.StatusCode)
			return
		default:
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	response := RegisterUserResponse{
		ReferenceID: user.ReferenceID,
		Username:    user.Username,
	}

	token, err := h.auth.GenerateToken(&auth.GenerateTokenRequest{
		Username: user.Username,
		Password: data.Password,
	})
	if err == nil {
		response.AccessToken = &RegisterUserToken{
			Token:   token.Token,
			Expires: token.Expires,
		}
	} else {
		log.Printf("WARN: failed to generate token: %v\n", err)
	}

	h.Respond(w, response)
}

// TokenRequest is a type used to unmarshal a token request's body to.
type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Token is a http.HandlerFunc used to provider an access token for
// a user, with the credentials given by the TokenRequest body.
func (h *AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	var data TokenRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	token, err := h.auth.GenerateToken(&auth.GenerateTokenRequest{
		Username: data.Username,
		Password: data.Password,
	})
	if err != nil {
		switch e := err.(type) {
		case *client.Error:
			h.RespondError(w, e, e.StatusCode)
			return
		default:
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	h.Respond(w, token)
}

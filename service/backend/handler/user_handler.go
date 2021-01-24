package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/reecerussell/open-social/client/auth"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
)

// UserHandler handles requests to the user domain.
type UserHandler struct {
	core.Handler
	client users.Client
	auth   auth.Client
}

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(client users.Client, auth auth.Client) *UserHandler {
	return &UserHandler{client: client, auth: auth}
}

type RegisterUserResponse struct {
	ReferenceID string             `json:"referenceId"`
	Username    string             `json:"username"`
	AccessToken *RegisterUserToken `json:"accessToken"`
}

type RegisterUserToken struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// Register handles requests to register a user.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var data users.CreateUserRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	user, err := h.client.Create(&data)
	if err != nil {
		switch e := err.(type) {
		case *users.Error:
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

	// TODO: redirect to authenticate and provide token for ui
	h.Respond(w, response)
}

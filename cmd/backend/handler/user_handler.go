package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client"
	"github.com/reecerussell/open-social/client/auth"
	"github.com/reecerussell/open-social/client/posts"
	"github.com/reecerussell/open-social/client/users"
)

// UserHandler handles requests to the user domain.
type UserHandler struct {
	core.Handler
	client users.Client
	auth   auth.Client
	posts  posts.Client
}

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(client users.Client, auth auth.Client, posts posts.Client) *UserHandler {
	return &UserHandler{
		client: client,
		auth:   auth,
		posts:  posts,
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
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var data users.CreateUserRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	user, err := h.client.Create(&data)
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

	// TODO: redirect to authenticate and provide token for ui
	h.Respond(w, response)
}

// GetProfileResponse returns a user's profile data.
type GetProfileResponse struct {
	users.Profile
	Feed []*posts.FeedItem `json:"feed"`
}

// GetProfile handles requests to get a user's profile.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	profile, err := h.client.GetProfile(username, userID)
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

	feed, err := h.posts.GetProfileFeed(username, userID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		switch e := err.(type) {
		case *client.Error:
			h.RespondError(w, e, e.StatusCode)
			return
		default:
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	resp := GetProfileResponse{
		Profile: *profile,
		Feed:    feed,
	}

	h.Respond(w, resp)
}

// GetInfo handles requests to get information for a user.
func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	info, err := h.client.GetInfo(userID)
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

	h.Respond(w, info)
}

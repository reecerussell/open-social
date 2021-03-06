package handler

import (
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

// Follow handles requests to make the current user follow the user with the given id.
func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userReferenceID := params["userReferenceID"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	err := h.client.Follow(userReferenceID, userID)
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

	h.Respond(w, nil)
}

// Unfollow handles requests to make the current user unfollow the user with the given id.
func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userReferenceID := params["userReferenceID"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	err := h.client.Unfollow(userReferenceID, userID)
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

	h.Respond(w, nil)
}

package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/posts/provider"
)

// ProfileFeedHandler is a http.Handler used to request a user's profile feed.
type ProfileFeedHandler struct {
	core.Handler
	provider provider.PostProvider
}

// NewProfileFeedHandler returns a new instance of ProfileFeedHandler.
func NewProfileFeedHandler(provider provider.PostProvider) *ProfileFeedHandler {
	return &ProfileFeedHandler{
		provider: provider,
	}
}

func (h *ProfileFeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	userReferenceID, err := uuid.Parse(params["userReferenceID"])
	if err != nil {
		h.RespondError(w, fmt.Errorf("user reference id must be a valid guid"), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	feed, err := h.provider.GetProfileFeed(ctx, username, userReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	h.Respond(w, feed)
}

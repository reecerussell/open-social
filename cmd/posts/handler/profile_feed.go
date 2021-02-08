package handler

import (
	"net/http"

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
	userReferenceID := params["userReferenceId"]

	ctx := r.Context()
	feed, err := h.provider.GetProfileFeed(ctx, username, userReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	h.Respond(w, feed)
}

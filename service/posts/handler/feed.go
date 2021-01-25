package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/posts/repository"
)

// FeedHandler is a http.Handler used to request a user's post feed.
type FeedHandler struct {
	core.Handler
	repo repository.PostRepository
}

// NewFeedHandler returns a new instance of FeedHandler.
func NewFeedHandler(repo repository.PostRepository) *FeedHandler {
	return &FeedHandler{repo: repo}
}

func (h *FeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userReferenceID := params["userReferenceId"]

	ctx := r.Context()
	feed, err := h.repo.GetFeed(ctx, userReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	h.Respond(w, feed)
}

package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/posts/provider"
)

// GetPostHandler is a http.Handler used to get a post for a user.
type GetPostHandler struct {
	core.Handler
	provider provider.PostProvider
}

// NewGetPostHandler returns a new instance of GetPostHandler.
func NewGetPostHandler(provider provider.PostProvider) *GetPostHandler {
	return &GetPostHandler{
		provider: provider,
	}
}

func (h *GetPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postReferenceID := params["postReferenceID"]
	userReferenceID := params["userReferenceID"]

	ctx := r.Context()
	post, err := h.provider.Get(ctx, postReferenceID, userReferenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == provider.ErrPostNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	h.Respond(w, post)
}

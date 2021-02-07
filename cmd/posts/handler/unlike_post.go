package handler

import (
	"encoding/json"
	"net/http"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/posts/repository"
)

// UnlikePostHandler is a http.Handler used to mark a post as unliked by a user.
type UnlikePostHandler struct {
	core.Handler
	repo  repository.PostRepository
	likes repository.LikeRepository
}

// UnlikePostRequest is the request body structure.
type UnlikePostRequest struct {
	PostReferenceID string `json:"postReferenceId"`
	UserReferenceID string `json:"userReferenceId"`
}

// NewUnlikePostHandler returns a new instance of UnlikePostHandler.
func NewUnlikePostHandler(repo repository.PostRepository, likes repository.LikeRepository) *UnlikePostHandler {
	return &UnlikePostHandler{
		repo:  repo,
		likes: likes,
	}
}

func (h *UnlikePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data UnlikePostRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	ctx := r.Context()
	post, err := h.repo.Get(ctx, data.PostReferenceID, data.UserReferenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrPostNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	err = post.CanUnlike()
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	err = h.likes.Delete(ctx, post.ID(), data.UserReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	h.Respond(w, nil)
}

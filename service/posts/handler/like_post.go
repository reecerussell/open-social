package handler

import (
	"encoding/json"
	"net/http"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/service/posts/repository"
)

// LikePostHandler is a http.Handler used to mark a post as liked by a user.
type LikePostHandler struct {
	core.Handler
	repo  repository.PostRepository
	likes repository.LikeRepository
}

// LikePostRequest is the request body structure.
type LikePostRequest struct {
	PostReferenceID string `json:"postReferenceId"`
	UserReferenceID string `json:"userReferenceId"`
}

// NewLikePostHandler returns a new instance of LikePostHandler.
func NewLikePostHandler(repo repository.PostRepository, likes repository.LikeRepository) *LikePostHandler {
	return &LikePostHandler{
		repo:  repo,
		likes: likes,
	}
}

func (h *LikePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data LikePostRequest
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

	err = post.CanLike()
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	err = h.likes.Create(ctx, post.ID(), data.UserReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	h.Respond(w, nil)
}

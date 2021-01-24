package handler

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/posts/model"
	"github.com/reecerussell/open-social/service/posts/repository"
)

type CreatePost struct {
	core.Handler
	repo  repository.PostRepository
	users users.Client
}

type CreatePostRequest struct {
	UserReferenceID string `json:"userReferenceId"`
	Caption         string `json:"caption"`
}

type CreatePostResponse struct {
	ReferenceID string `json:"referenceId"`
}

func NewCreatePost(repo repository.PostRepository, users users.Client) *CreatePost {
	return &CreatePost{
		repo:  repo,
		users: users,
	}
}

func (h *CreatePost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data CreatePostRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	userID, err := h.users.GetIDByReference(data.UserReferenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	post, err := model.NewPost(*userID, data.Caption)
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	response := CreatePostResponse{
		ReferenceID: post.ReferenceID(),
	}

	h.Respond(w, response)
}

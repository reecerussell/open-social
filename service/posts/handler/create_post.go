package handler

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/posts/model"
	"github.com/reecerussell/open-social/service/posts/repository"
)

// CreatePostHandler is a http.Handler used to handle POSt requests to create post records.
type CreatePostHandler struct {
	core.Handler
	repo  repository.PostRepository
	users users.Client
}

// CreatePostRequest is the body of a request.
type CreatePostRequest struct {
	UserReferenceID string `json:"userReferenceId"`
	Caption         string `json:"caption"`
}

// CreatePostResponse is the body of the response.
type CreatePostResponse struct {
	ReferenceID string `json:"referenceId"`
}

// NewCreatePostHandler returns a new instance of CreatePostHandler.
func NewCreatePostHandler(repo repository.PostRepository, users users.Client) *CreatePostHandler {
	return &CreatePostHandler{
		repo:  repo,
		users: users,
	}
}

func (h *CreatePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	err = h.repo.Create(r.Context(), post)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	response := CreatePostResponse{
		ReferenceID: post.ReferenceID(),
	}

	h.Respond(w, response)
}

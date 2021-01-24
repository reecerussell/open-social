package handler

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/open-social/client/posts"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
)

// PostHandler handles requests to the post domain.
type PostHandler struct {
	core.Handler
	client posts.Client
}

// NewPostHandler returns a new instance of PostHandler.
func NewPostHandler(client posts.Client) *PostHandler {
	return &PostHandler{client: client}
}

// CreatePostRequest is the body of the request.
type CreatePostRequest struct {
	Caption string `json:"caption"`
}

// CreatePostResponse contains the reference id of the newly created post.
type CreatePostResponse struct {
	ID string `json:"id"`
}

// Create handles requests to create a post.
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var data CreatePostRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	post, err := h.client.Create(&posts.CreateRequest{
		UserReferenceID: userID,
		Caption:         data.Caption,
	})
	if err != nil {
		switch e := err.(type) {
		case *users.Error:
			h.RespondError(w, e, e.StatusCode)
			return
		default:
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	response := CreatePostResponse{
		ID: post.ReferenceID,
	}

	h.Respond(w, response)
}

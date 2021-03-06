package handler

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client"
	"github.com/reecerussell/open-social/client/media"
	"github.com/reecerussell/open-social/client/posts"
)

// PostHandler handles requests to the post domain.
type PostHandler struct {
	core.Handler
	client posts.Client
	media  media.Client
}

// NewPostHandler returns a new instance of PostHandler.
func NewPostHandler(client posts.Client, media media.Client) *PostHandler {
	return &PostHandler{
		client: client,
		media:  media,
	}
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
	r.ParseMultipartForm(10 << 20)

	ctx := r.Context()
	mediaID, success := h.uploadMedia(ctx, w, r)
	if !success {
		return
	}

	userID := ctx.Value(core.ContextKey("uid")).(string)
	caption := r.FormValue("caption")

	post, err := h.client.Create(&posts.CreateRequest{
		UserReferenceID: userID,
		MediaID:         mediaID,
		Caption:         caption,
	})
	if err != nil {
		// TODO: queue media deletion
		h.handleError(w, err)
		return
	}

	response := CreatePostResponse{
		ID: post.ReferenceID,
	}

	h.Respond(w, response)
}

func (h *PostHandler) uploadMedia(ctx context.Context, w http.ResponseWriter, r *http.Request) (*int, bool) {
	file, header, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		h.RespondError(w, err, http.StatusInternalServerError)
		return nil, false
	}
	defer file.Close()

	fileData := make([]byte, header.Size)
	file.Read(fileData)
	contentType := http.DetectContentType(fileData)

	m, err := h.media.Create(&media.CreateRequest{
		ContentType: contentType,
		Content:     base64.StdEncoding.EncodeToString(fileData),
	})
	if err != nil {
		h.handleError(w, err)
		return nil, false
	}

	return &m.ID, true
}

// GetFeed returns a user's feed.
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	feed, err := h.client.GetFeed(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.Respond(w, feed)
}

// GetPost returns a post.
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	post, err := h.client.Get(id, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.Respond(w, post)
}

// Like marks a post as liked by the current user.
func (h *PostHandler) Like(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	err := h.client.LikePost(id, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.Respond(w, nil)
}

// Unlike marks a post as unliked by the current user.
func (h *PostHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	err := h.client.UnlikePost(id, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.Respond(w, nil)
}

func (h *PostHandler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *client.Error:
		h.RespondError(w, e, e.StatusCode)
		return
	default:
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}
}

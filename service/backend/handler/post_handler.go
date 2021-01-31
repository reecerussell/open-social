package handler

import (
	"net/http"

	"github.com/reecerussell/open-social/client"
	mediaSdk "github.com/reecerussell/open-social/client/media"
	"github.com/reecerussell/open-social/client/posts"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/core/media"
)

// PostHandler handles requests to the post domain.
type PostHandler struct {
	core.Handler
	client   posts.Client
	uploader media.Service
	media    mediaSdk.Client
}

// NewPostHandler returns a new instance of PostHandler.
func NewPostHandler(client posts.Client, uploader media.Service, media mediaSdk.Client) *PostHandler {
	return &PostHandler{
		client:   client,
		uploader: uploader,
		media:    media,
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

	file, header, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	var mediaID *int

	ctx := r.Context()
	if err != http.ErrMissingFile {
		defer file.Close()

		fileData := make([]byte, header.Size)
		file.Read(fileData)
		contentType := http.DetectContentType(fileData)

		m, err := h.media.Create(&mediaSdk.CreateRequest{
			ContentType: contentType,
		})
		if err != nil {
			switch e := err.(type) {
			case *client.Error:
				h.RespondError(w, e, e.StatusCode)
				return
			default:
				h.RespondError(w, err, http.StatusInternalServerError)
				return
			}
		}

		mediaID = &m.ID

		err = h.uploader.Upload(ctx, m.ReferenceID, fileData)
		if err != nil {
			// TODO: delete new media record @ m.ID
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	userID := ctx.Value(core.ContextKey("uid")).(string)

	post, err := h.client.Create(&posts.CreateRequest{
		UserReferenceID: userID,
		MediaID:         mediaID,
		Caption:         r.FormValue("caption"),
	})
	if err != nil {
		switch e := err.(type) {
		case *client.Error:
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

// GetFeed returns a user's feed.
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(core.ContextKey("uid")).(string)

	feed, err := h.client.GetFeed(userID)
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

	h.Respond(w, feed)
}

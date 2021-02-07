package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/media"
	"github.com/reecerussell/open-social/service/media/model"
	"github.com/reecerussell/open-social/service/media/repository"
)

// CreateMediaHandler is a http.Handler used to create a new Media record.
type CreateMediaHandler struct {
	core.Handler
	repo     repository.MediaRepository
	uploader media.Service
}

// CreateMediaRequest is the body of the request.
type CreateMediaRequest struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

// CreateMediaResponse is the body of the response.
type CreateMediaResponse struct {
	ID          int    `json:"id"`
	ReferenceID string `json:"referenceId"`
}

// NewCreateMediaHandler returns a new instance of CreateMediaHandler.
func NewCreateMediaHandler(repo repository.MediaRepository, uploader media.Service) *CreateMediaHandler {
	return &CreateMediaHandler{
		repo:     repo,
		uploader: uploader,
	}
}

func (h *CreateMediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data CreateMediaRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	media, err := model.NewMedia(data.ContentType)
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	save, err := h.repo.Create(ctx, media)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.upload(ctx, media.ReferenceID(), data.Content)
	if err != nil {
		save(false)
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	// Commit media
	save(true)

	resp := CreateMediaResponse{
		ID:          media.ID(),
		ReferenceID: media.ReferenceID(),
	}

	h.Respond(w, resp)
}

func (h *CreateMediaHandler) upload(ctx context.Context, key, content string) error {
	bytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		log.Printf("ERROR: failed to decode media base64 content: %v\n", err)
		return err
	}

	err = h.uploader.Upload(ctx, key, bytes)
	if err != nil {
		log.Printf("ERROR: failed to upload: %v\n", err)
		return err
	}

	return nil
}

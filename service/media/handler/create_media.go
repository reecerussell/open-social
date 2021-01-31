package handler

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/media/model"
	"github.com/reecerussell/open-social/service/media/repository"
)

// CreateMediaHandler is a http.Handler used to create a new Media record.
type CreateMediaHandler struct {
	core.Handler
	repo repository.MediaRepository
}

// CreateMediaRequest is the body of the request.
type CreateMediaRequest struct {
	ContentType string `json:"contentType"`
}

// CreateMediaResponse is the body of the response.
type CreateMediaResponse struct {
	ID          int    `json:"id"`
	ReferenceID string `json:"referenceId"`
}

// NewCreateMediaHandler returns a new instance of CreateMediaHandler.
func NewCreateMediaHandler(repo repository.MediaRepository) *CreateMediaHandler {
	return &CreateMediaHandler{repo: repo}
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
	err = h.repo.Create(ctx, media)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	resp := CreateMediaResponse{
		ID:          media.ID(),
		ReferenceID: media.ReferenceID(),
	}

	h.Respond(w, resp)
}

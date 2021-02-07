package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/media"
	"github.com/reecerussell/open-social/service/media/repository"
)

// GetMediaContentHandler is a http.Handler which serves a media's content.
type GetMediaContentHandler struct {
	core.Handler
	repo       repository.MediaRepository
	downloader media.Service
}

// GetMediaContentResponse represents the response body of the request.
type GetMediaContentResponse struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

// NewGetMediaContentHandler returns a new instance of GetMediaContentHandler.
func NewGetMediaContentHandler(repo repository.MediaRepository, downloader media.Service) *GetMediaContentHandler {
	return &GetMediaContentHandler{
		repo:       repo,
		downloader: downloader,
	}
}

func (h *GetMediaContentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	referenceID := params["referenceID"]

	ctx := r.Context()
	contentType, err := h.repo.GetContentType(ctx, referenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrMediaNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	data, err := h.downloader.Download(ctx, referenceID)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	resp := GetMediaContentResponse{
		ContentType: contentType,
		Content:     b64,
	}

	h.Respond(w, resp)
}

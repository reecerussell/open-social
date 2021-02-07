package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/media"
)

// DownloadHandler is a http.Handler used to download media.
type DownloadHandler struct {
	core.Handler
	client media.Client
}

// NewDownloadHandler returns a new instance of DownloadHandler.
func NewDownloadHandler(client media.Client) *DownloadHandler {
	return &DownloadHandler{
		client: client,
	}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	referenceID := params["referenceID"]

	contentType, content, err := h.client.GetContent(referenceID)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=3600")
	w.Write(content)
}

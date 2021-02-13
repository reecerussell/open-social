package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/provider"
)

// GetInfoHandler is a http.Handler used to serve a user's info.
type GetInfoHandler struct {
	core.Handler
	provider provider.UserProvider
}

// NewGetInfoHandler returns a new instance of GetInfoHandler.
func NewGetInfoHandler(provider provider.UserProvider) *GetInfoHandler {
	return &GetInfoHandler{
		provider: provider,
	}
}

func (h *GetInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userReferenceID := params["userReferenceID"]

	ctx := r.Context()
	profile, err := h.provider.GetInfo(ctx, userReferenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == provider.ErrProfileNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	h.Respond(w, profile)
}

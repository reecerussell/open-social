package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/provider"
)

// GetProfileHandler is a http.Handler used to serve a user's profile.
type GetProfileHandler struct {
	core.Handler
	provider provider.UserProvider
}

// NewGetProfileHandler returns a new instance of GetProfileHandler.
func NewGetProfileHandler(provider provider.UserProvider) *GetProfileHandler {
	return &GetProfileHandler{
		provider: provider,
	}
}

func (h *GetProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	userReferenceID := params["userReferenceID"]

	ctx := r.Context()
	profile, err := h.provider.GetProfile(ctx, username, userReferenceID)
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

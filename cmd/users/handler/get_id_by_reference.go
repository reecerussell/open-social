package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/repository"
)

// GetIDByReferenceHandler is a http.Handler which handles GET requests
// to retrieve a user's id, given their reference id.
type GetIDByReferenceHandler struct {
	core.Handler
	repo repository.UserRepository
}

// GetIDByReferenceResponse contains the returned user id.
type GetIDByReferenceResponse struct {
	ID int `json:"id"`
}

// NewGetIDByReferenceHandler returns a new instance of GetIDByReferenceHandler.
func NewGetIDByReferenceHandler(repo repository.UserRepository) *GetIDByReferenceHandler {
	return &GetIDByReferenceHandler{repo: repo}
}

func (h *GetIDByReferenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	referenceID := params["referenceId"]

	userID, err := h.repo.GetIDByReference(r.Context(), referenceID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			h.RespondError(w, err, http.StatusNotFound)
			return
		}

		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	response := GetIDByReferenceResponse{
		ID: *userID,
	}

	h.Respond(w, response)
}

package handler

import (
	"encoding/json"
	"net/http"

	hashpkg "github.com/reecerussell/adaptive-password-hasher"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/service/users/repository"
)

// GetClaimsHandler is a http.Handler used to get a user's claim values,
// given a username and password.
type GetClaimsHandler struct {
	core.Handler

	hasher hashpkg.Hasher
	repo   repository.UserRepository
}

// NewGetClaimsHandler returns a new instance of GetClaimsHandler,
// with the given dependcies.
func NewGetClaimsHandler(hasher hashpkg.Hasher, repo repository.UserRepository) *GetClaimsHandler {
	return &GetClaimsHandler{
		hasher: hasher,
		repo:   repo,
	}
}

// GetClaimsRequest represents the request body.
type GetClaimsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetClaimsResponse represents the response body.
type GetClaimsResponse struct {
	Claims map[string]interface{} `json:"claims"`
}

// ServeHTTP handles HTTP requests to get a user's claims.
func (h *GetClaimsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data GetClaimsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	ctx := r.Context()
	user, err := h.repo.GetUserByUsername(ctx, data.Username)
	if err != nil {
		if err == repository.ErrUserNotFound {
			h.RespondError(w, err, http.StatusBadRequest)
			return
		}

		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	err = user.VerifyPassword(data.Password, h.hasher)
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	claims := map[string]interface{}{
		"username": user.Username(),
		"uid":      user.ReferenceID(),
	}

	resp := GetClaimsResponse{
		Claims: claims,
	}

	h.Respond(w, resp)
}

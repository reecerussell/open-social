package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/repository"
)

// FollowUserHandler is a http.Handler used to handle requests to create a follow record.
type FollowUserHandler struct {
	core.Handler
	repo      repository.UserRepository
	followers repository.FollowerRepository
}

// NewFollowUserHandler returns a new instance of FollowUserHandler.
func NewFollowUserHandler(repo repository.UserRepository, followers repository.FollowerRepository) *FollowUserHandler {
	return &FollowUserHandler{
		repo:      repo,
		followers: followers,
	}
}

func (h *FollowUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userReferenceID := params["userReferenceId"]
	followerReferenceID := params["followerReferenceId"]

	ctx := r.Context()
	user, err := h.repo.GetUserByReference(ctx, userReferenceID, followerReferenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrUserNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	err = user.CanFollow()
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	err = h.followers.Create(ctx, user.ID(), followerReferenceID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrFollowerNotFound {
			status = http.StatusNotFound
		}

		h.RespondError(w, err, status)
		return
	}

	h.Respond(w, nil)
}

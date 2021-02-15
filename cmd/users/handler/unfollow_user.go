package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/repository"
)

// UnfollowUserHandler is a http.Handler used to handle requests to delete a follow record.
type UnfollowUserHandler struct {
	core.Handler
	repo      repository.UserRepository
	followers repository.FollowerRepository
}

// NewUnfollowUserHandler returns a new instance of UnfollowUserHandler.
func NewUnfollowUserHandler(repo repository.UserRepository, followers repository.FollowerRepository) *UnfollowUserHandler {
	return &UnfollowUserHandler{
		repo:      repo,
		followers: followers,
	}
}

func (h *UnfollowUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	err = user.CanUnfollow()
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	err = h.followers.Delete(ctx, user.ID(), followerReferenceID)
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

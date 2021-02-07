package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	hashpkg "github.com/reecerussell/adaptive-password-hasher"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/service/users/model"
	"github.com/reecerussell/open-social/service/users/password"
	"github.com/reecerussell/open-social/service/users/repository"
)

// CreateUserHandler is a http.Handler which handles requests to create users.
type CreateUserHandler struct {
	core.Handler

	val    password.Validator
	hasher hashpkg.Hasher
	repo   repository.UserRepository
}

// NewCreateUserHandler returns a new instance of CreateUserHandler,
// with the given dependencies.
func NewCreateUserHandler(val password.Validator, hasher hashpkg.Hasher, repo repository.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{
		val:    val,
		hasher: hasher,
		repo:   repo,
	}
}

// CreateUserRequest represents the request body.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUserResponse represents the response body.
type CreateUserResponse struct {
	ReferenceID string `json:"referenceId"`
	Username    string `json:"username"`
}

// ServeHTTP handles HTTP requests to create users.
func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data CreateUserRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	user, err := model.NewUser(data.Username, data.Password, h.val, h.hasher)
	if err != nil {
		h.RespondError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exists, err := h.repo.DoesUsernameExist(ctx, user.Username(), nil)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	if exists {
		h.RespondError(w, fmt.Errorf("the username '%s' is taken", user.Username()), http.StatusBadRequest)
		return
	}

	err = h.repo.Create(ctx, user)
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	resp := CreateUserResponse{
		ReferenceID: user.ReferenceID(),
		Username:    user.Username(),
	}

	h.Respond(w, resp)
}

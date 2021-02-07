package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/reecerussell/gojwt"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/users"
)

// TokenHandler handles HTTP POST requests to generate an access token.
type TokenHandler struct {
	core.Handler
	client        users.Client
	alg           gojwt.Algorithm
	expiryMinutes int
}

// TokenResponse contains a generated access token and it's expiry date.
type TokenResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// NewTokenHandler returns a new instance of TokenHandler.
func NewTokenHandler(client users.Client, alg gojwt.Algorithm, expiryMinutes int) *TokenHandler {
	return &TokenHandler{
		client:        client,
		alg:           alg,
		expiryMinutes: expiryMinutes,
	}
}

// ServeHTTP handles requests to generate an access token.
func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data users.GetClaimsRequest
	_ = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	claims, err := h.client.GetClaims(&data)
	if err != nil {
		switch e := err.(type) {
		case *users.Error:
			h.RespondError(w, e, e.StatusCode)
			return
		default:
			h.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	now := time.Now().UTC()
	exp := now.Add(time.Duration(h.expiryMinutes) * time.Minute)

	builder := gojwt.New(h.alg).
		AddClaims(claims.Claims).
		SetExpiry(exp)

	token, err := builder.Build()
	if err != nil {
		h.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		Token:   token,
		Expires: exp.Unix(),
	}

	h.Respond(w, response)
}

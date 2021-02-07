package core

import (
	"encoding/json"
	"net/http"
)

// Handler is a abstract handler which provides helpers to
// implementing handlers.
type Handler struct{}

// Respond writes an OK status code, as well as writing data
// to the response as JSON.
func (*Handler) Respond(w http.ResponseWriter, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

// RespondError writes an error to the response with the given status code.
func (*Handler) RespondError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{Message: err.Error()}
	_ = json.NewEncoder(w).Encode(resp)
}

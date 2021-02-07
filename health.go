package core

import "net/http"

// HealthCheckHandler is a basic health check handler.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

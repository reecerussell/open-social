package core

import (
	"context"
	"net/http"
)

// HealthCheck is used to check the status of a service.
type HealthCheck interface {
	Check(ctx context.Context) error
}

// HealthCheckHandler returns a new http.Handler which handles health checks,
// using the given health checks.
func HealthCheckHandler(healthChecks []HealthCheck) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, hc := range healthChecks {
			if err := hc.Check(ctx); err != nil {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte("Unhealthy"))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})
}

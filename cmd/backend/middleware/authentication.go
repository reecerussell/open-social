package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/reecerussell/gojwt"

	core "github.com/reecerussell/open-social"
)

var allowedPaths = []string{"/auth/register", "/auth/token", "/health"}

// Authentication is middleware used to authenticate HTTP requests.
type Authentication struct {
	core.Handler
	alg gojwt.Algorithm
}

// NewAuthentication returns a new instance of Authentication.
func NewAuthentication(alg gojwt.Algorithm) *Authentication {
	return &Authentication{alg: alg}
}

// Handle returns a new http.Handler, used to authenticate the given handler.
func (m *Authentication) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAllowedPath(r) {
			h.ServeHTTP(w, r)
			return
		}

		token, err := getAccessToken(r)
		if err != nil {
			m.RespondError(w, err, http.StatusUnauthorized)
			return
		}

		jwt, err := gojwt.Token(token)
		if err != nil {
			m.RespondError(w, err, http.StatusUnauthorized)
			return
		}

		err = jwt.Verify(m.alg)
		if err != nil {
			m.RespondError(w, err, http.StatusUnauthorized)
			return
		}

		ctx := contextWithClaims(r.Context(), jwt.Claims)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func isAllowedPath(r *http.Request) bool {
	for _, allowedPath := range allowedPaths {
		if r.URL.Path == allowedPath {
			return true
		}
	}

	return false
}

func getAccessToken(r *http.Request) (string, error) {
	value := r.Header.Get("Authorization")
	if value == "" {
		return "", errors.New("no auth header present")
	}

	if len(value) < 6 || value[:6] != "Bearer" {
		return "", errors.New("invalid auth scheme")
	}

	return value[7:], nil
}

func contextWithClaims(ctx context.Context, claims gojwt.Claims) context.Context {
	for k, v := range claims {
		ctx = context.WithValue(ctx, core.ContextKey(k), v)
	}

	return ctx
}

package core

import (
	"log"
	"net/http"
)

// Middleware is a http.Handler used as an intermediate
// step before handling a base handler.
type Middleware interface {
	Handle(h http.Handler) http.Handler
}

// ChainMiddleware chains a list of middleware into one http.Handler.
func ChainMiddleware(base http.Handler, middleware ...Middleware) http.Handler {
	h := base

	for i := len(middleware) - 1; i > 0; i-- {
		h = middleware[i].Handle(h)
	}

	return h
}

// LoggingMiddleware is middleware which logs all incoming requests.
type LoggingMiddleware struct{}

// NewLoggingMiddleware returns a new instance of LoggingMiddleware.
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Handle returns a new http.Handler which logs all requests.
func (*LoggingMiddleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s\n", r.Method, r.RequestURI)

		h.ServeHTTP(w, r)
	})
}

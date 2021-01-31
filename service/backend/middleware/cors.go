package middleware

import "net/http"

// Cors is middleware used to add CORS headers to the HTTP response.
type Cors struct{}

// NewCors returns a new instance of Cors.
func NewCors() *Cors {
	return &Cors{}
}

// Handle returns a new http.Handler which will add the CORS headers.
func (*Cors) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, r)
	})
}

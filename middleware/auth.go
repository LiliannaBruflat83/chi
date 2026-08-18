// Package middleware provides HTTP middleware that halts the handler chain
// after writing a response, preventing superfluous WriteHeader warnings and
// unintended execution of downstream handlers.
package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware rejects requests without a valid bearer token. It writes the
// 401 response and returns immediately so next.ServeHTTP is never called for
// unauthorized requests.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // halt the chain
		}
		next.ServeHTTP(w, r)
	})
}

// isAuthorized reports whether the request carries a valid bearer token.
func isAuthorized(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	return strings.HasPrefix(auth, "Bearer ") && strings.TrimPrefix(auth, "Bearer ") != ""
}

package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware creates a middleware that checks for authorization.
// If unauthorized, it writes http.StatusUnauthorized and immediately halts the middleware chain.
func AuthMiddleware(isAuthorized func(r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAuthorized != nil && !isAuthorized(r) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return // Explicit early return halts middleware chain execution
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireHeader ensures that a required header is present and non-empty.
// If missing, writes http.StatusBadRequest and returns immediately.
func RequireHeader(headerName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.TrimSpace(r.Header.Get(headerName)) == "" {
				http.Error(w, "Missing required header: "+headerName, http.StatusBadRequest)
				return // Explicit early return halts middleware chain execution
			}
			next.ServeHTTP(w, r)
		})
	}
}

// BasicAuth creates a basic authentication middleware.
// If credentials don't match or are missing, returns StatusUnauthorized with WWW-Authenticate header and halts execution.
func BasicAuth(realm string, credentials map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return // Explicit early return halts middleware chain execution
			}

			expectedPass, userExists := credentials[user]
			if !userExists || expectedPass != pass {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return // Explicit early return halts middleware chain execution
			}

			next.ServeHTTP(w, r)
		})
	}
}

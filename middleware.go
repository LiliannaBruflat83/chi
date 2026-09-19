package main

import (
	"net/http"
)

// AuthMiddleware demonstrates the fix for chi middleware chain execution after http.Error.
// Before fix: calling http.Error without return would fall through to next.ServeHTTP, causing
// "superfluous response.WriteHeader call", unintended downstream execution, and potential panics.
// After fix: return immediately after http.Error, preventing downstream handlers from running.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example auth check — replace with real logic
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // FIX: halt chain, do not call next.ServeHTTP
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware shows the same pattern for any middleware that writes a response
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: reject bad requests early
		if r.Header.Get("X-Blocked") == "1" {
			http.Error(w, "Blocked", http.StatusForbidden)
			return // FIX: prevent downstream execution
		}
		next.ServeHTTP(w, r)
	})
}

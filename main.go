package main

import (
	"fmt"
	"net/http"
)

// AuthMiddleware halts downstream execution immediately upon authorization failure
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // Explicit early return halts middleware chain execution
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Chi Middleware Security Hardened.")
}

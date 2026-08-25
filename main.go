package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/LiliannaBruflat83/chi/middleware"
)

func main() {
	fmt.Println("Running Chi Middleware Verification Example...")

	authMw := middleware.AuthMiddleware(func(r *http.Request) bool {
		return r.Header.Get("Authorization") == "Bearer secret-token"
	})

	handler := authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Access Granted"))
	}))

	// Test unauthorized request
	req1 := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	fmt.Printf("Unauthorized Request -> Status: %d, Body: %s\n", rec1.Code, rec1.Body.String())

	// Test authorized request
	req2 := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req2.Header.Set("Authorization", "Bearer secret-token")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	fmt.Printf("Authorized Request -> Status: %d, Body: %s\n", rec2.Code, rec2.Body.String())
}

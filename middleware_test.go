package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_BlocksDownstream(t *testing.T) {
	downstreamCalled := false
	downstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamCalled = true
		w.Write([]byte("downstream"))
	})

	handler := AuthMiddleware(downstream)

	req := httptest.NewRequest("GET", "/", nil)
	// No Authorization header -> should block
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if downstreamCalled {
		t.Fatal("downstream handler was called after http.Error — fix missing return")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	// Ensure no superfluous WriteHeader: second write should not panic
}

func TestAuthMiddleware_AllowsAuthorized(t *testing.T) {
	downstreamCalled := false
	downstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handler := AuthMiddleware(downstream)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if !downstreamCalled {
		t.Fatal("downstream should be called for authorized request")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

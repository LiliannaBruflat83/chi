package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthMiddlewareRejectsUnauthorized verifies that when a request is
// rejected, the downstream handler is completely bypassed and the correct
// 401 response is written.
func TestAuthMiddlewareRejectsUnauthorized(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if got := rec.Body.String(); got != "Unauthorized\n" {
		t.Fatalf("expected body %q, got %q", "Unauthorized\n", got)
	}
	if nextCalled {
		t.Fatal("downstream handler was called for an unauthorized request")
	}
}

// TestAuthMiddlewareAllowsAuthorized verifies that authorized requests pass
// through to the downstream handler.
func TestAuthMiddlewareAllowsAuthorized(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !nextCalled {
		t.Fatal("downstream handler was not called for an authorized request")
	}
}

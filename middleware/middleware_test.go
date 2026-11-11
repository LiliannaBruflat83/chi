package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthMiddlewareHaltsChainOnUnauthorized(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := AuthMiddleware(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Unauthorized") {
		t.Fatalf("expected Unauthorized body, got %q", rec.Body.String())
	}
	if nextCalled {
		t.Fatalf("expected nextHandler NOT to be called after http.Error")
	}
}

func TestAPIKeyMiddlewareHaltsChainOnForbidden(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := APIKeyMiddleware("secret-key-123")(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
	if nextCalled {
		t.Fatalf("expected nextHandler NOT to be called on invalid API key")
	}
}

func TestContentTypeJSONMiddlewareHaltsChainOnInvalidContentType(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := ContentTypeJSONMiddleware(nextHandler)
	req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader("plain text"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status %d, got %d", http.StatusUnsupportedMediaType, rec.Code)
	}
	if nextCalled {
		t.Fatalf("expected nextHandler NOT to be called on invalid Content-Type")
	}
}

func TestValidRequestProceedsToNextHandler(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := AuthMiddleware(ContentTypeJSONMiddleware(nextHandler))
	req := httptest.NewRequest(http.MethodPost, "/api/items", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !nextCalled {
		t.Fatalf("expected nextHandler to be called on valid request")
	}
}

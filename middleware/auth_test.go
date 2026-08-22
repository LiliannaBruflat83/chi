package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_Unauthorized(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mw := AuthMiddleware(func(r *http.Request) bool {
		return false // Reject request
	})

	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if nextCalled {
		t.Errorf("expected downstream handler to be bypassed, but nextHandler was called")
	}
}

func TestAuthMiddleware_Authorized(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mw := AuthMiddleware(func(r *http.Request) bool {
		return true // Allow request
	})

	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !nextCalled {
		t.Errorf("expected downstream handler to be called, but nextCalled was false")
	}
}

func TestRequireHeader_MissingHeader(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	mw := RequireHeader("X-API-Key")
	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if nextCalled {
		t.Errorf("expected downstream handler to be bypassed when header is missing")
	}
}

func TestRequireHeader_PresentHeader(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	mw := RequireHeader("X-API-Key")
	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	req.Header.Set("X-API-Key", "secret-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if nextCalled == false {
		t.Errorf("expected downstream handler to be called when header is present")
	}
}

func TestBasicAuth_Unauthorized(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	creds := map[string]string{
		"admin": "password123",
	}
	mw := BasicAuth("Restricted", creds)
	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if nextCalled {
		t.Errorf("expected downstream handler to be bypassed on failed Basic Auth")
	}
}

func TestBasicAuth_Authorized(t *testing.T) {
	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	creds := map[string]string{
		"admin": "password123",
	}
	mw := BasicAuth("Restricted", creds)
	handler := mw(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.SetBasicAuth("admin", "password123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Errorf("expected downstream handler to be called on successful Basic Auth")
	}
}

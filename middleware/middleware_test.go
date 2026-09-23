package middleware

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestAuthMiddlewareBlocksAndStopsChain(t *testing.T) {
    nextCalled := false
    handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        nextCalled = true
    }))

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)

    handler.ServeHTTP(rr, req)

    if nextCalled {
        t.Fatalf("next handler should not be called on unauthorized requests")
    }
    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected status %d got %d", http.StatusUnauthorized, rr.Code)
    }
    if !strings.Contains(rr.Body.String(), "Unauthorized") {
        t.Fatalf("expected unauthorized body, got %q", rr.Body.String())
    }
}

func TestBlockHeaderMiddlewareStopsChain(t *testing.T) {
    nextCalled := false
    handler := BlockHeaderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        nextCalled = true
    }))

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("X-Block", "true")

    handler.ServeHTTP(rr, req)

    if nextCalled {
        t.Fatalf("next handler should not be called when X-Block header present")
    }
    if rr.Code != http.StatusForbidden {
        t.Fatalf("expected status %d got %d", http.StatusForbidden, rr.Code)
    }
    if rr.Body.String() != "blocked" {
        t.Fatalf("expected body %q got %q", "blocked", rr.Body.String())
    }
}

func TestAuthMiddlewareAllowsThrough(t *testing.T) {
    nextCalled := false
    handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        nextCalled = true
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    }))

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("Authorization", "token")

    handler.ServeHTTP(rr, req)

    if !nextCalled {
        t.Fatalf("next handler should have been called for authorized request")
    }
    if rr.Code != http.StatusOK {
        t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
    }
    if rr.Body.String() != "ok" {
        t.Fatalf("expected body %q got %q", "ok", rr.Body.String())
    }
}

package middleware

import (
	"net/http"
	"strings"
)

// ResponseOnceWriter wraps http.ResponseWriter to track whether headers/body
// have already been committed and prevent superfluous WriteHeader calls.
type ResponseOnceWriter struct {
	http.ResponseWriter
	wroteHeader bool
	statusCode  int
}

func NewResponseOnceWriter(w http.ResponseWriter) *ResponseOnceWriter {
	return &ResponseOnceWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *ResponseOnceWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseOnceWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseOnceWriter) Committed() bool {
	return rw.wroteHeader
}

// AuthMiddleware validates the Authorization Bearer token and returns immediately
// after calling http.Error so next.ServeHTTP is never executed on unauthorized requests.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") || len(strings.TrimPrefix(auth, "Bearer ")) == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// APIKeyMiddleware validates the X-API-Key header against a required key and
// halts the chain immediately after writing a 403 Forbidden response.
func APIKeyMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expectedKey == "" || r.Header.Get("X-API-Key") != expectedKey {
				http.Error(w, "Forbidden: invalid API key", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ContentTypeJSONMiddleware enforces application/json on POST/PUT/PATCH requests
// and returns immediately after writing 415 Unsupported Media Type.
func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if !strings.HasPrefix(ct, "application/json") {
				http.Error(w, "Unsupported Media Type: application/json required", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RecovererMiddleware catches panics, writes 500 Internal Server Error only if
// the response has not yet been committed, and returns immediately.
func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		safeW := NewResponseOnceWriter(w)
		defer func() {
			if rec := recover(); rec != nil {
				if !safeW.Committed() {
					http.Error(safeW, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
		}()
		next.ServeHTTP(safeW, r)
	})
}

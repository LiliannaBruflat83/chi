package main

import (
	"fmt"
	"net/http"

	"github.com/LiliannaBruflat83/chi/middleware"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	var h http.Handler = mux
	h = middleware.ContentTypeJSONMiddleware(h)
	h = middleware.AuthMiddleware(h)
	h = middleware.RecovererMiddleware(h)
	return h
}

func main() {
	fmt.Println("Chi Middleware Server Ready")
}

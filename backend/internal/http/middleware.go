package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return middleware.RequestID(middleware.Recoverer(next))
}

func JSONOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type")[:16] != "application/json" {
				http.Error(w, `{"error":"unsupported_media_type"}`, http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func LogStartup(addr string) {
	log.Printf("API listening on http://localhost%s\n", addr)
}

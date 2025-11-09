// UTF-8
package http

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"desafio/internal/store"
)

func csvEnv(key string, fallback []string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func NewRouter(s store.TaskStore) http.Handler {
	r := chi.NewRouter()

	// CORS: read from env (CSV). Default: localhost:5173
	allowedOrigins := csvEnv("ALLOWED_ORIGINS", []string{"http://localhost:5173"})
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Optional: strip trailing slashes for route leniency
	r.Use(StripSlashes())

	// Health/root
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Tasks routes
	th := TaskHandlers{Store: s}
	r.Mount("/tasks", th.Routes())

	return r
}

// StripSlashes removes a single trailing slash ("/") from the request path.
func StripSlashes() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := r.URL.Path
			if len(p) > 1 && strings.HasSuffix(p, "/") {
				r.URL.Path = strings.TrimRight(p, "/")
			}
			next.ServeHTTP(w, r)
		})
	}
}

// (Example) reasonable server timeouts can be used where needed.
// Keeping here if you decide to build your own *http.Server in this package.
var (
	readHeaderTimeout = 5 * time.Second
)

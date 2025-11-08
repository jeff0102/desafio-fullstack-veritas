package httpx

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"desafio/internal/store"
)

func NewRouter(s store.TaskStore) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(CommonMiddleware)
	r.Use(JSONOnly)
	r.Use(middleware.StripSlashes) // tolerate trailing/double slashes

	// CORS
	allowed := os.Getenv("ALLOWED_ORIGINS")
	if strings.TrimSpace(allowed) == "" {
		allowed = "http://localhost:5173"
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowed},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	// Tasks routes
	th := TaskHandlers{Store: s}
	r.Mount("/tasks", th.Routes())

	return r
}

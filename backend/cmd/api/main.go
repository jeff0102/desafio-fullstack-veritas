package main

import (
"encoding/json"
"log"
"net/http"
"time"

"github.com/go-chi/chi/v5"
"github.com/go-chi/cors"
)

type Task struct {
ID          string    `json:"id"`
Title       string    `json:"title"`
Description string    `json:"description"`
Status      string    `json:"status"` // todo | doing | done
CreatedAt   time.Time `json:"createdAt"`
UpdatedAt   time.Time `json:"updatedAt"`
}

func main() {
r := chi.NewRouter()

// CORS for the Vite frontend (http://localhost:5173)
r.Use(cors.Handler(cors.Options{
AllowedOrigins:   []string{"http://localhost:5173"},
AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
ExposedHeaders:   []string{"Link"},
AllowCredentials: false,
MaxAge:           300,
}))

// Seed tasks so the frontend shows dummy data
now := time.Now()
seed := []Task{
	{ID: "t1", Title: "Set up environment", Status: "todo", CreatedAt: now, UpdatedAt: now},
	{ID: "t2", Title: "Backend skeleton", Status: "doing", CreatedAt: now, UpdatedAt: now},
	{ID: "t3", Title: "Frontend skeleton", Status: "done", CreatedAt: now, UpdatedAt: now},
}

// GET /tasks
r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json; charset=utf-8")
_ = json.NewEncoder(w).Encode(seed)
})

srv := &http.Server{
Addr:              ":8080",
Handler:           r,
ReadHeaderTimeout: 5 * time.Second,
}

log.Println("API listening on http://localhost:8080")
log.Fatal(srv.ListenAndServe())
}

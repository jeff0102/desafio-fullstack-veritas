package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"desafio/internal/core"
	httpx "desafio/internal/http"
	"desafio/internal/store"
)

func main() {
	// Configuration
	jsonPath := os.Getenv("TASKS_JSON_PATH")
	if jsonPath == "" {
		jsonPath = "tasks.json" // default path; with "go run -C backend ..." this lives in /backend
	}

	// Seed (fallback)
	now := time.Now()
	seed := []core.Task{
		{ID: "t1", Title: "Set up environment", Status: "todo", CreatedAt: now, UpdatedAt: now},
		{ID: "t2", Title: "Backend skeleton", Status: "doing", CreatedAt: now, UpdatedAt: now},
		{ID: "t3", Title: "Frontend skeleton", Status: "done", CreatedAt: now, UpdatedAt: now},
	}

	// Try to load persisted data
	loaded, err := store.LoadFromJSON(jsonPath)
	if err != nil {
		log.Printf("WARN: could not load %s: %v (using seed)", jsonPath, err)
	}

	var initial []core.Task
	if len(loaded) > 0 {
		initial = loaded
	} else {
		initial = seed
	}

	// Build store: in-memory + persistence wrapper
	mem := store.NewMemoryStore(initial)
	persisted := store.NewPersistedStore(mem, jsonPath)

	// Build router
	r := httpx.NewRouter(persisted)

	addr := ":" + getPort()
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	httpx.LogStartup(addr)
	log.Printf("Persistence: %s", jsonPath)
	log.Fatal(srv.ListenAndServe())
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// UTF-8
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"desafio/internal/core"
	httpx "desafio/internal/http"
	"desafio/internal/store"
)

// envOr returns the value of key or fallback if blank.
func envOr(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func main() {
	// --- config by env ---
	tasksPath := envOr("TASKS_JSON_PATH", "tasks.json")
	addr := envOr("PORT", ":8080")
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	// --- load existing tasks (if file exists) ---
	var initial []core.Task
	if tasks, err := store.LoadFromJSON(tasksPath); err != nil {
		log.Fatalf("failed to load tasks from %s: %v", tasksPath, err)
	} else {
		initial = tasks // may be nil/empty -> starts blank
	}

	// --- wire store (memory seeded with initial), then wrap with persistence ---
	mem := store.NewMemoryStore(initial)
	persisted := store.NewPersistedStore(mem, tasksPath)

	// --- router (CORS is handled inside router using env ALLOWED_ORIGINS) ---
	r := httpx.NewRouter(persisted)

	// --- server ---
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("API listening on http://localhost%s", addr)
	log.Printf("Persistence file: %s", tasksPath)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

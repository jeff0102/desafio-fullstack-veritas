package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"desafio/internal/core"
	"desafio/internal/store"

	"github.com/go-chi/chi/v5"
)

type TaskHandlers struct {
	Store store.TaskStore
}

func (h TaskHandlers) Routes() chi.Router {
	r := chi.NewRouter()
	// /tasks
	r.Get("/", h.List)
	r.Post("/", h.Create)
	// /tasks/{id}
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	// /tasks/{id}/reorder
	r.Put("/{id}/reorder", h.Reorder)
	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

// GET /tasks[?status=todo|doing|done]
func (h TaskHandlers) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "list_failed")
		return
	}

	// Filter by status if provided
	if raw := strings.TrimSpace(r.URL.Query().Get("status")); raw != "" {
		switch raw {
		case "todo", "doing", "done":
			filtered := make([]core.Task, 0, len(items))
			for _, t := range items {
				if t.Status == raw {
					filtered = append(filtered, t)
				}
			}
			items = filtered
		default:
			writeErr(w, http.StatusBadRequest, "invalid_status")
			return
		}
	}

	// Sort by order asc, then updatedAt desc (stable by id if equal)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Order == items[j].Order {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		return items[i].Order < items[j].Order
	})

	writeJSON(w, http.StatusOK, items)
}

func (h TaskHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.Store.Get(r.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			writeErr(w, http.StatusNotFound, "not_found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "get_failed")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h TaskHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var in core.NewTask
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if err := core.ValidateNew(in); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.Store.Create(context.Background(), in)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "create_failed")
		return
	}
	w.Header().Set("Location", "/tasks/"+item.ID)
	writeJSON(w, http.StatusCreated, item)
}

func (h TaskHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in core.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if err := core.ValidateUpdate(in); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.Store.Update(r.Context(), id, in)
	if err != nil {
		if err == store.ErrNotFound {
			writeErr(w, http.StatusNotFound, "not_found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "update_failed")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h TaskHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.Store.Delete(r.Context(), id); err != nil {
		if err == store.ErrNotFound {
			writeErr(w, http.StatusNotFound, "not_found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "delete_failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type reorderReq struct {
	Status string `json:"status"`
	Index  int    `json:"index"` // 0-based
}

// PUT /tasks/{id}/reorder
func (h TaskHandlers) Reorder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req reorderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json")
		return
	}
	switch req.Status {
	case "todo", "doing", "done":
	default:
		writeErr(w, http.StatusBadRequest, "invalid_status")
		return
	}
	item, err := h.Store.Reorder(r.Context(), id, req.Status, req.Index)
	if err != nil {
		if err == store.ErrNotFound {
			writeErr(w, http.StatusNotFound, "not_found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "reorder_failed")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

package httpx

import (
	"context"
	"encoding/json"
	"net/http"

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

func (h TaskHandlers) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.List(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "list_failed")
		return
	}
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

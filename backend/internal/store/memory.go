package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"desafio/internal/core"
)

var (
	ErrNotFound = errors.New("not_found")
)

type TaskStore interface {
	List(ctx context.Context) ([]core.Task, error)
	Get(ctx context.Context, id string) (core.Task, error)
	Create(ctx context.Context, in core.NewTask) (core.Task, error)
	Update(ctx context.Context, id string, in core.UpdateTask) (core.Task, error)
	Delete(ctx context.Context, id string) error
}

type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]core.Task
}

func NewMemoryStore(seed []core.Task) *MemoryStore {
	m := &MemoryStore{
		tasks: make(map[string]core.Task),
	}
	for _, t := range seed {
		m.tasks[t.ID] = t
	}
	return m
}

func (m *MemoryStore) List(ctx context.Context) ([]core.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]core.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		out = append(out, t)
	}
	return out, nil
}

func (m *MemoryStore) Get(ctx context.Context, id string) (core.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tasks[id]
	if !ok {
		return core.Task{}, ErrNotFound
	}
	return t, nil
}

func (m *MemoryStore) Create(ctx context.Context, in core.NewTask) (core.Task, error) {
	now := time.Now()
	t := core.Task{
		ID:        genID(),
		Title:     in.Title,
		Status:    "todo",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if in.Description != nil {
		t.Description = *in.Description
	}

	m.mu.Lock()
	m.tasks[t.ID] = t
	m.mu.Unlock()

	return t, nil
}

func (m *MemoryStore) Update(ctx context.Context, id string, in core.UpdateTask) (core.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	if !ok {
		return core.Task{}, ErrNotFound
	}
	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.Status != nil {
		t.Status = *in.Status
	}
	t.UpdatedAt = time.Now()
	m.tasks[id] = t
	return t, nil
}

func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}

// genID creates a short, unique ID (base32-like) without extra deps.
func genID() string {
	// Timestamp + counter fallback. For simplicity and no external deps.
	// Good enough for this challenge; can be replaced by uuid later.
	return time.Now().Format("20060102150405.000000000")
}

package store

import (
	"context"
	"errors"
	"sort"
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
	Reorder(ctx context.Context, id string, newStatus string, index int) (core.Task, error)
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
	// Ensure order is initialized per status
	m.normalizeAll()
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
	// new task goes to "todo" at the bottom of that column
	status := "todo"
	maxOrder := m.maxOrderLocked(status)

	t := core.Task{
		ID:        genID(),
		Title:     in.Title,
		Status:    status,
		Order:     maxOrder + 1,
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
	if in.Status != nil && t.Status != *in.Status {
		// Move to another column: append to bottom there; source column will be normalized lazily on List or by Reorder later
		t.Status = *in.Status
		t.Order = m.maxOrderLocked(t.Status) + 1
	}
	t.UpdatedAt = time.Now()
	m.tasks[t.ID] = t
	return t, nil
}

func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	if !ok {
		return ErrNotFound
	}
	delete(m.tasks, id)
	// reindex source column for compactness
	m.reindexStatusLocked(t.Status)
	return nil
}

// Reorder moves a task to `newStatus` and inserts it at `index` (0-based), reindexing that column.
func (m *MemoryStore) Reorder(ctx context.Context, id string, newStatus string, index int) (core.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return core.Task{}, ErrNotFound
	}

	// Collect destination column tasks (excluding the moving one)
	dst := m.byStatusLocked(newStatus, id)

	// Bound index
	if index < 0 {
		index = 0
	}
	if index > len(dst) {
		index = len(dst)
	}

	// Insert placeholder at index
	dst = append(dst, core.Task{})
	copy(dst[index+1:], dst[index:])
	dst[index] = t

	// If status changed, update and reindex both columns; else just reindex one
	if t.Status != newStatus {
		oldStatus := t.Status
		// Update item in place
		dst[index].Status = newStatus
		dst[index].UpdatedAt = time.Now()
		// Persist into map with temporary zero order; reindex below
		m.tasks[id] = dst[index]

		// Reindex both columns
		m.reindexStatusLocked(oldStatus)
		m.reindexSliceLocked(newStatus, dst)
	} else {
		// Same column reorder
		m.reindexSliceLocked(newStatus, dst)
	}

	return m.tasks[id], nil
}

// --- helpers (locked) ---

func (m *MemoryStore) byStatusLocked(status string, excludeID string) []core.Task {
	var items []core.Task
	for _, t := range m.tasks {
		if t.Status == status && t.ID != excludeID {
			items = append(items, t)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Order == items[j].Order {
			return items[i].UpdatedAt.Before(items[j].UpdatedAt)
		}
		return items[i].Order < items[j].Order
	})
	return items
}

func (m *MemoryStore) reindexStatusLocked(status string) {
	var items []core.Task
	for _, t := range m.tasks {
		if t.Status == status {
			items = append(items, t)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Order == items[j].Order {
			return items[i].UpdatedAt.Before(items[j].UpdatedAt)
		}
		return items[i].Order < items[j].Order
	})
	for i := range items {
		items[i].Order = i + 1
		items[i].UpdatedAt = time.Now()
		m.tasks[items[i].ID] = items[i]
	}
}

func (m *MemoryStore) reindexSliceLocked(status string, items []core.Task) {
	for i := range items {
		items[i].Status = status
		items[i].Order = i + 1
		items[i].UpdatedAt = time.Now()
		m.tasks[items[i].ID] = items[i]
	}
}

func (m *MemoryStore) normalizeAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range []string{"todo", "doing", "done"} {
		m.reindexStatusLocked(s)
	}
}

func (m *MemoryStore) maxOrderLocked(status string) int {
	max := 0
	for _, t := range m.tasks {
		if t.Status == status && t.Order > max {
			max = t.Order
		}
	}
	return max
}

// genID creates a short, unique ID without extra deps.
func genID() string {
	return time.Now().Format("20060102150405.000000000")
}

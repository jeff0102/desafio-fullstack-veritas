package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"desafio/internal/core"
)

// LoadFromJSON loads tasks from a JSON file. If the file does not exist, returns (nil, nil).
func LoadFromJSON(path string) ([]core.Task, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // no data yet
		}
		return nil, err
	}
	var tasks []core.Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// SaveToJSON writes tasks to a JSON file (pretty-printed).
// It ensures the directory exists and writes to a temp file, then renames.
func SaveToJSON(path string, tasks []core.Task) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "tasks-*.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tasks); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	// Replace target atomically where possible; on Windows we remove first.
	_ = os.Remove(path)
	return os.Rename(tmp.Name(), path)
}

// PersistedStore wraps a TaskStore and writes to disk on mutating operations.
type PersistedStore struct {
	inner TaskStore
	path  string
}

func NewPersistedStore(inner TaskStore, path string) *PersistedStore {
	return &PersistedStore{inner: inner, path: path}
}

func (p *PersistedStore) List(ctx context.Context) ([]core.Task, error) {
	return p.inner.List(ctx)
}

func (p *PersistedStore) Get(ctx context.Context, id string) (core.Task, error) {
	return p.inner.Get(ctx, id)
}

func (p *PersistedStore) Create(ctx context.Context, in core.NewTask) (core.Task, error) {
	t, err := p.inner.Create(ctx, in)
	if err != nil {
		return t, err
	}
	return t, p.saveSnapshot(ctx)
}

func (p *PersistedStore) Update(ctx context.Context, id string, in core.UpdateTask) (core.Task, error) {
	t, err := p.inner.Update(ctx, id, in)
	if err != nil {
		return t, err
	}
	return t, p.saveSnapshot(ctx)
}

func (p *PersistedStore) Delete(ctx context.Context, id string) error {
	if err := p.inner.Delete(ctx, id); err != nil {
		return err
	}
	return p.saveSnapshot(ctx)
}

func (p *PersistedStore) saveSnapshot(ctx context.Context) error {
	all, err := p.inner.List(ctx)
	if err != nil {
		return err
	}
	return SaveToJSON(p.path, all)
}

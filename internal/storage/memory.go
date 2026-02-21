package storage

import (
	"context"
	"encoding/json"
	"sync"
)

// memoryStore is a fast, ephemeral in-memory engine.
// Used EXCLUSIVELY in unit tests — zero disk I/O, zero setup.
// Tests run 100x faster than any file/DB engine.
type memoryStore struct {
	mu   sync.RWMutex
	data map[string]map[string][]byte // collection → key → serialized bytes
}

// NewMemoryStore constructs an in-memory engine.
// No path needed — lives and dies with the process.
func NewMemoryStore() Store {
	return &memoryStore{
		data: make(map[string]map[string][]byte),
	}
}

func (m *memoryStore) Save(_ context.Context, collection, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	encoded, err := json.Marshal(value)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	if _, ok := m.data[collection]; !ok {
		m.data[collection] = make(map[string][]byte)
	}
	m.data[collection][key] = encoded
	return nil
}

func (m *memoryStore) Get(_ context.Context, collection, key string, dest any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	col, ok := m.data[collection]
	if !ok {
		return ErrNotFound{Collection: collection, Key: key}
	}
	raw, ok := col[key]
	if !ok {
		return ErrNotFound{Collection: collection, Key: key}
	}
	return json.Unmarshal(raw, dest)
}

func (m *memoryStore) Delete(_ context.Context, collection, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	col, ok := m.data[collection]
	if !ok {
		return ErrNotFound{Collection: collection, Key: key}
	}
	if _, ok := col[key]; !ok {
		return ErrNotFound{Collection: collection, Key: key}
	}
	delete(col, key)
	return nil
}

func (m *memoryStore) List(_ context.Context, collection string, dest any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	col := m.data[collection]
	raws := make([]json.RawMessage, 0, len(col))
	for _, v := range col {
		raws = append(raws, v)
	}

	combined, err := json.Marshal(raws)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}
	return json.Unmarshal(combined, dest)
}

func (m *memoryStore) Close() error { return nil }

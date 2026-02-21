package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// jsonStore is the concrete JSON file-based implementation of Store.
// It satisfies the Store interface 100%.
// Task/Note modules never import this — they only know Store.
type jsonStore struct {
	baseDir string       // root folder where all JSON files live
	mu      sync.RWMutex // guards concurrent reads/writes
}

// NewJSONStore constructs a JSON engine rooted at baseDir.
// This is the only place JSON is mentioned outside this file.
func NewJSONStore(baseDir string) (Store, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, ErrStorageInit{Cause: err}
	}
	return &jsonStore{baseDir: baseDir}, nil
}

// collectionPath returns the file path for a collection.
// e.g. collection="tasks" → /data/tasks.json
func (s *jsonStore) collectionPath(collection string) string {
	return filepath.Join(s.baseDir, collection+".json")
}

// loadCollection reads the entire collection file into a raw map.
// Returns empty map if file doesn't exist yet (first run).
func (s *jsonStore) loadCollection(collection string) (map[string]json.RawMessage, error) {
	path := s.collectionPath(collection)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]json.RawMessage), nil // first time = empty store
	}
	if err != nil {
		return nil, ErrReadFailed{Collection: collection, Cause: err}
	}

	var records map[string]json.RawMessage
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, ErrCorruptData{Collection: collection, Cause: err}
	}
	return records, nil
}

// flushCollection writes the full map back to disk atomically.
func (s *jsonStore) flushCollection(collection string, records map[string]json.RawMessage) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	// Write to temp file first, then rename — atomic write pattern.
	// Prevents data corruption if process crashes mid-write.
	tmpPath := s.collectionPath(collection) + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}
	return os.Rename(tmpPath, s.collectionPath(collection))
}

// Save serializes value and stores it under collection[key].
func (s *jsonStore) Save(_ context.Context, collection, key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	records, err := s.loadCollection(collection)
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	records[key] = encoded
	return s.flushCollection(collection, records)
}

// Get deserializes a single record into dest.
func (s *jsonStore) Get(_ context.Context, collection, key string, dest any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records, err := s.loadCollection(collection)
	if err != nil {
		return err
	}

	raw, exists := records[key]
	if !exists {
		return ErrNotFound{Collection: collection, Key: key}
	}

	return json.Unmarshal(raw, dest)
}

// Delete removes a key from the collection.
func (s *jsonStore) Delete(_ context.Context, collection, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	records, err := s.loadCollection(collection)
	if err != nil {
		return err
	}

	if _, exists := records[key]; !exists {
		return ErrNotFound{Collection: collection, Key: key}
	}

	delete(records, key)
	return s.flushCollection(collection, records)
}

// List deserializes all records in a collection into dest (must be *[]T).
func (s *jsonStore) List(_ context.Context, collection string, dest any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records, err := s.loadCollection(collection)
	if err != nil {
		return err
	}

	// Collect all raw JSON values into a slice, then unmarshal into dest.
	raws := make([]json.RawMessage, 0, len(records))
	for _, v := range records {
		raws = append(raws, v)
	}

	combined, err := json.Marshal(raws)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	return json.Unmarshal(combined, dest)
}

// Close is a no-op for JSON (no connection to close).
// Required by interface — future engines will use this.
func (s *jsonStore) Close() error {
	return nil
}

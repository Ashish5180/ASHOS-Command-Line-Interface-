package storage

import "context"

// Store is the universal data engine contract.
// Any module (task, note, focus) talks ONLY to this interface.
// The actual engine (JSON, Bolt, SQLite) lives behind this wall.
type Store interface {
	// Save persists a value under a collection and unique key.
	// collection = "tasks", "notes", "focus_sessions"
	// key        = unique ID of the record
	// value      = any serializable Go struct
	Save(ctx context.Context, collection, key string, value any) error

	// Get retrieves a single record by key into dest (pointer to your struct).
	Get(ctx context.Context, collection, key string, dest any) error

	// Delete removes a record permanently.
	Delete(ctx context.Context, collection, key string) error

	// List fetches ALL records in a collection into dest (pointer to slice).
	List(ctx context.Context, collection string, dest any) error

	// Close gracefully shuts down the storage engine.
	// Critical for file flush, DB connection pool teardown, etc.
	Close() error
}

package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// sqliteStore is the concrete SQLite implementation of Store.
// It maps collections to rows in a single generic table.
type sqliteStore struct {
	db *sql.DB
}

// NewSQLiteStore constructs a SQLite engine at the given path.
// It ensures the schema is ready for use.
func NewSQLiteStore(dbPath string) (Store, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, ErrStorageInit{Cause: err}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, ErrStorageInit{Cause: err}
	}

	// Create generic table for collection-based storage
	schema := `
	CREATE TABLE IF NOT EXISTS records (
		collection TEXT,
		key TEXT,
		value BLOB,
		PRIMARY KEY (collection, key)
	);`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, ErrStorageInit{Cause: err}
	}

	return &sqliteStore{db: db}, nil
}

// Save serializes value and performs an UPSERT into the records table.
func (s *sqliteStore) Save(ctx context.Context, collection, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	query := `INSERT OR REPLACE INTO records (collection, key, value) VALUES (?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, collection, key, data)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	return nil
}

// Get retrieves a single record by key and deserializes it into dest.
func (s *sqliteStore) Get(ctx context.Context, collection, key string, dest any) error {
	query := `SELECT value FROM records WHERE collection = ? AND key = ?`
	var data []byte
	err := s.db.QueryRowContext(ctx, query, collection, key).Scan(&data)
	if err == sql.ErrNoRows {
		return ErrNotFound{Collection: collection, Key: key}
	}
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return ErrCorruptData{Collection: collection, Cause: err}
	}

	return nil
}

// Delete removes a record by key.
func (s *sqliteStore) Delete(ctx context.Context, collection, key string) error {
	query := `DELETE FROM records WHERE collection = ? AND key = ?`
	res, err := s.db.ExecContext(ctx, query, collection, key)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	if rows == 0 {
		return ErrNotFound{Collection: collection, Key: key}
	}

	return nil
}

// List fetches all values for a collection and unmarshals them into dest (slice pointer).
func (s *sqliteStore) List(ctx context.Context, collection string, dest any) error {
	query := `SELECT value FROM records WHERE collection = ?`
	rows, err := s.db.QueryContext(ctx, query, collection)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}
	defer rows.Close()

	var raws []json.RawMessage
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return ErrReadFailed{Collection: collection, Cause: err}
		}
		raws = append(raws, json.RawMessage(data))
	}

	if err := rows.Err(); err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	// Unmarshal the slice of raw JSONs into the dest pointer
	combined, err := json.Marshal(raws)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	return json.Unmarshal(combined, dest)
}

// GetDB returns the underlying *sql.DB connection.
// Used for modules that need direct SQL access (like AI vector search).
func (s *sqliteStore) GetDB() *sql.DB {
	return s.db
}

// Close closes the database connection.
func (s *sqliteStore) Close() error {
	return s.db.Close()
}

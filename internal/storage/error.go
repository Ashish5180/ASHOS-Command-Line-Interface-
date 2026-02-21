package storage

import "fmt"

// ErrNotFound is returned when a key does not exist in a collection.
// Callers can type-check this to handle "not found" vs "system error" differently.
type ErrNotFound struct {
	Collection string
	Key        string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("storage: record not found [collection=%s key=%s]", e.Collection, e.Key)
}

// ErrStorageInit is returned when the storage engine fails to initialize.
type ErrStorageInit struct {
	Cause error
}

func (e ErrStorageInit) Error() string {
	return fmt.Sprintf("storage: initialization failed: %v", e.Cause)
}

// ErrReadFailed is returned when a collection cannot be read.
type ErrReadFailed struct {
	Collection string
	Cause      error
}

func (e ErrReadFailed) Error() string {
	return fmt.Sprintf("storage: read failed [collection=%s]: %v", e.Collection, e.Cause)
}

// ErrWriteFailed is returned when a record cannot be persisted.
type ErrWriteFailed struct {
	Collection string
	Cause      error
}

func (e ErrWriteFailed) Error() string {
	return fmt.Sprintf("storage: write failed [collection=%s]: %v", e.Collection, e.Cause)
}

// ErrCorruptData is returned when stored data cannot be deserialized.
type ErrCorruptData struct {
	Collection string
	Cause      error
}

func (e ErrCorruptData) Error() string {
	return fmt.Sprintf("storage: corrupt data [collection=%s]: %v", e.Collection, e.Cause)
}

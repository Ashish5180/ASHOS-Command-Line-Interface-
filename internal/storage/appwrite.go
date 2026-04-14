package storage

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
)

// appwriteStore implements the Store interface using Appwrite Databases.
// Note: We do NOT store *appwrite.Client here because it is an internal type
// in the SDK — only databases.Databases is needed after initialization.
type appwriteStore struct {
	db   *databases.Databases
	dbID string
}

// NewAppwriteStore creates a new storage engine backed by Appwrite.
func NewAppwriteStore(endpoint, projectID, apiKey, dbID string) Store {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(endpoint),
		appwrite.WithProject(projectID),
		appwrite.WithKey(apiKey),
	)

	return &appwriteStore{
		db:   databases.New(client),
		dbID: dbID,
	}
}

// Aur bhi clean approach — directly upsert karo
func (s *appwriteStore) Save(_ context.Context, collection, key string, value any) error {
	data, err := toMap(value)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	// Create ya Update automatically handle karta hai!
	_, err = s.db.UpsertDocument(
		s.dbID,
		collection,
		key,
		s.db.WithUpsertDocumentData(data),
	)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}

	return nil
}

// Get fetches a single document and deserializes it into dest.
func (s *appwriteStore) Get(_ context.Context, collection, key string, dest any) error {
	doc, err := s.db.GetDocument(s.dbID, collection, key)
	if err != nil {
		return ErrNotFound{Collection: collection, Key: key}
	}

	// Marshaling full doc to handle internal private data
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return ErrCorruptData{Collection: collection, Cause: err}
	}

	// Logic to unwrap array if it was stored inside a map wrapper (as a Stringified JSON)
	var temp map[string]any
	if err := json.Unmarshal(jsonData, &temp); err == nil {
		if val, ok := temp["data"]; ok {
			var inner []byte
			if strVal, ok := val.(string); ok {
				// Agar string hai (Appwrite String Attribute), toh directly use karo
				inner = []byte(strVal)
			} else {
				// Fallback: Agar kisi wajah se map format mein hi mil gaya
				inner, _ = json.Marshal(val)
			}
			return json.Unmarshal(inner, dest)
		}
	}

	return json.Unmarshal(jsonData, dest)
}

// ... (Rest of the methods remain same, only toMap and Get/List logic changes)

func (s *appwriteStore) Delete(_ context.Context, collection, key string) error {
	_, err := s.db.DeleteDocument(s.dbID, collection, key)
	if err != nil {
		return ErrWriteFailed{Collection: collection, Cause: err}
	}
	return nil
}

func (s *appwriteStore) List(_ context.Context, collection string, dest any) error {
	result, err := s.db.ListDocuments(s.dbID, collection)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	var items []json.RawMessage
	for _, doc := range result.Documents {
		data, err := json.Marshal(doc)
		if err == nil {
			// Check if we need to unwrap stringified JSON
			var temp map[string]any
			if json.Unmarshal(data, &temp) == nil {
				if val, ok := temp["data"]; ok {
					var inner []byte
					if strVal, ok := val.(string); ok {
						inner = []byte(strVal)
					} else {
						inner, _ = json.Marshal(val)
					}
					items = append(items, inner)
					continue
				}
			}
			items = append(items, data)
		}
	}

	combined, err := json.Marshal(items)
	if err != nil {
		return ErrReadFailed{Collection: collection, Cause: err}
	}

	return json.Unmarshal(combined, dest)
}

func (s *appwriteStore) Close() error {
	return nil
}

func toMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Agar data array hai, toh isse JSON String mein convert karke wrap karo
	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "[") {
		return map[string]any{"data": string(data)}, nil
	}

	var res map[string]any
	err = json.Unmarshal(data, &res)
	return res, err
}

package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// SyncStore handles dual-write logic (Local + Cloud).
// It implements the Store interface, making it transparent to services.
type SyncStore struct {
	local Store
	cloud Store
}

// NewSyncStore coordinates two different storage engines.
// local: usually SQLite for offline-first performance.
// cloud: usually Appwrite for cross-device synchronization.
func NewSyncStore(local, cloud Store) Store {
	return &SyncStore{
		local: local,
		cloud: cloud,
	}
}

func (s *SyncStore) Save(ctx context.Context, collection, key string, value any) error {
	// 1. Persist locally first (Source of Truth for this machine)
	if err := s.local.Save(ctx, collection, key, value); err != nil {
		return fmt.Errorf("local save failed: %w", err)
	}

	// 2. Cloud Sync
	if s.cloud != nil {
		if cloudErr := s.cloud.Save(ctx, collection, key, value); cloudErr != nil {
			fmt.Printf("⚠️  Cloud Sync Failed: %v\n", cloudErr)
		}
	}

	return nil
}

func (s *SyncStore) Get(ctx context.Context, collection, key string, dest any) error {
	// Always read local for speed and offline availability
	return s.local.Get(ctx, collection, key, dest)
}

func (s *SyncStore) List(ctx context.Context, collection string, dest any) error {
	return s.local.List(ctx, collection, dest)
}

func (s *SyncStore) Delete(ctx context.Context, collection, key string) error {
	if err := s.local.Delete(ctx, collection, key); err != nil {
		return err
	}

	if s.cloud != nil {
		go func() {
			_ = s.cloud.Delete(ctx, collection, key)
		}()
	}
	return nil
}

func (s *SyncStore) Close() error {
	err := s.local.Close()
	if s.cloud != nil {
		if cloudErr := s.cloud.Close(); cloudErr != nil {
			err = cloudErr
		}
	}
	return err
}

// GetDB satisfies the interface for AI module direct SQL access.
func (s *SyncStore) GetDB() *sql.DB {
	if getter, ok := s.local.(interface{ GetDB() *sql.DB }); ok {
		return getter.GetDB()
	}
	return nil
}

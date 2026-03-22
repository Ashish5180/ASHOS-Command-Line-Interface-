package ai

import (
	"context"
)

// Repository defines the contract for storing and searching AI context.
type Repository interface {
	SaveEmbedding(ctx context.Context, record EmbeddingRecord, embedding []float32) error
	GetRecordByHash(ctx context.Context, hash string) (*EmbeddingRecord, error)
	SearchEmbeddings(ctx context.Context, embedding []float32, limit int) ([]EmbeddingRecord, error)
	GetRecentActions(ctx context.Context, limit int) ([]EmbeddingRecord, error)
}

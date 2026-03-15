package ai

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)


// NewSQLiteRepository creates a new AI repository using the provided database.
func NewSQLiteRepository(db *sql.DB) (Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Try to initialize sqlite-vec table
	schema := `
	CREATE VIRTUAL TABLE IF NOT EXISTS ai_embeddings USING vec0(
		collection TEXT,
		key TEXT,
		content TEXT,
		metadata TEXT,
		created_at TEXT,
		embedding FLOAT[1536]
	);`

	if _, err := db.Exec(schema); err != nil {
		// Fallback to standard table if vec0 is not available
		// Logged silently or as debug for clean CLI
		fallbackSchema := `
		CREATE TABLE IF NOT EXISTS ai_embeddings_fallback (
			collection TEXT,
			key TEXT,
			content TEXT,
			metadata TEXT,
			created_at TEXT,
			embedding BLOB,
			PRIMARY KEY (collection, key)
		);`
		if _, err := db.Exec(fallbackSchema); err != nil {
			return nil, fmt.Errorf("failed to initialize fallback table: %w", err)
		}
		return &sqliteRepo{db: db, isFallback: true}, nil
	}

	return &sqliteRepo{db: db, isFallback: false}, nil
}

type sqliteRepo struct {
	db         *sql.DB
	isFallback bool
}

func (r *sqliteRepo) serializeEmbedding(embedding []float32) []byte {
	buf := make([]byte, len(embedding)*4)
	for i, f := range embedding {
		bits := math.Float32bits(f)
		copy(buf[i*4:], []byte{
			byte(bits),
			byte(bits >> 8),
			byte(bits >> 16),
			byte(bits >> 24),
		})
	}
	return buf
}

func (r *sqliteRepo) deserializeEmbedding(buf []byte) []float32 {
	embedding := make([]float32, len(buf)/4)
	for i := 0; i < len(embedding); i++ {
		bits := uint32(buf[i*4]) |
			uint32(buf[i*4+1])<<8 |
			uint32(buf[i*4+2])<<16 |
			uint32(buf[i*4+3])<<24
		embedding[i] = math.Float32frombits(bits)
	}
	return embedding
}

func (r *sqliteRepo) SaveEmbedding(ctx context.Context, record EmbeddingRecord, embedding []float32) error {
	embData := r.serializeEmbedding(embedding)

	tableName := "ai_embeddings"
	if r.isFallback {
		tableName = "ai_embeddings_fallback"
	}

	query := fmt.Sprintf(`INSERT OR REPLACE INTO %s(collection, key, content, metadata, created_at, embedding) 
	          VALUES (?, ?, ?, ?, ?, ?)`, tableName)
	
	_, err := r.db.ExecContext(ctx, query, 
		record.Collection, 
		record.Key, 
		record.Content, 
		record.Metadata, 
		record.CreatedAt.Format(time.RFC3339),
		embData,
	)
	return err
}

func (r *sqliteRepo) SearchEmbeddings(ctx context.Context, embedding []float32, limit int) ([]EmbeddingRecord, error) {
	if r.isFallback {
		return r.searchFallback(ctx, embedding, limit)
	}

	embData := r.serializeEmbedding(embedding)
	query := `
		SELECT collection, key, content, metadata, created_at
		FROM ai_embeddings
		WHERE embedding MATCH ?
		ORDER BY distance
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, embData, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *sqliteRepo) searchFallback(ctx context.Context, queryEmb []float32, limit int) ([]EmbeddingRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT collection, key, content, metadata, created_at, embedding FROM ai_embeddings_fallback`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type scoredRecord struct {
		record EmbeddingRecord
		score  float64
	}
	var scored []scoredRecord

	for rows.Next() {
		var rec EmbeddingRecord
		var createdAtStr string
		var embData []byte
		if err := rows.Scan(&rec.Collection, &rec.Key, &rec.Content, &rec.Metadata, &createdAtStr, &embData); err != nil {
			return nil, err
		}
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		
		itemEmb := r.deserializeEmbedding(embData)
		score := r.cosineSimilarity(queryEmb, itemEmb)
		scored = append(scored, scoredRecord{rec, score})
	}

	// Sort by score descending (higher is better for cosine similarity)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	var result []EmbeddingRecord
	for i := 0; i < len(scored) && i < limit; i++ {
		result = append(result, scored[i].record)
	}
	return result, nil
}

func (r *sqliteRepo) cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (r *sqliteRepo) scanRows(rows *sql.Rows) ([]EmbeddingRecord, error) {
	var records []EmbeddingRecord
	for rows.Next() {
		var rec EmbeddingRecord
		var createdAtStr string
		if err := rows.Scan(&rec.Collection, &rec.Key, &rec.Content, &rec.Metadata, &createdAtStr); err != nil {
			return nil, err
		}
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		records = append(records, rec)
	}
	return records, nil
}

func (r *sqliteRepo) GetRecentActions(ctx context.Context, limit int) ([]EmbeddingRecord, error) {
	tableName := "ai_embeddings"
	if r.isFallback {
		tableName = "ai_embeddings_fallback"
	}
	query := fmt.Sprintf(`
		SELECT collection, key, content, metadata, created_at
		FROM %s
		ORDER BY created_at DESC
		LIMIT ?`, tableName)

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

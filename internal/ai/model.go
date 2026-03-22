package ai

import "time"

// EmbeddingRecord represents a vector embedding stored in the database.
type EmbeddingRecord struct {
	Collection  string    `json:"collection"`
	SourceID    string    `json:"source_id"`
	SourceType  string    `json:"source_type"`
	ContentHash string    `json:"content_hash"`
	Content     string    `json:"content"`
	Metadata    string    `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
	Score       float64   `json:"score,omitempty"`
}

// ConversationMessage represents a single message in an LLM conversation.
type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

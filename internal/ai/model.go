package ai

import "time"

// EmbeddingRecord represents a vector embedding stored in the database.
type EmbeddingRecord struct {
	Collection string    `json:"collection"`
	Key        string    `json:"key"`
	Content    string    `json:"content"`
	Metadata   string    `json:"metadata"`
	CreatedAt  time.Time `json:"created_at"`
}

// ConversationMessage represents a single message in an LLM conversation.
type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

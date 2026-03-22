package ai

import "github.com/sashabaranov/go-openai"

type ConversationMemory struct {
	SessionID string
	History   []openai.ChatCompletionMessage
	MaxTurns  int
}

func NewMemory(maxTurns int) *ConversationMemory {
	return &ConversationMemory{
		History:  make([]openai.ChatCompletionMessage, 0),
		MaxTurns: maxTurns,
	}
}

func (m *ConversationMemory) AddMessage(role string, content string) {
	m.History = append(m.History, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})

	// Keep only the last MaxTurns * 2 (user + assistant) messages
	limit := m.MaxTurns * 2
	if len(m.History) > limit {
		m.History = m.History[len(m.History)-limit:]
	}
}

func (m *ConversationMemory) GetHistory() []openai.ChatCompletionMessage {
	return m.History
}

package ai

import (
	"context"
	"fmt"
	"time"

	"ashos/internal/core/event"
	"ashos/internal/system"
	"github.com/sashabaranov/go-openai"
)

type service struct {
	client *openai.Client
	repo   Repository
	bus    *event.EventBus
	sys    *system.Service
}

// NewService creates a new AI service and wires it to the event bus.
func NewService(repo Repository, bus *event.EventBus, sys *system.Service) Service {
	// API key hardcoded for now as requested
	apiKey := "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"
	client := openai.NewClient(apiKey)

	s := &service{
		client: client,
		repo:   repo,
		bus:    bus,
		sys:    sys,
	}

	// Subscribe to events for auto-ingestion
	bus.Subscribe(event.TaskCreated{}, func(e any) {
		evt := e.(event.TaskCreated)
		s.IngestTask(context.Background(), evt.ID, evt.Title)
	})

	bus.Subscribe(event.FocusEnded{}, func(e any) {
		evt := e.(event.FocusEnded)
		s.IngestFocusSession(context.Background(), evt.StartTime, evt.Duration, evt.Summary)
	})

	bus.Subscribe(event.NoteCreated{}, func(e any) {
		evt := e.(event.NoteCreated)
		s.IngestNote(context.Background(), evt.ID, evt.Content)
	})

	bus.Subscribe(event.SprintEnded{}, func(e any) {
		evt := e.(event.SprintEnded)
		s.IngestSprint(context.Background(), evt.ID, evt.Title, evt.Summary)
	})

	return s
}

func (s *service) Ask(ctx context.Context, query string) (string, error) {
	if s.repo == nil {
		return "", fmt.Errorf("AI brain is currently disabled (storage error)")
	}
	// 1. Generate embedding for query
	queryEmb, err := s.generateEmbedding(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// 2. Search relevant context
	records, err := s.repo.SearchEmbeddings(ctx, queryEmb, 5)
	if err != nil {
		return "", fmt.Errorf("failed to search context: %w", err)
	}

	// 3. Build context string
	contextStr := ""
	for _, r := range records {
		contextStr += fmt.Sprintf("[%s] %s (at %s)\n", r.Collection, r.Content, r.CreatedAt.Format(time.RFC822))
	}

	// 4. Get System Status for quantitative context
	status := s.sys.GetStatus()
	sysContext := fmt.Sprintf("Current Time: %s\nUptime: %v\nPending Tasks: %d\nFocus Today: %v", 
		time.Now().Format(time.RFC1123), status.Uptime, status.PendingTasks, status.FocusToday)

	// 5. Call LLM
	prompt := fmt.Sprintf(`You are ASHOS AI Brain, a personalized intelligence layer for the user's CLI OS.
Answer based on the provided System Context and Archive History.

IMPORTANT: Do NOT use markdown stars (**) for bolding. Use plain text or CAPITALIZATION for emphasis.
Keep the formatting clean for a terminal.

System Context:
%s

Archive History (Relevant snippets):
%s

Question: %s`, sysContext, contextStr, query)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *service) GenerateStandup(ctx context.Context) (string, error) {
	if s.repo == nil {
		return "", fmt.Errorf("AI brain is currently disabled (storage error)")
	}
	// Get recent actions for the last 24 hours
	records, err := s.repo.GetRecentActions(ctx, 20)
	if err != nil {
		return "", err
	}

	contextStr := ""
	for _, r := range records {
		contextStr += fmt.Sprintf("[%s] %s\n", r.Collection, r.Content)
	}

	status := s.sys.GetStatus()
	sysContext := fmt.Sprintf("Current Time: %s\nTasks Pending: %d\nFocus Today: %v", 
		time.Now().Format(time.RFC822), status.PendingTasks, status.FocusToday)

	prompt := fmt.Sprintf(`Generate a professional and concise daily standup report.
Do NOT use markdown stars (**) or markdown bullet points. Use simple dashes (-) or numbers.
Keep it premium but plain-text friendly. 

Today's Stats:
%s

Recent Activity Feed:
%s`, sysContext, contextStr)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *service) SuggestNextTask(ctx context.Context) (string, error) {
	if s.repo == nil {
		return "", fmt.Errorf("AI brain is currently disabled (storage error)")
	}
	records, err := s.repo.GetRecentActions(ctx, 10)
	if err != nil {
		return "", err
	}

	contextStr := ""
	for _, r := range records {
		contextStr += fmt.Sprintf("[%s] %s\n", r.Collection, r.Content)
	}

	status := s.sys.GetStatus()
	prompt := fmt.Sprintf(`Based on the user's recent patterns and current status, suggest what to work on next.
Do NOT use markdown stars (**) for bolding. 
Current Stats: %d pending tasks, %v focused today.

Recent Activity:
%s`, status.PendingTasks, status.FocusToday, contextStr)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *service) IngestTask(ctx context.Context, taskID int, title string) error {
	if s.repo == nil {
		return nil // Silent fail for events if repo is disabled
	}
	fmt.Printf("🧠 AI Brain: Ingesting task '%s'...\n", title)
	content := fmt.Sprintf("Created task: %s", title)
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection: "tasks",
		Key:        fmt.Sprintf("task_%d", taskID),
		Content:    content,
		CreatedAt:  time.Now(),
	}, embedding)
}

func (s *service) IngestFocusSession(ctx context.Context, startTime time.Time, duration time.Duration, summary string) error {
	if s.repo == nil {
		return nil
	}
	fmt.Printf("🧠 AI Brain: Ingesting focus session (%v)...\n", duration.Round(time.Second))
	content := fmt.Sprintf("Focused for %s. Summary: %s", duration.Round(time.Minute), summary)
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection: "focus",
		Key:        fmt.Sprintf("focus_%d", startTime.Unix()),
		Content:    content,
		CreatedAt:  time.Now(),
	}, embedding)
}

func (s *service) IngestNote(ctx context.Context, noteID int, content string) error {
	if s.repo == nil {
		return nil
	}
	fmt.Printf("🧠 AI Brain: Ingesting journal entry...\n")
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection: "notes",
		Key:        fmt.Sprintf("note_%d", noteID),
		Content:    content,
		CreatedAt:  time.Now(),
	}, embedding)
}

func (s *service) IngestSprint(ctx context.Context, sprintID int, title, summary string) error {
	if s.repo == nil {
		return nil
	}
	fmt.Printf("🧠 AI Brain: Ingesting sprint review '%s'...\n", title)
	content := fmt.Sprintf("Sprint: %s. Review: %s", title, summary)
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection: "sprints",
		Key:        fmt.Sprintf("sprint_%d", sprintID),
		Content:    content,
		CreatedAt:  time.Now(),
	}, embedding)
}

func (s *service) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

type Service interface {
	Ask(ctx context.Context, query string) (string, error)
	GenerateStandup(ctx context.Context) (string, error)
	SuggestNextTask(ctx context.Context) (string, error)
	IngestTask(ctx context.Context, taskID int, title string) error
	IngestFocusSession(ctx context.Context, startTime time.Time, duration time.Duration, summary string) error
	IngestNote(ctx context.Context, noteID int, content string) error
	IngestSprint(ctx context.Context, sprintID int, title, summary string) error
}

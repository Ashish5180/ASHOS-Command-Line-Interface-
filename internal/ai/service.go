package ai

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"ashos/internal/core/event"
	"ashos/internal/system"
	"github.com/sashabaranov/go-openai"
	"net/http"
)

type service struct {
	client   *openai.Client
	repo     Repository
	bus      *event.EventBus
	sys      *system.Service
	memory   *ConversationMemory
	useLocal bool
}

// NewService creates a new AI service and wires it to the event bus.
func NewService(repo Repository, bus *event.EventBus, sys *system.Service) Service {
	// API key hardcoded for now as requested
	apiKey := "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"
	client := openai.NewClient(apiKey)

	s := &service{
		client:   client,
		repo:     repo,
		bus:      bus,
		sys:      sys,
		memory:   NewMemory(10), // Remember last 10 turns
		useLocal: false,         // Default to Cloud
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
func (s *service) SetLocalMode(local bool) {
	s.useLocal = local
	if local {
		config := openai.DefaultConfig("ollama")
		config.BaseURL = "http://localhost:11434/v1"
		s.client = openai.NewClientWithConfig(config)
		
		// Attempt a quick ping to confirm Ollama is running
		resp, err := http.Get("http://localhost:11434/")
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("⚠️ ASHOS AI: Local mode failed. Ollama not detected at localhost:11434.")
			fmt.Println("💡 Tip: Run 'ollama serve' to enable local intelligence.")
			s.useLocal = false 
			// Re-initialize cloud client
			apiKey := "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"
			s.client = openai.NewClient(apiKey)
			return
		}
		
		fmt.Println("🚀 AI Brain: Successfully switched to Local Mode (Ollama)")
	}
}

func (s *service) hasInternet() bool {
	// Simple check: can we reach google?
	// In a real app, you'd use a more robust check or just try the API call.
	return true 
}
func (s *service) Ask(ctx context.Context, query string) (string, []EmbeddingRecord, error) {
	if s.repo == nil {
		return "", nil, fmt.Errorf("AI brain is currently disabled (storage error)")
	}
	// 1. Generate embedding for query
	queryEmb, err := s.generateEmbedding(ctx, query)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// 2. Search relevant context
	records, err := s.repo.SearchEmbeddings(ctx, queryEmb, 5)
	if err != nil {
		return "", nil, fmt.Errorf("failed to search context: %w", err)
	}

	// 3. Build context string
	contextStr := ""
	for _, r := range records {
		contextStr += fmt.Sprintf("[%s] %s (at %s)\n", r.Collection, r.Content, r.CreatedAt.Format(time.RFC822))
	}

	// 4. Get System Status for quantitative context
	status := s.sys.GetStatus()
	sysContext := fmt.Sprintf("Current Time: %s\nPending Tasks: %d\nFocus Today: %v",
		time.Now().Format(time.RFC1123), status.PendingTasks, status.FocusToday)

	// 5. Build Chat History
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are ASHOS AI - a personal productivity intelligence.
You have access to: tasks, focus sessions, notes, streaks.
Don't just list data. Analyze it. Give insights. Be direct.
If productivity is low, say it. Suggest next action always.
Use plain text. NEVER use markdown bolding (**) or headers. 
Maintain a sleek, premium, and direct tone.`,
		},
	}

	// Add Archive Knowledge as context
	knowledgePrompt := fmt.Sprintf("SYSTEM CONTEXT:\n%s\n\nARCHIVE KNOWLEDGE:\n%s", sysContext, contextStr)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: knowledgePrompt,
	})

	// Add Conversation Memory
	messages = append(messages, s.memory.GetHistory()...)

	// Add Current Query
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: query,
	})

	model := openai.GPT4oMini
	if s.useLocal {
		model = "llama3.2" // Ollama default
	}

	// 6. Call LLM
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", nil, err
	}

	answer := resp.Choices[0].Message.Content

	// 7. Save to memory
	s.memory.AddMessage(openai.ChatMessageRoleUser, query)
	s.memory.AddMessage(openai.ChatMessageRoleAssistant, answer)

	return answer, records, nil
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
FORMAT:
Yesterday : Summary of 24h activity
Today     : Based on pending tasks and notes
Blockers  : Any potential issues
Focus     : %v logged today

Do NOT use markdown stars (**) or markdown headers. 
Use the following format style:

🗒️ Daily Standup — [Date]
━━━━━━━━━━━━━━━━━━━━━━━━━
Yesterday : [Content]
Today     : [Content]
Blockers  : [Content]
Focus     : %v logged
━━━━━━━━━━━━━━━━━━━━━━━━━

Current Stats:
%s

Recent Activity Feed:
%s`, status.FocusToday, status.FocusToday, sysContext, contextStr)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
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
	records, err := s.repo.GetRecentActions(ctx, 15)
	if err != nil {
		return "", err
	}

	contextStr := ""
	for _, r := range records {
		contextStr += fmt.Sprintf("[%s] %s\n", r.Collection, r.Content)
	}

	status := s.sys.GetStatus()
	prompt := fmt.Sprintf(`You are ASHOS AI. Based on the user's recent patterns, pending tasks, and current focus state, suggest the absolute best thing to work on right now.
Be direct. Give a reasoning. 

Current Stats: %d pending tasks, %v focused today.
Recent Activity:
%s

Format:
Suggestion: [Task Name]
Reasoning: [Why this matters now]
Tip: [A small productivity tip for this task]`, status.PendingTasks, status.FocusToday, contextStr)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
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

	content := fmt.Sprintf("Created task: %s", title)
	hash := calculateHash(content)

	// Check if already archived
	existing, _ := s.repo.GetRecordByHash(ctx, hash)
	if existing != nil {
		return nil
	}

	fmt.Printf("🧠 AI Brain: Ingesting task '%s'...\n", title)
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection:  "tasks",
		SourceID:    fmt.Sprintf("%d", taskID),
		SourceType:  "task",
		ContentHash: hash,
		Content:     content,
		CreatedAt:   time.Now(),
	}, embedding)
}

func (s *service) IngestFocusSession(ctx context.Context, startTime time.Time, duration time.Duration, summary string) error {
	if s.repo == nil {
		return nil
	}

	content := fmt.Sprintf("Focused for %s. Summary: %s", duration.Round(time.Minute), summary)
	hash := calculateHash(content)

	existing, _ := s.repo.GetRecordByHash(ctx, hash)
	if existing != nil {
		return nil
	}

	fmt.Printf("🧠 AI Brain: Ingesting focus session (%v)...\n", duration.Round(time.Second))
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection:  "focus",
		SourceID:    fmt.Sprintf("%d", startTime.Unix()),
		SourceType:  "focus",
		ContentHash: hash,
		Content:     content,
		CreatedAt:   time.Now(),
	}, embedding)
}

func (s *service) IngestNote(ctx context.Context, noteID int, content string) error {
	if s.repo == nil {
		return nil
	}

	hash := calculateHash(content)
	existing, _ := s.repo.GetRecordByHash(ctx, hash)
	if existing != nil {
		return nil
	}

	fmt.Printf("🧠 AI Brain: Ingesting journal entry...\n")
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection:  "notes",
		SourceID:    fmt.Sprintf("%d", noteID),
		SourceType:  "note",
		ContentHash: hash,
		Content:     content,
		CreatedAt:   time.Now(),
	}, embedding)
}

func (s *service) IngestSprint(ctx context.Context, sprintID int, title, summary string) error {
	if s.repo == nil {
		return nil
	}

	content := fmt.Sprintf("Sprint: %s. Review: %s", title, summary)
	hash := calculateHash(content)
	existing, _ := s.repo.GetRecordByHash(ctx, hash)
	if existing != nil {
		return nil
	}

	fmt.Printf("🧠 AI Brain: Ingesting sprint review '%s'...\n", title)
	embedding, err := s.generateEmbedding(ctx, content)
	if err != nil {
		return err
	}

	return s.repo.SaveEmbedding(ctx, EmbeddingRecord{
		Collection:  "sprints",
		SourceID:    fmt.Sprintf("%d", sprintID),
		SourceType:  "sprint",
		ContentHash: hash,
		Content:     content,
		CreatedAt:   time.Now(),
	}, embedding)
}

func calculateHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (s *service) GenerateDailyDigest(ctx context.Context) (string, error) {
	if s.repo == nil {
		return "", fmt.Errorf("AI brain is currently disabled (storage error)")
	}

	status := s.sys.GetStatus()
	records, err := s.repo.GetRecentActions(ctx, 20)
	if err != nil {
		return "", err
	}

	activityStr := ""
	for _, r := range records {
		activityStr += fmt.Sprintf("[%s] %s\n", r.Collection, r.Content)
	}

	prompt := fmt.Sprintf(`Generate "ASHOS Daily Intelligence" report.
Data:
Focus Time: %v
Tasks Completed: %d (estimated from baseline)
Activity: %s

FORMAT:
📊 ASHOS Daily Intelligence — [Date]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Productivity Score : [0-100]/100
Focus Time         : %v
Tasks Completed    : [X]
AI Queries Made    : [Estimate]

💡 Insight: [One high-lvl behavioral insight about their work patterns]

🎯 Tomorrow's Suggestion:
   1. [Suggestion 1]
   2. [Suggestion 2]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Make it professional, data-driven, and insightful. Avoid markdown headers.`, status.FocusToday, status.PendingTasks, activityStr, status.FocusToday)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *service) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	model := openai.AdaEmbeddingV2
	if s.useLocal {
		model = "nomic-embed-text" // Standard Ollama embedding
	}

	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: model,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func (s *service) Reset(ctx context.Context) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.Reset(ctx)
}

type Service interface {
	Ask(ctx context.Context, query string) (string, []EmbeddingRecord, error)
	GenerateStandup(ctx context.Context) (string, error)
	SuggestNextTask(ctx context.Context) (string, error)
	GenerateDailyDigest(ctx context.Context) (string, error)
	SetLocalMode(local bool)
	IngestTask(ctx context.Context, taskID int, title string) error
	IngestFocusSession(ctx context.Context, startTime time.Time, duration time.Duration, summary string) error
	IngestNote(ctx context.Context, noteID int, content string) error
	IngestSprint(ctx context.Context, sprintID int, title, summary string) error
	Reset(ctx context.Context) error
}

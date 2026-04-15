package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
)

// appwriteStore implements the Store interface using Appwrite Databases.
// Note: We do NOT store *appwrite.Client here because it is an internal type
// in the SDK — only databases.Databases is needed after initialization.
type appwriteStore struct {
	db        *databases.Databases
	dbID      string
	endpoint  string
	projectID string
	apiKey    string
}

// 1. Interfaces (To avoid tight coupling)
// Ye kisi bhi Storage (Appwrite, SQL, etc.) ke sath kaam karega
type DataRepository interface {
	GetWeeklyNotes(ctx context.Context, startTime time.Time) ([]string, error)
	GetWeeklyTasks(ctx context.Context, startTime time.Time) ([]string, error)
	UpdatePersonalContext(ctx context.Context, summary string) error
	GetAboutMe(ctx context.Context) (string, error)
}

// Ye kisi bhi AI (OpenAI, Anthropic, etc.) ke sath kaam karega
type Summarizer interface {
	SummarizeWeek(ctx context.Context, currentContext string, data string) (string, error)
}

// 2. Main Service (Independent Logic)
type WeeklySyncService struct {
	repo      DataRepository
	aiService Summarizer
}

// NewWeeklySyncService creates a new instance of the weekly synchronization service.
func NewWeeklySyncService(r DataRepository, ai Summarizer) *WeeklySyncService {
	return &WeeklySyncService{
		repo:      r,
		aiService: ai,
	}
}

// NewAppwriteStore creates a new storage engine backed by Appwrite.
func NewAppwriteStore(endpoint, projectID, apiKey, dbID string) Store {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(endpoint),
		appwrite.WithProject(projectID),
		appwrite.WithKey(apiKey),
	)

	return &appwriteStore{
		db:        databases.New(client),
		dbID:      dbID,
		endpoint:  endpoint,
		projectID: projectID,
		apiKey:    apiKey,
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

	// MASTER HACK: Marshal the whole doc to reveal hidden fields
	bytes, _ := json.Marshal(doc)
	var temp map[string]any
	json.Unmarshal(bytes, &temp)

	if val, ok := temp["data"]; ok {
		var inner []byte
		if strVal, ok := val.(string); ok {
			inner = []byte(strVal)
		} else {
			inner, _ = json.Marshal(val)
		}
		return json.Unmarshal(inner, dest)
	}

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return ErrCorruptData{Collection: collection, Cause: err}
	}
	return json.Unmarshal(jsonData, dest)
}

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

// Core Logic: Hafte bhar ka data summarize karke context update karna
func (s *WeeklySyncService) RunWeeklySync() error {
	ctx := context.Background()
	lastWeek := time.Now().AddDate(0, 0, -7)
	fmt.Printf("📂 Starting Weekly Sync (From: %v)\n", lastWeek.Format("2006-01-02"))

	// A. Data Fetch karo
	fmt.Println("🔍 Fetching weekly notes...")
	notes, err := s.repo.GetWeeklyNotes(ctx, lastWeek)
	if err != nil {
		return fmt.Errorf("failed to fetch notes: %v", err)
	}
	fmt.Printf("✅ Found %d notes.\n", len(notes))

	fmt.Println("🔍 Fetching weekly tasks...")
	tasks, err := s.repo.GetWeeklyTasks(ctx, lastWeek)
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %v", err)
	}
	fmt.Printf("✅ Found %d tasks.\n", len(tasks))

	// B. Purana context (About Me) uthao
	fmt.Println("🧠 Reading current 'About Me' context...")
	oldContext, _ := s.repo.GetAboutMe(ctx)

	// C. Summarize karo logic
	fmt.Println("🤖 Requesting AI to analyze your week and update personality...")
	combinedData := fmt.Sprintf("Notes: %v \n Tasks: %v", notes, tasks)

	newSummary, err := s.aiService.SummarizeWeek(ctx, oldContext, combinedData)
	if err != nil {
		return fmt.Errorf("ai summarization failed: %v", err)
	}
	fmt.Println("✨ AI Summary Generated Successfully.")

	// D. "Main Brain" update karo
	fmt.Println("💾 Updating your Personal OS Brain in Appwrite...")
	err = s.repo.UpdatePersonalContext(ctx, newSummary)
	if err != nil {
		return fmt.Errorf("failed to update brain: %v", err)
	}

	fmt.Println("🎉 Weekly Sync Complete! Your AI OS is now more personalized.")
	return nil
}

func (s *appwriteStore) GetWeeklyNotes(ctx context.Context, startTime time.Time) ([]string, error) {
	// 1. Direct API Request (Bypassing SDK limitations)
	url := fmt.Sprintf("%s/databases/%s/collections/notes/documents", s.endpoint, s.dbID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Appwrite-Project", s.projectID)
	req.Header.Set("X-Appwrite-Key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return nil, nil }
	defer resp.Body.Close()

	var result struct {
		Documents []map[string]any `json:"documents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil
	}

	var notes []string
	for _, doc := range result.Documents {
		// Clean data extraction from Raw JSON Map
		if dataAttr, ok := doc["data"]; ok {
			var rawEntries []any
			if str, ok := dataAttr.(string); ok {
				json.Unmarshal([]byte(str), &rawEntries)
			} else if slice, ok := dataAttr.([]any); ok {
				rawEntries = slice
			}
			for _, item := range rawEntries {
				if entry, ok := item.(map[string]any); ok {
					if content, ok := entry["content"].(string); ok { notes = append(notes, content) }
				}
			}
		}
	}
	fmt.Printf("✅ [Direct API] Found %d notes.\n", len(notes))
	return notes, nil
}

func (s *appwriteStore) GetWeeklyTasks(ctx context.Context, startTime time.Time) ([]string, error) {
	url := fmt.Sprintf("%s/databases/%s/collections/task_data/documents", s.endpoint, s.dbID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Appwrite-Project", s.projectID)
	req.Header.Set("X-Appwrite-Key", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return nil, nil }
	defer resp.Body.Close()

	var result struct {
		Documents []map[string]any `json:"documents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil
	}

	var tasks []string
	for _, doc := range result.Documents {
		if dataAttr, ok := doc["data"]; ok {
			var rawEntries []any
			if str, ok := dataAttr.(string); ok {
				json.Unmarshal([]byte(str), &rawEntries)
			} else if slice, ok := dataAttr.([]any); ok {
				rawEntries = slice
			}
			for _, item := range rawEntries {
				if entry, ok := item.(map[string]any); ok {
					if t, ok := entry["Title"].(string); ok { tasks = append(tasks, t) }
					if t, ok := entry["title"].(string); ok { tasks = append(tasks, t) }
				}
			}
		}
	}
	fmt.Printf("✅ [Direct API] Found %d tasks.\n", len(tasks))
	return tasks, nil
}

// GetAboutMe: 'context' collection se current description uthana
func (s *appwriteStore) GetAboutMe(ctx context.Context) (string, error) {
	doc, err := s.db.GetDocument(s.dbID, "context", "me_document_id") // <-- [DASHBOARD SE CONTEXT ID YAHAN PASTE KAR]
	if err != nil {
		return "", nil // Agar nahi milta toh empty string bhej do
	}

	// Fix #3: marshal to access unexported data field
	raw, err := json.Marshal(doc)
	if err != nil {
		return "", nil
	}
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		return "", nil
	}
	if bio, ok := fields["bio"].(string); ok {
		return bio, nil
	}
	return "", nil
}

// UpdatePersonalContext: AI ke naye summary se profile update karna
func (s *appwriteStore) UpdatePersonalContext(ctx context.Context, summary string) error {
	// Fix #4: Pass data as a functional option, not a raw map argument
	_, err := s.db.UpdateDocument(
		s.dbID,
		"context", // <-- [YAHAN BHI SAME CONTEXT ID DAAL]
		"me_document_id",
		s.db.WithUpdateDocumentData(map[string]any{
			"bio":       summary,
			"updatedAt": time.Now().Format(time.RFC3339),
		}),
	)
	return err
}

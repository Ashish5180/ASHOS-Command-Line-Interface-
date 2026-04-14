package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"ashos/cmd"
	"ashos/internal/ai"
	"ashos/internal/core/event"
	"ashos/internal/focus"
	"ashos/internal/note"
	"ashos/internal/slack"
	"ashos/internal/sprint"
	"ashos/internal/storage"
	"ashos/internal/system"
	"ashos/internal/task"
)

func main() {
	// 1. Storage engine
	localStore, err := storage.NewSQLiteStore("data/ashos.db")
	if err != nil {
		fmt.Println("❌ Storage init failed:", err)
		os.Exit(1)
	}

	// --- ☁️ APPWRITE CONFIG (MANUAL) ---
	// Bhai, yahan apni details bhar do:
	endpoint := "http://localhost/v1"
	projectID := "69dcfde50000c78f72a5"
	apiKey := "standard_b8eff7aa5a4460d1758312cc5d0572ef4ea2a8aa3dda5d9acf6df43566c5735d967d3020ca62a6566b28ff033edbd43db58ddc8109ff1e7213ca5384fbaef1e222877bf6f327d3da734259dbda66d3195d1bde7f9bf298be98cfd46d70d2706bd7e3ba158fdf0b492f22d8a57b5fc666fc6dd8bd59f6abbd6fa7982ba32da52c"
	dbID := "69dd003b00190ae7eaa4"

	var store storage.Store = localStore
	if projectID != "YOUR_PROJECT_ID" {
		fmt.Println("☁️  Appwrite Sync Enabled (Hardcoded)!")
		cloudStore := storage.NewAppwriteStore(endpoint, projectID, apiKey, dbID)
		store = storage.NewSyncStore(localStore, cloudStore)
	} else {
		fmt.Println("ℹ️  Running in Local-Only mode (Fill hardcoded variables in main.go to enable sync)")
	}
	// ----------------------------------

	defer store.Close()

	// 2. Event Bus — Central messaging system
	bus := event.NewEventBus()

	// 3. Services setup
	taskRepo := task.NewStoreRepository(store)
	taskService := task.NewService(taskRepo, bus)

	sysService := system.NewService(taskService, store)
	focusManager := focus.NewManager(store, bus)
	noteService := note.NewService(store, bus)
	sprintService := sprint.NewService(store, bus)
	slackService := slack.NewService(store, bus)

	// 4. AI Module — RAG on Your Own Data
	var db *sql.DB
	if s, ok := store.(interface{ GetDB() *sql.DB }); ok {
		db = s.GetDB()
	}

	aiRepo, err := ai.NewSQLiteRepository(db)
	if err != nil {
		fmt.Printf("⚠️  AI module disabled: %v\n", err)
	}
	aiService := ai.NewService(aiRepo, bus, sysService)

	// --- 🧠 Event Subscriptions (Cross-Module Reactions) ---
	bus.Subscribe(event.TaskCompleted{}, func(e any) {
		ev := e.(event.TaskCompleted)
		fmt.Printf("\n✨ [Event] Productivity Win! You finished: %s\n", ev.Title)
	})

	bus.Subscribe(event.FocusEnded{}, func(e any) {
		ev := e.(event.FocusEnded)
		fmt.Printf("\n🔋 [Event] Focus session complete: %v tracked. Recovery mode on!\n", ev.Duration.Truncate(time.Second))
	})

	bus.Subscribe(event.NoteCreated{}, func(e any) {
		fmt.Printf("\n✍️  [Event] Personal insight archived.\n")
	})

	bus.Subscribe(event.SprintEnded{}, func(e any) {
		ev := e.(event.SprintEnded)
		fmt.Printf("\n🌊 [Event] Sprint '%s' closed. Performance analyzed.\n", ev.Title)
	})

	// 5. Dependency container — holds all services for CLI
	app := &cmd.App{
		TaskService:   taskService,
		SystemService: sysService,
		FocusManager:  focusManager,
		NoteService:   noteService,
		SprintService: sprintService,
		AIService:     aiService,
		SlackService:  slackService,
	}

	// 5. Wire commands — inject dependencies into command tree
	cmd.RegisterCommands(app)

	// 6. Execute CLI
	if err := cmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

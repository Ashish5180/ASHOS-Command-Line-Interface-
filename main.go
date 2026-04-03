package main

import (
	"fmt"
	"os"

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
	"database/sql"
	"time"
)

func main() {
	// 1. Storage engine — Ab hum SQLite use kar rahe hain
	store, err := storage.NewSQLiteStore("data/ashos.db")
	if err != nil {
		fmt.Println("❌ Storage init failed:", err)
		os.Exit(1)
	}
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

package main

import (
	"fmt"
	"os"

	"ashos/cmd"
	"ashos/internal/core/event"
	"ashos/internal/focus"
	"ashos/internal/storage"
	"ashos/internal/system"
	"ashos/internal/task"
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

	// 3. Task repository — bridges task module ↔ storage engine
	taskRepo := task.NewStoreRepository(store)

	// 4. Task service — business logic layer
	taskService := task.NewService(taskRepo, bus)

	// 5. System & Focus services
	sysService := system.NewService(taskService, store)
	focusManager := focus.NewManager(store, bus)

	// --- 🧠 Event Subscriptions (Cross-Module Reactions) ---
	// In the future, these can be moved to dedicated modules (e.g., streak_manager)
	bus.Subscribe(event.TaskCompleted{}, func(e any) {
		ev := e.(event.TaskCompleted)
		fmt.Printf("\n✨ [Event] Productivity Win! You finished: %s\n", ev.Title)
	})

	bus.Subscribe(event.FocusEnded{}, func(e any) {
		ev := e.(event.FocusEnded)
		fmt.Printf("\n🔋 [Event] Focus session complete: %v tracked. Recovery mode on!\n", ev.Duration.Truncate(time.Second))
	})

	// 6. Dependency container — holds all services for CLI
	app := &cmd.App{
		TaskService:   taskService,
		SystemService: sysService,
		FocusManager:  focusManager,
	}

	// 5. Wire commands — inject dependencies into command tree
	cmd.RegisterCommands(app)

	// 6. Execute CLI
	if err := cmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"ashos/cmd"
	"ashos/internal/focus"
	"ashos/internal/storage"
	"ashos/internal/system"
	"ashos/internal/task"
)

func main() {
	// ---------------------------------------------------------------
	// COMPOSITION ROOT — Sab dependencies yahan banti hain.
	// Sirf yeh file jaanti hai ki JSON use ho raha hai.
	// Agar kal SQLite lagana ho, toh sirf YEH file change hogi.
	// ---------------------------------------------------------------

	// 1. Storage engine — JSON files "data/" folder mein store hongi
	store, err := storage.NewJSONStore("data")
	if err != nil {
		fmt.Println("❌ Storage init failed:", err)
		os.Exit(1)
	}
	defer store.Close()

	// 2. Task repository — bridges task module ↔ storage engine
	taskRepo := task.NewStoreRepository(store)

	// 3. Task service — business logic layer
	taskService := task.NewService(taskRepo)

	// 4. System & Focus services
	sysService := system.NewService(taskService, store)
	focusManager := focus.NewManager(store)

	// 5. Dependency container — holds all services for CLI
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

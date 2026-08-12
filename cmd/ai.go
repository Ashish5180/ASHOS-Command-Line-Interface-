package cmd

import (
	"ashos/internal/storage"
	"ashos/internal/ui"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func aiCmd(app *App) *cobra.Command {
	ai := &cobra.Command{
		Use:   "ai",
		Short: "🧠 AI Brain commands (RAG)",
	}

	var local bool
	ask := &cobra.Command{
		Use:   "ask [query]",
		Short: "Ask a question about your personal OS data",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if local {
				app.AIService.SetLocalMode(true)
			}
			query := strings.Join(args, " ")

			fmt.Println(ui.SubTitleStyle.Render("🧠 ASHOS AI BRAIN IS THINKING..."))
			fmt.Println(ui.GrayStyle.Render("Query: ") + query)

			resp, records, err := app.AIService.Ask(context.Background(), query)
			if err != nil {
				return err
			}

			// Show Retrieval Scores
			if len(records) > 0 {
				fmt.Println("\n" + ui.SubTitleStyle.Foreground(ui.SecondaryColor).Render("Retrieved Context (Relevance Scores):"))
				for _, r := range records {
					scorePct := int(r.Score * 100)
					if scorePct > 100 {
						scorePct = 99 // Cap for display
					}
					fmt.Printf("  → %-40s [%d%% match]\n", truncateString(r.Content, 40), scorePct)
				}
			}

			header := ui.TitleStyle.Render("🤖 AI RESPONSE")
			content := ui.GrayStyle.Foreground(ui.WhiteColor).Render(resp)

			fmt.Printf("\n%s\n\n%s\n", header, content)
			return nil
		},
	}
	ask.Flags().BoolVarP(&local, "local", "l", false, "Use local LLM (Ollama)")

	standup := &cobra.Command{
		Use:   "standup",
		Short: "Auto-generate your daily standup",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ui.SubTitleStyle.Render("📋 GATHERING DATA FOR STANDUP..."))

			resp, err := app.AIService.GenerateStandup(context.Background())
			if err != nil {
				return err
			}

			fmt.Printf("\n%s\n", resp)

			fmt.Print("\nCopy to clipboard? (y/n): ")
			var input string
			fmt.Scanln(&input)

			if strings.ToLower(input) == "y" {
				// Use pbcopy on macOS
				copyCmd := exec.Command("pbcopy")
				copyCmd.Stdin = strings.NewReader(resp)
				if err := copyCmd.Run(); err != nil {
					fmt.Println(ui.ErrorStyle.Render("Failed to copy to clipboard"))
				} else {
					fmt.Println(ui.SuccessStyle.Render("✅ Copied to clipboard!"))
				}
			}

			return nil
		},
	}

	suggest := &cobra.Command{
		Use:   "suggest",
		Short: "Get suggestions on what to work on next",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ui.SubTitleStyle.Render("💡 ANALYZING PATTERNS..."))

			resp, err := app.AIService.SuggestNextTask(context.Background())
			if err != nil {
				return err
			}

			header := ui.TitleStyle.Background(ui.AccentColor).Render("🚀 SMART SUGGESTION")
			fmt.Printf("\n%s\n\n%s\n", header, resp)
			return nil
		},
	}

	digest := &cobra.Command{
		Use:   "daily-digest",
		Short: "Generate your daily productivity intelligence digest",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ui.SubTitleStyle.Render("📊 GENERATING DAILY INTELLIGENCE..."))

			resp, err := app.AIService.GenerateDailyDigest(context.Background())
			if err != nil {
				return err
			}

			fmt.Printf("\n%s\n", resp)
			return nil
		},
	}

	var syncWeeklyCmd = &cobra.Command{
		Use:   "weekly-sync",
		Short: "Sync and summarize your weekly progress locally and to the cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("⏳ Starting Appwrite weekly sync and brain update...")

			// Appwrite config (Inhe ek jagah define karna better hota hai, but filhal main.go wale use kar rahe hain)
			endpoint := "http://localhost/v1"
			projectID := "69dcfde50000c78f72a5"
			apiKey := "standard_b8eff7aa5a4460d1758312cc5d0572ef4ea2a8aa3dda5d9acf6df43566c5735d967d3020ca62a6566b28ff033edbd43db58ddc8109ff1e7213ca5384fbaef1e222877bf6f327d3da734259dbda66d3195d1bde7f9bf298be98cfd46d70d2706bd7e3ba158fdf0b492f22d8a57b5fc666fc6dd8bd59f6abbd6fa7982ba32da52c"
			dbID := "69dd003b00190ae7eaa4"

			cloudStore := storage.NewAppwriteStore(endpoint, projectID, apiKey, dbID)

			// Weekly Sync service initialize karein
			// Note: cloudStore ko DataRepository interface mein cast kar rahe hain
			syncService := storage.NewWeeklySyncService(cloudStore.(storage.DataRepository), app.AIService)

			// Sync run karein
			err := syncService.RunWeeklySync()
			if err != nil {
				fmt.Printf("❌ Sync failed: %v\n", err)
				return err
			}

			return nil
		},
	}
	var reingestAll bool
	var reingestType string
	reingest := &cobra.Command{
		Use:   "reingest",
		Short: "Force re-ingestion of data into AI vector store",
		RunE: func(cmd *cobra.Command, args []string) error {
			if local {
				app.AIService.SetLocalMode(true)
			}
			ctx := context.Background()

			if reingestAll {
				fmt.Println("🧠 Clearing vector store for full re-sync...")
				app.AIService.Reset(ctx)
			}

			if reingestAll || reingestType == "tasks" {
				fmt.Println("🧠 Re-ingesting TASKS...")
				tasks, _ := app.TaskService.ListTasks()
				for _, t := range tasks {
					app.AIService.IngestTask(ctx, t.ID, t.Title)
				}
			}

			if reingestAll || reingestType == "notes" {
				fmt.Println("🧠 Re-ingesting NOTES...")
				notes, _ := app.NoteService.ListNotes(ctx)
				for _, n := range notes {
					app.AIService.IngestNote(ctx, n.ID, n.Content)
				}
			}

			if reingestAll || reingestType == "focus" {
				fmt.Println("🧠 Re-ingesting FOCUS SESSIONS...")
				sessions, _ := app.FocusManager.ListSessions(ctx)
				for _, s := range sessions {
					app.AIService.IngestFocusSession(ctx, s.StartTime, s.Duration, s.Summary)
				}
			}

			if reingestAll || reingestType == "sprints" {
				fmt.Println("🧠 Re-ingesting SPRINTS...")
				sprints, _ := app.SprintService.ListSprints(ctx)
				for _, s := range sprints {
					app.AIService.IngestSprint(ctx, s.ID, s.Title, s.Summary)
				}
			}

			fmt.Println("✅ Re-ingestion complete!")
			return nil
		},
	}
	reingest.Flags().BoolVarP(&reingestAll, "all", "a", false, "Re-ingest everything")
	reingest.Flags().StringVarP(&reingestType, "type", "t", "", "Re-ingest specific type (tasks, notes, focus, sprints)")
	reingest.Flags().BoolVarP(&local, "local", "l", false, "Use local LLM (Ollama) for re-ingestion")

	// -------------------------------------------------------------------------
	// ash ai ingest-about-me
	// -------------------------------------------------------------------------
	ingestAboutMe := &cobra.Command{
		Use:   "ingest-about-me",
		Short: "🧬 Embed your personal context (about_me.json) into the AI brain",
		Long: `Reads data/about_me.json and embeds every entry into the RAG vector store.
The AI will use these entries to answer questions about Ashish — identity,
quirks, projects, goals, communication style — with personal accuracy.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// --local flag support (shares the var declared earlier in scope)
			if local {
				app.AIService.SetLocalMode(true)
			}

			// Read the JSON file
			path := "data/about_me.json"
			raw, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("❌ Cannot read %s: %w\nRun from the project root directory.", path, err)
			}

			// Parse structure
			type Entry struct {
				Key      string `json:"key"`
				Category string `json:"category"`
				Content  string `json:"content"`
			}
			var payload struct {
				Entries []Entry `json:"entries"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				return fmt.Errorf("❌ Invalid JSON in about_me.json: %w", err)
			}

			if len(payload.Entries) == 0 {
				return fmt.Errorf("❌ No entries found in about_me.json")
			}

			fmt.Println(ui.TitleStyle.Render(" 🧬 ASHOS ABOUT-ME INGESTION "))
			fmt.Printf("\n📂 Found %d entries. Embedding into AI brain...\n\n", len(payload.Entries))

			success, skipped, failed := 0, 0, 0
			for _, e := range payload.Entries {
				err := app.AIService.IngestAboutMe(ctx, e.Key, e.Category, e.Content)
				if err != nil {
					fmt.Printf("  ❌ Failed [%s]: %v\n", e.Key, err)
					failed++
				} else {
					// IngestAboutMe prints skip/success internally
					if strings.Contains(e.Key, "skip") {
						skipped++
					} else {
						success++
					}
				}
			}

			fmt.Println()
			fmt.Println(ui.BoxStyle.Render(fmt.Sprintf(
				"✅ Ingested : %d\n⏭️  Skipped  : %d\n❌ Failed   : %d\n\n🎉 AI brain now knows you personally!",
				success+skipped, 0, failed,
			)))
			return nil
		},
	}
	ingestAboutMe.Flags().BoolVarP(&local, "local", "l", false, "Use local Ollama for embeddings")

	ai.AddCommand(ask, standup, suggest, digest, syncWeeklyCmd, reingest, ingestAboutMe)
	return ai
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

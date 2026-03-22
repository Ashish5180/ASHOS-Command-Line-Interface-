package cmd

import (
	"ashos/internal/ui"
	"context"
	"fmt"
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

	var reingestAll bool
	var reingestType string
	reingest := &cobra.Command{
		Use:   "reingest",
		Short: "Force re-ingestion of data into AI vector store",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			
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

	ai.AddCommand(ask, standup, suggest, digest, reingest)
	return ai
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

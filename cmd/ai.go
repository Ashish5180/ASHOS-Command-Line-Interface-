package cmd

import (
	"ashos/internal/ui"
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func aiCmd(app *App) *cobra.Command {
	ai := &cobra.Command{
		Use:   "ai",
		Short: "🧠 AI Brain commands (RAG)",
	}

	ask := &cobra.Command{
		Use:   "ask [query]",
		Short: "Ask a question about your personal OS data",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			
			// Premium Thinking Message
			fmt.Println(ui.SubTitleStyle.Render("🧠 ASHOS AI BRAIN IS THINKING..."))
			fmt.Println(ui.GrayStyle.Render("Query: ") + query)
			
			resp, err := app.AIService.Ask(context.Background(), query)
			if err != nil {
				return err
			}
			
			// Premium Balanced Layout (No Box)
			header := ui.TitleStyle.Render("🤖 AI RESPONSE")
			content := ui.GrayStyle.Foreground(ui.WhiteColor).Render(resp)
			
			fmt.Printf("\n%s\n\n%s\n", header, content)
			return nil
		},
	}

	standup := &cobra.Command{
		Use:   "standup",
		Short: "Auto-generate your daily standup",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ui.SubTitleStyle.Render("📋 GATHERING DATA FOR STANDUP..."))
			
			resp, err := app.AIService.GenerateStandup(context.Background())
			if err != nil {
				return err
			}
			
			header := ui.TitleStyle.Background(ui.SecondaryColor).Render("✨ DAILY STANDUP REPORT")
			fmt.Printf("\n%s\n\n%s\n", header, resp)
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

	ai.AddCommand(ask, standup, suggest)
	return ai
}

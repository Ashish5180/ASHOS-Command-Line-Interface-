package cmd

import (
	"ashos/internal/ui"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)

// NewFocusCommand returns the 'focus' command.
func NewFocusCommand(app *App) *cobra.Command {
	focusCmd := &cobra.Command{
		Use:   "focus",
		Short: "Manage focus sessions",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a focus session",
		Run: func(cmd *cobra.Command, args []string) {
			// 🧠 TECHNIQUE: Interrupt Handling
			// Catch Ctrl+C to stop focus session gracefully.
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			app.FocusManager.StartSession()

			fmt.Println("Press Ctrl+C to end focus session.")

			// Wait for interrupt
			<-sigs
			fmt.Print("\n📝 What did you focus on? (Optional): ")
			
			scanner := bufio.NewScanner(os.Stdin)
			var summary string
			if scanner.Scan() {
				summary = scanner.Text()
			}
			
			if strings.TrimSpace(summary) == "" {
				summary = "No summary provided."
			}

			app.FocusManager.StopSession(summary)
			fmt.Printf("Total Focus Time: %v\n", app.FocusManager.GetStats())
		},
	}

	analyticsCmd := &cobra.Command{
		Use:   "analytics",
		Short: "Show focus stats and productivity heatmap",
		Run: func(cmd *cobra.Command, args []string) {
			report, err := app.FocusManager.GetAnalyticsReport(7)
			if err != nil {
				fmt.Printf("❌ Failed to load analytics: %v\n", err)
				return
			}

			fmt.Println("\n" + ui.TitleStyle.Render("📊 FOCUS ANALYTICS - LAST 7 DAYS"))

			// --- 📈 Trend Graph (asciigraph) ---
			data := make([]float64, len(report.DailyStats))
			for i, s := range report.DailyStats {
				data[i] = s.Duration.Minutes()
			}
			
			if len(data) > 0 {
				graph := asciigraph.Plot(data, 
					asciigraph.Height(8), 
					asciigraph.Width(40),
					asciigraph.Caption("Focus Intensity (Minutes)"),
				)
				fmt.Println(ui.GrayStyle.Render(graph) + "\n")
			}

			// --- 🔥 Heatmap (Bar Chart) ---
			maxSeconds := 0.0
			for _, s := range report.DailyStats {
				if s.Duration.Seconds() > maxSeconds {
					maxSeconds = s.Duration.Seconds()
				}
			}

			for _, s := range report.DailyStats {
				day := s.Date.Format("Mon")
				duration := s.Duration.Truncate(time.Minute)
				
				// Standardize bar length (max 30 chars)
				barLen := 0
				if maxSeconds > 0 {
					barLen = int((s.Duration.Seconds() / maxSeconds) * 30)
				}
				bar := strings.Repeat("█", barLen)
				
				// Color based on activity
				color := ui.GrayColor
				suffix := ""
				if s.Duration >= 4*time.Hour {
					color = ui.SecondaryColor
					suffix = " 🏆"
				} else if s.Duration >= 1*time.Hour {
					color = ui.PrimaryColor
					suffix = " 🔥"
				}

				dayStyle := lipgloss.NewStyle().Width(5).Foreground(ui.WhiteColor).Bold(true)
				barStyle := lipgloss.NewStyle().Foreground(color)

				fmt.Printf("%s %s %s%s\n", 
					dayStyle.Render(day), 
					barStyle.Render(bar), 
					ui.GrayStyle.Render(duration.String()),
					suffix)
			}

			fmt.Println("\n" + ui.BoxStyle.Render(fmt.Sprintf(
				"Deep Work Score: %s\nCurrent Streak:   %s days",
				ui.HighlightStyle.Render(fmt.Sprintf("%d/100", report.DeepWorkScore)),
				ui.HighlightStyle.Render(fmt.Sprintf("%d", report.Streak)),
			)))
		},
	}

	focusCmd.AddCommand(startCmd)
	focusCmd.AddCommand(analyticsCmd)
	return focusCmd
}

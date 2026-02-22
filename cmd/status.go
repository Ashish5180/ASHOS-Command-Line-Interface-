package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewStatusCommand returns the 'status' command.
func NewStatusCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show system status dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			status := app.SystemService.GetStatus()

			fmt.Println("╔══════════════════════╗")
			fmt.Println("║       ASHOS          ║")
			fmt.Println("╚══════════════════════╝")
			fmt.Println()
			fmt.Printf("🧠 Focus Today: %v\n", status.FocusToday)
			fmt.Printf("📌 Pending Tasks: %d\n", status.PendingTasks)
			fmt.Printf("💻 Uptime: %v\n", status.Uptime)
			fmt.Println("🔥 Systems compound daily.")
		},
	}
}

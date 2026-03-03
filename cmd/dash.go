package cmd

import (
	"ashos/internal/ui"

	"github.com/spf13/cobra"
)

// newDashCommand returns the 'dash' command.
func newDashCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "dash",
		Short: "Launch the interactive dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ui.RunDashboard(app.TaskService, app.SystemService, app.FocusManager)
		},
	}
}

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
			app.FocusManager.StopSession()
			fmt.Printf("Total Focus Time: %v\n", app.FocusManager.GetStats())
		},
	}

	focusCmd.AddCommand(startCmd)
	return focusCmd
}

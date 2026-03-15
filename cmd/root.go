package cmd

import (
	"ashos/internal/ai"
	"ashos/internal/focus"
	"ashos/internal/system"
	"ashos/internal/task"
	"fmt"

	"github.com/spf13/cobra"
)

// App is the dependency container.
// main.go mein banata hai, commands mein inject hota hai.
// Koi bhi command file directly storage/repo nahi banayegi — sab kuch yahan se milega.
type App struct {
	TaskService   *task.Service
	SystemService *system.Service
	FocusManager  *focus.Manager
	AIService     ai.Service
}

// rootCmd is the base "ash" command.
var rootCmd = &cobra.Command{
	Use:   "ash",
	Short: "ASHOS - Your Personal Command Center",
	Long: `
⚙️ ASHOS - Personal CLI Operating System
Your digital cockpit.
Control tasks, notes, focus and system insights.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 ASHOS is online.")
		fmt.Println("Type 'ash --help' to see available commands.")
	},
}

// RegisterCommands wires all module commands with injected dependencies.
// main.go calls this AFTER creating App — ensures DI flows top-down.
//
// Kyon alag function hai?
// Kyunki rootCmd global hai (Cobra pattern), lekin dependencies runtime pe banti hain.
// RegisterCommands wo bridge hai jo runtime dependencies ko static command tree se jodta hai.
func RegisterCommands(app *App) {
	rootCmd.AddCommand(NewTaskCommand(app))
	rootCmd.AddCommand(NewStatusCommand(app))
	rootCmd.AddCommand(NewFocusCommand(app))
	rootCmd.AddCommand(newDashCommand(app))
	rootCmd.AddCommand(aiCmd(app))
}

// Execute runs the root command. Called from main.go.
func Execute() error {
	return rootCmd.Execute()
}

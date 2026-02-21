package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

func Execute() error {
	return rootCmd.Execute()
}

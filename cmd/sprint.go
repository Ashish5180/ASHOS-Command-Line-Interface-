package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewSprintCommand(app *App) *cobra.Command {
	sprintCmd := &cobra.Command{
		Use:   "sprint",
		Short: "Manage productivity sprints",
	}

	endCmd := &cobra.Command{
		Use:   "end",
		Short: "Summarize and end a work sprint",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Print("🏆 Sprint Title: ")
			var title string
			if scanner.Scan() {
				title = strings.TrimSpace(scanner.Text())
			}

			fmt.Print("📝 Sprint Review Summary: ")
			var summary string
			if scanner.Scan() {
				summary = strings.TrimSpace(scanner.Text())
			}

			fmt.Print("✅ Tasks Completed (number): ")
			var tasksStr string
			if scanner.Scan() {
				tasksStr = scanner.Text()
			}
			tasksDone, _ := strconv.Atoi(tasksStr)

			fmt.Print("⏱️ Total Focus Time (e.g., 2h30m): ")
			var focusStr string
			if scanner.Scan() {
				focusStr = scanner.Text()
			}
			focusTime, _ := time.ParseDuration(focusStr)

			err := app.SprintService.EndSprint(cmd.Context(), title, summary, tasksDone, focusTime)
			if err != nil {
				fmt.Printf("❌ Failed to end sprint: %v\n", err)
				return
			}
			fmt.Println("🌊 Sprint summarized and archived in the brain.")
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List sprint history",
		Run: func(cmd *cobra.Command, args []string) {
			sprints, err := app.SprintService.ListSprints(cmd.Context())
			if err != nil {
				fmt.Printf("❌ Failed to list sprints: %v\n", err)
				return
			}
			if len(sprints) == 0 {
				fmt.Println("No sprints found.")
				return
			}
			for _, s := range sprints {
				fmt.Printf("[%d] %s: %s (%d tasks, %v focus)\n", s.ID, s.Title, s.Summary, s.TasksDone, s.FocusTime)
			}
		},
	}

	sprintCmd.AddCommand(endCmd)
	sprintCmd.AddCommand(listCmd)
	return sprintCmd
}

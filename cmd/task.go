package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// NewTaskCommand creates the `ash task` command group.
// App inject hota hai — ye function KABHI directly repo/storage nahi banayega.
func NewTaskCommand(app *App) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:   "task",
		Short: "Manage your tasks",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	taskCmd.AddCommand(newAddTaskCmd(app))
	taskCmd.AddCommand(newListTaskCmd(app))
	taskCmd.AddCommand(newDoneTaskCmd(app))
	taskCmd.AddCommand(newDeleteTaskCmd(app))

	return taskCmd
}

// ---------------------------------------------------------------------------
// ash task add "Buy groceries"
// ---------------------------------------------------------------------------
func newAddTaskCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "add [title...]",
		Short: "Add a new task",
		Example: `  ash task add Buy groceries
  ash task add "Read 10 pages of Go book"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := strings.Join(args, " ")

			id, err := app.TaskService.AddTask(title)
			if err != nil {
				return err
			}

			fmt.Printf("✅ Task #%d added: %s\n", id, title)
			return nil
		},
	}
}

// ---------------------------------------------------------------------------
// ash task list
// ---------------------------------------------------------------------------
func newListTaskCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, err := app.TaskService.ListTasks()
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				fmt.Println("📭 No tasks found. Add one with: ash task add <title>")
				return nil
			}

			fmt.Println("📝 Your Tasks:")
			fmt.Println(strings.Repeat("─", 40))
			for _, t := range tasks {
				status := "[ ]"
				if t.Completed {
					status = "[✓]"
				}
				fmt.Printf("  %s #%d — %s\n", status, t.ID, t.Title)
			}
			fmt.Println(strings.Repeat("─", 40))

			pending, _ := app.TaskService.PendingCount()
			fmt.Printf("📊 Pending: %d | Total: %d\n", pending, len(tasks))
			return nil
		},
	}
}

// ---------------------------------------------------------------------------
// ash task done 3
// ---------------------------------------------------------------------------
func newDoneTaskCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:     "done [id]",
		Short:   "Mark a task as completed",
		Example: "  ash task done 3",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %q — must be a number", args[0])
			}

			if err := app.TaskService.CompleteTask(id); err != nil {
				return err
			}

			fmt.Printf("🎉 Task #%d marked as done!\n", id)
			return nil
		},
	}
}

// ---------------------------------------------------------------------------
// ash task delete 3
// ---------------------------------------------------------------------------
func newDeleteTaskCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [id]",
		Short:   "Delete a task permanently",
		Example: "  ash task delete 3",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID: %q — must be a number", args[0])
			}

			if err := app.TaskService.DeleteTask(id); err != nil {
				return err
			}

			fmt.Printf("🗑️  Task #%d deleted!\n", id)
			return nil
		},
	}
}

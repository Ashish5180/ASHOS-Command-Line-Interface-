package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewNoteCommand(app *App) *cobra.Command {
	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage journal entries and personal notes",
	}

	addCmd := &cobra.Command{
		Use:   "add <text>",
		Short: "Add a new journal entry",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			content := strings.Join(args, " ")
			err := app.NoteService.AddNote(cmd.Context(), content)
			if err != nil {
				fmt.Printf("❌ Failed to add note: %v\n", err)
				return
			}
			fmt.Println("📝 Note added to your brain.")
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List recent notes",
		Run: func(cmd *cobra.Command, args []string) {
			notes, err := app.NoteService.ListNotes(cmd.Context())
			if err != nil {
				fmt.Printf("❌ Failed to list notes: %v\n", err)
				return
			}
			if len(notes) == 0 {
				fmt.Println("No notes found.")
				return
			}
			for _, n := range notes {
				fmt.Printf("[%d] %s (%s)\n", n.ID, n.Content, n.CreatedAt.Format("2006-01-02 15:04"))
			}
		},
	}

	noteCmd.AddCommand(addCmd)
	noteCmd.AddCommand(listCmd)
	return noteCmd
}

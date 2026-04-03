package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newSlackCmd(app *App) *cobra.Command {
	slackCmd := &cobra.Command{
		Use:   "slack",
		Short: "🪸 Control and access Slack",
		Long:  `⚙️ Connect ASHOS to your Slack workspace and channels.`,
	}

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "🔗 Link your Slack Bot Token and Channel ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, _ := cmd.Flags().GetString("token")
			channel, _ := cmd.Flags().GetString("channel")

			if token == "" || channel == "" {
				return fmt.Errorf("both --token and --channel are required")
			}

			err := app.SlackService.Connect(context.Background(), token, channel)
			if err != nil {
				return err
			}

			fmt.Println("✅ Slack successfully linked. Ready to interact.")
			return nil
		},
	}
	connectCmd.Flags().StringP("token", "t", "", "Slack Bot User OAuth Token (xoxb-...)")
	connectCmd.Flags().StringP("channel", "c", "", "Slack Channel ID")

	sendCmd := &cobra.Command{
		Use:   "send [message]",
		Short: "📤 Send a message to the linked channel",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := strings.Join(args, " ")
			err := app.SlackService.SendMessage(context.Background(), msg)
			if err != nil {
				return err
			}
			fmt.Println("🚀 Message sent to Slack.")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "📥 Fetch recent messages from the linked channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			msgs, err := app.SlackService.GetRecentMessages(context.Background(), 5)
			if err != nil {
				return err
			}

			fmt.Printf("📥 Last %d messages from Slack:\n", len(msgs))
			for i, m := range msgs {
				fmt.Printf("[%d] %s: %s\n", i+1, m.User, m.Text)
			}
			return nil
		},
	}

	slackCmd.AddCommand(connectCmd)
	slackCmd.AddCommand(sendCmd)
	slackCmd.AddCommand(listCmd)

	return slackCmd
}

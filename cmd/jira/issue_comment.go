package jira

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var issueCommentCmd = &cobra.Command{
	Use:   "comment [issue-key] [body...]",
	Short: "Add a comment to an issue",
	Args:  cobra.MinimumNArgs(2),
	Example: `  orbit jira issue comment PROJ-123 "This is fixed now"
  orbit jira issue comment PROJ-123 "Looks good to me"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		body := strings.Join(args[1:], " ")

		if err := client.AddComment(key, body); err != nil {
			return err
		}

		fmt.Printf("Comment added to %s\n", key)
		return nil
	},
}

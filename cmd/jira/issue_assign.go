package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueAssignCmd = &cobra.Command{
	Use:   "assign [issue-key] [assignee]",
	Short: "Assign or unassign an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit jira issue assign PROJ-123 john
  orbit jira issue assign PROJ-123 x   # unassign`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		assignee := args[1]

		if assignee == "x" {
			if err := client.UnassignIssue(key); err != nil {
				return err
			}
			fmt.Printf("Unassigned %s\n", key)
		} else {
			if err := client.AssignIssue(key, assignee); err != nil {
				return err
			}
			fmt.Printf("Assigned %s to %s\n", key, assignee)
		}
		return nil
	},
}

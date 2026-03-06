package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueUnlinkCmd = &cobra.Command{
	Use:   "unlink [inward-key] [outward-key]",
	Short: "Remove a link between two issues",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit jira issue unlink PROJ-100 PROJ-200`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if err := client.UnlinkIssues(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Unlinked %s from %s\n", args[0], args[1])
		return nil
	},
}

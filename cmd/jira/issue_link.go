package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueLinkCmd = &cobra.Command{
	Use:   "link [inward-key] [outward-key] [link-type]",
	Short: "Link two issues together",
	Args:  cobra.ExactArgs(3),
	Example: `  orbit jira issue link PROJ-100 PROJ-200 Blocks
  orbit jira issue link PROJ-100 PROJ-200 Duplicates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if err := client.LinkIssues(args[0], args[1], args[2]); err != nil {
			return err
		}

		fmt.Printf("Linked %s -> %s (%s)\n", args[0], args[1], args[2])
		return nil
	},
}

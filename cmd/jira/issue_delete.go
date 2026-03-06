package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueDeleteOpts struct {
	cascade bool
}

var issueDeleteCmd = &cobra.Command{
	Use:     "delete [issue-key]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an issue",
	Args:    cobra.ExactArgs(1),
	Example: `  orbit jira issue delete PROJ-123
  orbit jira issue delete PROJ-123 --cascade`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if err := client.DeleteIssue(args[0], issueDeleteOpts.cascade); err != nil {
			return err
		}

		fmt.Printf("Deleted %s\n", args[0])
		return nil
	},
}

func init() {
	issueDeleteCmd.Flags().BoolVar(&issueDeleteOpts.cascade, "cascade", false, "delete subtasks too")
}

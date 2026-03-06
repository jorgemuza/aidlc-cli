package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueWorklogOpts struct {
	comment string
}

var issueWorklogCmd = &cobra.Command{
	Use:   "worklog [issue-key] [time-spent]",
	Short: "Log time spent on an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit jira issue worklog PROJ-123 "2h 30m"
  orbit jira issue worklog PROJ-123 "1d" --comment "Code review"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if err := client.AddWorklog(args[0], args[1], issueWorklogOpts.comment); err != nil {
			return err
		}

		fmt.Printf("Logged %s on %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	issueWorklogCmd.Flags().StringVar(&issueWorklogOpts.comment, "comment", "", "comment for worklog entry")
}

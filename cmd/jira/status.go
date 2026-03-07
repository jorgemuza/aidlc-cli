package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Manage Jira workflow statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var statusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflow statuses",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		statuses, err := client.ListStatuses()
		if err != nil {
			return err
		}

		for _, s := range statuses {
			category := ""
			if s.StatusCategory != nil {
				category = fmt.Sprintf(" [%s]", s.StatusCategory.Name)
			}
			fmt.Printf("%-10s %-30s%s\n", s.ID, s.Name, category)
		}
		return nil
	},
}

var issueTypeListCmd = &cobra.Command{
	Use:   "issuetype-list",
	Short: "List all issue types",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		issueTypes, err := client.ListIssueTypes()
		if err != nil {
			return err
		}

		for _, it := range issueTypes {
			fmt.Printf("%-10s %s\n", it.ID, it.Name)
		}
		return nil
	},
}

func init() {
	statusCmd.AddCommand(statusListCmd)
}

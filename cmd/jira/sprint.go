package jira

import (
	"fmt"
	"strconv"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var sprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "Manage sprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var sprintListOpts struct {
	boardID    int
	state      string
	maxResults int
}

var sprintListCmd = &cobra.Command{
	Use:     "list [sprint-id]",
	Aliases: []string{"ls"},
	Short:   "List sprints or issues in a sprint",
	Args:    cobra.MaximumNArgs(1),
	Example: `  orbit jira sprint list --board-id 42
  orbit jira sprint list 123
  orbit jira sprint list --board-id 42 --state active`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			// List issues in sprint
			sprintID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid sprint ID: %s", args[0])
			}

			result, err := client.GetSprintIssues(sprintID)
			if err != nil {
				return err
			}

			headers := []string{"KEY", "TYPE", "STATUS", "PRIORITY", "ASSIGNEE", "SUMMARY"}
			rowFn := func() [][]string {
				var rows [][]string
				for _, issue := range result.Issues {
					assignee := ""
					if issue.Fields.Assignee != nil {
						assignee = issue.Fields.Assignee.DisplayName
					}
					rows = append(rows, []string{
						issue.Key, issue.Fields.IssueType.Name, issue.Fields.Status.Name,
						issue.Fields.Priority.Name, assignee, issue.Fields.Summary,
					})
				}
				return rows
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Issues in sprint: %d of %d\n", len(result.Issues), result.Total)
			return output.Print(cmdutil.OutputFormat(cmd), result.Issues, headers, rowFn)
		}

		// List sprints for a board
		if sprintListOpts.boardID == 0 {
			return fmt.Errorf("--board-id is required when listing sprints")
		}

		sprints, err := client.ListSprints(sprintListOpts.boardID, sprintListOpts.state)
		if err != nil {
			return err
		}

		headers := []string{"ID", "NAME", "STATE", "START", "END"}
		rowFn := func() [][]string {
			var rows [][]string
			for _, s := range sprints {
				rows = append(rows, []string{
					fmt.Sprintf("%d", s.ID), s.Name, s.State, s.StartDate, s.EndDate,
				})
			}
			return rows
		}
		return output.Print(cmdutil.OutputFormat(cmd), sprints, headers, rowFn)
	},
}

var sprintAddCmd = &cobra.Command{
	Use:   "add [sprint-id] [issue-keys...]",
	Short: "Add issues to a sprint (max 50)",
	Args:  cobra.MinimumNArgs(2),
	Example: `  orbit jira sprint add 42 PROJ-101 PROJ-102`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		sprintID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %s", args[0])
		}

		issueKeys := args[1:]
		if len(issueKeys) > 50 {
			return fmt.Errorf("max 50 issues at once, got %d", len(issueKeys))
		}

		if err := client.AddIssuesToSprint(sprintID, issueKeys); err != nil {
			return err
		}

		fmt.Printf("Added %d issues to sprint %d\n", len(issueKeys), sprintID)
		return nil
	},
}

func init() {
	sprintCmd.AddCommand(sprintListCmd)
	sprintCmd.AddCommand(sprintAddCmd)

	sprintListCmd.Flags().IntVar(&sprintListOpts.boardID, "board-id", 0, "board ID (required for listing sprints)")
	sprintListCmd.Flags().StringVar(&sprintListOpts.state, "state", "", "filter by state: future, active, closed")
	sprintListCmd.Flags().IntVar(&sprintListOpts.maxResults, "max-results", 50, "max results")
}

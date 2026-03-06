package jira

import (
	"fmt"
	"strings"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var epicCmd = &cobra.Command{
	Use:   "epic",
	Short: "Manage epics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var epicListOpts struct {
	project    string
	status     []string
	maxResults int
}

var epicListCmd = &cobra.Command{
	Use:     "list [epic-key]",
	Aliases: []string{"ls"},
	Short:   "List epics or issues within an epic",
	Args:    cobra.MaximumNArgs(1),
	Example: `  orbit jira epic list --project PROJ
  orbit jira epic list PROJ-50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			// List issues in epic
			jql := fmt.Sprintf("parent = %s ORDER BY created DESC", args[0])
			result, err := client.SearchIssues(jql, 0, epicListOpts.maxResults)
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
			fmt.Fprintf(cmd.ErrOrStderr(), "Issues in epic %s: %d of %d\n", args[0], len(result.Issues), result.Total)
			return output.Print(cmdutil.OutputFormat(cmd), result.Issues, headers, rowFn)
		}

		// List epics
		var clauses []string
		clauses = append(clauses, "issuetype = Epic")
		if epicListOpts.project != "" {
			clauses = append(clauses, fmt.Sprintf("project = %q", epicListOpts.project))
		}
		if len(epicListOpts.status) > 0 {
			quoted := make([]string, len(epicListOpts.status))
			for i, s := range epicListOpts.status {
				quoted[i] = fmt.Sprintf("%q", s)
			}
			clauses = append(clauses, fmt.Sprintf("status in (%s)", strings.Join(quoted, ", ")))
		}
		jql := strings.Join(clauses, " AND ") + " ORDER BY created DESC"

		result, err := client.SearchIssues(jql, 0, epicListOpts.maxResults)
		if err != nil {
			return err
		}

		headers := []string{"KEY", "STATUS", "PRIORITY", "SUMMARY"}
		rowFn := func() [][]string {
			var rows [][]string
			for _, issue := range result.Issues {
				rows = append(rows, []string{
					issue.Key, issue.Fields.Status.Name,
					issue.Fields.Priority.Name, issue.Fields.Summary,
				})
			}
			return rows
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Showing %d of %d epics\n", len(result.Issues), result.Total)
		return output.Print(cmdutil.OutputFormat(cmd), result.Issues, headers, rowFn)
	},
}

var epicCreateOpts struct {
	project     string
	name        string
	summary     string
	description string
	priority    string
	labels      []string
	components  []string
}

var epicCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new epic",
	Example: `  orbit jira epic create --project PROJ --name "Q1 Auth" --summary "Revamp auth system"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Epic create is just issue create with type=Epic
		issueCreateOpts.project = epicCreateOpts.project
		issueCreateOpts.issueType = "Epic"
		issueCreateOpts.summary = epicCreateOpts.summary
		if issueCreateOpts.summary == "" {
			issueCreateOpts.summary = epicCreateOpts.name
		}
		issueCreateOpts.description = epicCreateOpts.description
		issueCreateOpts.priority = epicCreateOpts.priority
		issueCreateOpts.labels = epicCreateOpts.labels
		issueCreateOpts.components = epicCreateOpts.components
		return issueCreateCmd.RunE(cmd, args)
	},
}

var epicAddCmd = &cobra.Command{
	Use:   "add [epic-key] [issue-keys...]",
	Short: "Add issues to an epic (max 50)",
	Args:  cobra.MinimumNArgs(2),
	Example: `  orbit jira epic add PROJ-50 PROJ-101 PROJ-102 PROJ-103`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		epicKey := args[0]
		issueKeys := args[1:]
		if len(issueKeys) > 50 {
			return fmt.Errorf("max 50 issues at once, got %d", len(issueKeys))
		}

		if err := client.AddIssuesToEpic(epicKey, issueKeys); err != nil {
			return err
		}

		fmt.Printf("Added %d issues to epic %s\n", len(issueKeys), epicKey)
		return nil
	},
}

var epicRemoveCmd = &cobra.Command{
	Use:   "remove [issue-keys...]",
	Short: "Remove issues from their epic (max 50)",
	Args:  cobra.MinimumNArgs(1),
	Example: `  orbit jira epic remove PROJ-101 PROJ-102`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if len(args) > 50 {
			return fmt.Errorf("max 50 issues at once, got %d", len(args))
		}

		if err := client.RemoveIssuesFromEpic(args); err != nil {
			return err
		}

		fmt.Printf("Removed %d issues from their epic\n", len(args))
		return nil
	},
}

func init() {
	epicCmd.AddCommand(epicListCmd)
	epicCmd.AddCommand(epicCreateCmd)
	epicCmd.AddCommand(epicAddCmd)
	epicCmd.AddCommand(epicRemoveCmd)

	epicListCmd.Flags().StringVar(&epicListOpts.project, "project", "", "filter by project key")
	epicListCmd.Flags().StringSliceVarP(&epicListOpts.status, "status", "s", nil, "filter by status")
	epicListCmd.Flags().IntVar(&epicListOpts.maxResults, "max-results", 50, "max results")

	f := epicCreateCmd.Flags()
	f.StringVar(&epicCreateOpts.project, "project", "", "project key (required)")
	f.StringVarP(&epicCreateOpts.name, "name", "n", "", "epic name (required)")
	f.StringVarP(&epicCreateOpts.summary, "summary", "s", "", "epic summary (defaults to name)")
	f.StringVarP(&epicCreateOpts.description, "body", "b", "", "epic description")
	f.StringVarP(&epicCreateOpts.priority, "priority", "y", "", "priority")
	f.StringSliceVarP(&epicCreateOpts.labels, "label", "l", nil, "labels")
	f.StringSliceVarP(&epicCreateOpts.components, "component", "C", nil, "components")
	_ = epicCreateCmd.MarkFlagRequired("project")
	_ = epicCreateCmd.MarkFlagRequired("name")
}

package jira

import (
	"fmt"
	"strings"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var issueListOpts struct {
	jql           string
	issueType     string
	status        []string
	priority      string
	assignee      string
	reporter      string
	labels        []string
	component     string
	parent        string
	project       string
	created       string
	updated       string
	createdAfter  string
	createdBefore string
	updatedAfter  string
	updatedBefore string
	orderBy       string
	reverse       bool
	startAt       int
	maxResults    int
}

var issueListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "search"},
	Short:   "List issues with filtering",
	Example: `  orbit jira issue list --project PROJ
  orbit jira issue list --assignee me --status "In Progress"
  orbit jira issue list --jql "project = PROJ AND sprint in openSprints()"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		jql := buildJQL()
		result, err := client.SearchIssues(jql, issueListOpts.startAt, issueListOpts.maxResults)
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
					issue.Key,
					issue.Fields.IssueType.Name,
					issue.Fields.Status.Name,
					issue.Fields.Priority.Name,
					assignee,
					issue.Fields.Summary,
				})
			}
			return rows
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Showing %d of %d issues\n", len(result.Issues), result.Total)
		return output.Print(cmdutil.OutputFormat(cmd), result.Issues, headers, rowFn)
	},
}

func init() {
	f := issueListCmd.Flags()
	f.StringVarP(&issueListOpts.jql, "jql", "q", "", "raw JQL query (overrides other filters)")
	f.StringVarP(&issueListOpts.issueType, "type", "t", "", "filter by issue type")
	f.StringSliceVarP(&issueListOpts.status, "status", "s", nil, "filter by status (repeatable)")
	f.StringVarP(&issueListOpts.priority, "priority", "y", "", "filter by priority")
	f.StringVarP(&issueListOpts.assignee, "assignee", "a", "", "filter by assignee")
	f.StringVarP(&issueListOpts.reporter, "reporter", "r", "", "filter by reporter")
	f.StringSliceVarP(&issueListOpts.labels, "label", "l", nil, "filter by label (repeatable)")
	f.StringVarP(&issueListOpts.component, "component", "C", "", "filter by component")
	f.StringVarP(&issueListOpts.parent, "parent", "P", "", "filter by parent issue key")
	f.StringVar(&issueListOpts.project, "project", "", "filter by project key")
	f.StringVar(&issueListOpts.created, "created", "", "created date filter")
	f.StringVar(&issueListOpts.updated, "updated", "", "updated date filter")
	f.StringVar(&issueListOpts.createdAfter, "created-after", "", "created after date")
	f.StringVar(&issueListOpts.createdBefore, "created-before", "", "created before date")
	f.StringVar(&issueListOpts.updatedAfter, "updated-after", "", "updated after date")
	f.StringVar(&issueListOpts.updatedBefore, "updated-before", "", "updated before date")
	f.StringVar(&issueListOpts.orderBy, "order-by", "created", "order by field")
	f.BoolVar(&issueListOpts.reverse, "reverse", false, "reverse display order")
	f.IntVar(&issueListOpts.startAt, "start-at", 0, "pagination start index")
	f.IntVar(&issueListOpts.maxResults, "max-results", 50, "max results to return")
}

func buildJQL() string {
	if issueListOpts.jql != "" {
		return issueListOpts.jql
	}

	var clauses []string
	if issueListOpts.project != "" {
		clauses = append(clauses, fmt.Sprintf("project = %q", issueListOpts.project))
	}
	if issueListOpts.issueType != "" {
		clauses = append(clauses, fmt.Sprintf("issuetype = %q", issueListOpts.issueType))
	}
	if len(issueListOpts.status) > 0 {
		quoted := make([]string, len(issueListOpts.status))
		for i, s := range issueListOpts.status {
			quoted[i] = fmt.Sprintf("%q", s)
		}
		clauses = append(clauses, fmt.Sprintf("status in (%s)", strings.Join(quoted, ", ")))
	}
	if issueListOpts.priority != "" {
		clauses = append(clauses, fmt.Sprintf("priority = %q", issueListOpts.priority))
	}
	if issueListOpts.assignee != "" {
		if issueListOpts.assignee == "me" {
			clauses = append(clauses, "assignee = currentUser()")
		} else {
			clauses = append(clauses, fmt.Sprintf("assignee = %q", issueListOpts.assignee))
		}
	}
	if issueListOpts.reporter != "" {
		clauses = append(clauses, fmt.Sprintf("reporter = %q", issueListOpts.reporter))
	}
	if len(issueListOpts.labels) > 0 {
		for _, l := range issueListOpts.labels {
			clauses = append(clauses, fmt.Sprintf("labels = %q", l))
		}
	}
	if issueListOpts.component != "" {
		clauses = append(clauses, fmt.Sprintf("component = %q", issueListOpts.component))
	}
	if issueListOpts.parent != "" {
		clauses = append(clauses, fmt.Sprintf("parent = %q", issueListOpts.parent))
	}
	if issueListOpts.created != "" {
		clauses = append(clauses, fmt.Sprintf("created >= %q", issueListOpts.created))
	}
	if issueListOpts.updated != "" {
		clauses = append(clauses, fmt.Sprintf("updated >= %q", issueListOpts.updated))
	}
	if issueListOpts.createdAfter != "" {
		clauses = append(clauses, fmt.Sprintf("created >= %q", issueListOpts.createdAfter))
	}
	if issueListOpts.createdBefore != "" {
		clauses = append(clauses, fmt.Sprintf("created <= %q", issueListOpts.createdBefore))
	}
	if issueListOpts.updatedAfter != "" {
		clauses = append(clauses, fmt.Sprintf("updated >= %q", issueListOpts.updatedAfter))
	}
	if issueListOpts.updatedBefore != "" {
		clauses = append(clauses, fmt.Sprintf("updated <= %q", issueListOpts.updatedBefore))
	}

	jql := strings.Join(clauses, " AND ")

	order := "ORDER BY " + issueListOpts.orderBy
	if issueListOpts.reverse {
		order += " ASC"
	} else {
		order += " DESC"
	}
	if jql != "" {
		jql += " " + order
	} else {
		jql = order
	}
	return jql
}

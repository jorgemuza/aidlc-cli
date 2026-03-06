package jira

import (
	"fmt"
	"strings"

	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var issueCreateOpts struct {
	project         string
	issueType       string
	summary         string
	description     string
	priority        string
	assignee        string
	reporter        string
	parent          string
	labels          []string
	components      []string
	fixVersions     []string
	affectsVersions []string
	estimate        string
	epicName        string
	fields          []string
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Example: `  orbit jira issue create --project PROJ --type Story --summary "Add login page"
  orbit jira issue create --project PROJ --type Bug --summary "Fix timeout" --priority High
  orbit jira issue create --project PROJ --type Sub-task --parent PROJ-123 --summary "Add validation"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		req := &jirasvc.CreateIssueRequest{
			Fields: jirasvc.CreateIssueFields{
				Project:   map[string]string{"key": issueCreateOpts.project},
				IssueType: map[string]string{"name": issueCreateOpts.issueType},
				Summary:   issueCreateOpts.summary,
			},
		}

		if issueCreateOpts.description != "" {
			req.Fields.Description = issueCreateOpts.description
		}
		if issueCreateOpts.priority != "" {
			req.Fields.Priority = map[string]string{"name": issueCreateOpts.priority}
		}
		if issueCreateOpts.assignee != "" {
			req.Fields.Assignee = map[string]string{"name": issueCreateOpts.assignee}
		}
		if issueCreateOpts.reporter != "" {
			req.Fields.Reporter = map[string]string{"name": issueCreateOpts.reporter}
		}
		if issueCreateOpts.parent != "" {
			req.Fields.Parent = map[string]string{"key": issueCreateOpts.parent}
		}
		if len(issueCreateOpts.labels) > 0 {
			req.Fields.Labels = issueCreateOpts.labels
		}
		if len(issueCreateOpts.components) > 0 {
			comps := make([]map[string]string, len(issueCreateOpts.components))
			for i, c := range issueCreateOpts.components {
				comps[i] = map[string]string{"name": c}
			}
			req.Fields.Components = comps
		}
		if len(issueCreateOpts.fixVersions) > 0 {
			versions := make([]map[string]string, len(issueCreateOpts.fixVersions))
			for i, v := range issueCreateOpts.fixVersions {
				versions[i] = map[string]string{"name": v}
			}
			req.Fields.FixVersions = versions
		}
		if issueCreateOpts.estimate != "" {
			req.Fields.TimeTracking = &jirasvc.TimeTracking{
				OriginalEstimate: issueCreateOpts.estimate,
			}
		}

		// Custom fields
		customFields := make(map[string]any)

		// Auto-set Epic Name (customfield_11523) from summary if creating an Epic
		if strings.EqualFold(issueCreateOpts.issueType, "Epic") {
			epicName := issueCreateOpts.epicName
			if epicName == "" {
				epicName = issueCreateOpts.summary
			}
			customFields["customfield_11523"] = epicName
		} else if issueCreateOpts.epicName != "" {
			customFields["customfield_11523"] = issueCreateOpts.epicName
		}

		// Parse --field flags (key=value)
		for _, f := range issueCreateOpts.fields {
			parts := strings.SplitN(f, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid field format %q, expected key=value", f)
			}
			customFields[parts[0]] = parts[1]
		}

		if len(customFields) > 0 {
			req.Fields.CustomFields = customFields
		}

		result, err := client.CreateIssue(req)
		if err != nil {
			return err
		}

		fmt.Printf("Created %s\n", result.Key)
		return nil
	},
}

func init() {
	f := issueCreateCmd.Flags()
	f.StringVar(&issueCreateOpts.project, "project", "", "project key (required)")
	f.StringVarP(&issueCreateOpts.issueType, "type", "t", "", "issue type (required)")
	f.StringVarP(&issueCreateOpts.summary, "summary", "s", "", "issue summary (required)")
	f.StringVarP(&issueCreateOpts.description, "body", "b", "", "issue description")
	f.StringVarP(&issueCreateOpts.priority, "priority", "y", "", "priority")
	f.StringVarP(&issueCreateOpts.assignee, "assignee", "a", "", "assignee username")
	f.StringVarP(&issueCreateOpts.reporter, "reporter", "r", "", "reporter username")
	f.StringVarP(&issueCreateOpts.parent, "parent", "P", "", "parent issue key (for subtasks)")
	f.StringSliceVarP(&issueCreateOpts.labels, "label", "l", nil, "labels (repeatable)")
	f.StringSliceVarP(&issueCreateOpts.components, "component", "C", nil, "components (repeatable)")
	f.StringSliceVar(&issueCreateOpts.fixVersions, "fix-version", nil, "fix versions (repeatable)")
	f.StringSliceVar(&issueCreateOpts.affectsVersions, "affects-version", nil, "affects versions (repeatable)")
	f.StringVarP(&issueCreateOpts.estimate, "original-estimate", "e", "", "time estimate (e.g. 2d 3h)")
	f.StringVar(&issueCreateOpts.epicName, "epic-name", "", "epic name (auto-set from summary for Epic type)")
	f.StringSliceVarP(&issueCreateOpts.fields, "field", "F", nil, "custom fields as key=value (repeatable)")
	_ = issueCreateCmd.MarkFlagRequired("project")
	_ = issueCreateCmd.MarkFlagRequired("type")
	_ = issueCreateCmd.MarkFlagRequired("summary")
}

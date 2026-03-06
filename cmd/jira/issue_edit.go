package jira

import (
	"fmt"
	"strings"

	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var issueEditOpts struct {
	summary     string
	description string
	priority    string
	labels      []string
	components  []string
	fixVersions []string
}

var issueEditCmd = &cobra.Command{
	Use:     "edit [issue-key]",
	Aliases: []string{"update"},
	Short:   "Edit an existing issue",
	Args:    cobra.ExactArgs(1),
	Example: `  orbit jira issue edit PROJ-123 --summary "Updated title"
  orbit jira issue edit PROJ-123 --priority Critical
  orbit jira issue edit PROJ-123 --label new-label --label -old-label`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		req := &jirasvc.EditIssueRequest{
			Fields: map[string]any{},
			Update: map[string]any{},
		}

		if issueEditOpts.summary != "" {
			req.Fields["summary"] = issueEditOpts.summary
		}
		if issueEditOpts.description != "" {
			req.Fields["description"] = issueEditOpts.description
		}
		if issueEditOpts.priority != "" {
			req.Fields["priority"] = map[string]string{"name": issueEditOpts.priority}
		}

		if len(issueEditOpts.labels) > 0 {
			var ops []map[string]string
			for _, l := range issueEditOpts.labels {
				if strings.HasPrefix(l, "-") {
					ops = append(ops, map[string]string{"remove": l[1:]})
				} else {
					ops = append(ops, map[string]string{"add": l})
				}
			}
			req.Update["labels"] = ops
		}

		if len(issueEditOpts.components) > 0 {
			var ops []map[string]any
			for _, c := range issueEditOpts.components {
				if strings.HasPrefix(c, "-") {
					ops = append(ops, map[string]any{"remove": map[string]string{"name": c[1:]}})
				} else {
					ops = append(ops, map[string]any{"add": map[string]string{"name": c}})
				}
			}
			req.Update["components"] = ops
		}

		if len(issueEditOpts.fixVersions) > 0 {
			var ops []map[string]any
			for _, v := range issueEditOpts.fixVersions {
				if strings.HasPrefix(v, "-") {
					ops = append(ops, map[string]any{"remove": map[string]string{"name": v[1:]}})
				} else {
					ops = append(ops, map[string]any{"add": map[string]string{"name": v}})
				}
			}
			req.Update["fixVersions"] = ops
		}

		if err := client.EditIssue(args[0], req); err != nil {
			return err
		}

		fmt.Printf("Updated %s\n", args[0])
		return nil
	},
}

func init() {
	f := issueEditCmd.Flags()
	f.StringVarP(&issueEditOpts.summary, "summary", "s", "", "new summary")
	f.StringVarP(&issueEditOpts.description, "body", "b", "", "new description")
	f.StringVarP(&issueEditOpts.priority, "priority", "y", "", "new priority")
	f.StringSliceVarP(&issueEditOpts.labels, "label", "l", nil, "add/remove labels (prefix with - to remove)")
	f.StringSliceVarP(&issueEditOpts.components, "component", "C", nil, "add/remove components (prefix with - to remove)")
	f.StringSliceVar(&issueEditOpts.fixVersions, "fix-version", nil, "add/remove fix versions (prefix with - to remove)")
}

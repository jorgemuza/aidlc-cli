package jira

import (
	"fmt"
	"strings"

	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var issueMoveOpts struct {
	comment    string
	resolution string
}

var issueMoveCmd = &cobra.Command{
	Use:     "move [issue-key] [state]",
	Aliases: []string{"transition"},
	Short:   "Transition an issue to a new workflow state",
	Args:    cobra.ExactArgs(2),
	Example: `  orbit jira issue move PROJ-123 "In Progress"
  orbit jira issue move PROJ-123 Done --comment "Fixed in v2.1"
  orbit jira issue move PROJ-123 Done --resolution Fixed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		key := args[0]
		targetState := args[1]

		transitions, err := client.GetTransitions(key)
		if err != nil {
			return err
		}

		var transitionID string
		for _, t := range transitions {
			if strings.EqualFold(t.Name, targetState) || strings.EqualFold(t.To.Name, targetState) {
				transitionID = t.ID
				break
			}
		}
		if transitionID == "" {
			available := make([]string, len(transitions))
			for i, t := range transitions {
				available[i] = t.Name
			}
			return fmt.Errorf("transition %q not available for %s; available: %s", targetState, key, strings.Join(available, ", "))
		}

		req := &jirasvc.TransitionRequest{
			Transition: map[string]string{"id": transitionID},
		}

		if issueMoveOpts.comment != "" {
			req.Update = map[string]any{
				"comment": []map[string]any{
					{"add": map[string]string{"body": issueMoveOpts.comment}},
				},
			}
		}
		if issueMoveOpts.resolution != "" {
			if req.Fields == nil {
				req.Fields = map[string]any{}
			}
			req.Fields["resolution"] = map[string]string{"name": issueMoveOpts.resolution}
		}

		if err := client.TransitionIssue(key, req); err != nil {
			return err
		}

		fmt.Printf("Moved %s to %s\n", key, targetState)
		return nil
	},
}

func init() {
	issueMoveCmd.Flags().StringVar(&issueMoveOpts.comment, "comment", "", "add comment during transition")
	issueMoveCmd.Flags().StringVarP(&issueMoveOpts.resolution, "resolution", "R", "", "set resolution (e.g. Fixed)")
}

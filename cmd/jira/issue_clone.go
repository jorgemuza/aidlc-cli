package jira

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var issueCloneOpts struct {
	summary string
	replace []string
}

var issueCloneCmd = &cobra.Command{
	Use:   "clone [issue-key]",
	Short: "Clone an issue with optional modifications",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit jira issue clone PROJ-123
  orbit jira issue clone PROJ-123 --summary "Cloned: new title"
  orbit jira issue clone PROJ-123 --replace "v1:v2"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		replaceMap := map[string]string{}
		for _, r := range issueCloneOpts.replace {
			parts := strings.SplitN(r, ":", 2)
			if len(parts) == 2 {
				replaceMap[parts[0]] = parts[1]
			}
		}

		result, err := client.CloneIssue(args[0], issueCloneOpts.summary, replaceMap)
		if err != nil {
			return err
		}

		fmt.Printf("Cloned %s -> %s\n", args[0], result.Key)
		return nil
	},
}

func init() {
	issueCloneCmd.Flags().StringVarP(&issueCloneOpts.summary, "summary", "s", "", "override summary")
	issueCloneCmd.Flags().StringSliceVarP(&issueCloneOpts.replace, "replace", "H", nil, "replace text find:replace")
}

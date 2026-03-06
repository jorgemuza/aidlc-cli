package jira

import (
	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible projects",
	Example: `  orbit jira project list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		projects, err := client.ListProjects()
		if err != nil {
			return err
		}

		headers := []string{"KEY", "NAME", "LEAD"}
		rowFn := func() [][]string {
			var rows [][]string
			for _, p := range projects {
				lead := ""
				if p.Lead != nil {
					lead = p.Lead.DisplayName
				}
				rows = append(rows, []string{p.Key, p.Name, lead})
			}
			return rows
		}
		return output.Print(cmdutil.OutputFormat(cmd), projects, headers, rowFn)
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}

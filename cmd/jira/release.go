package jira

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage project releases/versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var releaseListOpts struct {
	project string
}

var releaseListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List project versions/releases",
	Example: `  orbit jira release list --project PROJ`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		if releaseListOpts.project == "" {
			return fmt.Errorf("--project is required")
		}

		versions, err := client.ListVersions(releaseListOpts.project)
		if err != nil {
			return err
		}

		headers := []string{"ID", "NAME", "RELEASED", "DATE", "DESCRIPTION"}
		rowFn := func() [][]string {
			var rows [][]string
			for _, v := range versions {
				released := ""
				if v.Released {
					released = "yes"
				}
				rows = append(rows, []string{v.ID, v.Name, released, v.ReleaseDate, v.Description})
			}
			return rows
		}
		return output.Print(cmdutil.OutputFormat(cmd), versions, headers, rowFn)
	},
}

func init() {
	releaseCmd.AddCommand(releaseListCmd)
	releaseListCmd.Flags().StringVar(&releaseListOpts.project, "project", "", "project key (required)")
	_ = releaseListCmd.MarkFlagRequired("project")
}

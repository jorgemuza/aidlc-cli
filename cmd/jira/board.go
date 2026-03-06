package jira

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Manage boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var boardListOpts struct {
	project string
}

var boardListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List boards",
	Example: `  orbit jira board list
  orbit jira board list --project PROJ`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		boards, err := client.ListBoards(boardListOpts.project)
		if err != nil {
			return err
		}

		headers := []string{"ID", "NAME", "TYPE"}
		rowFn := func() [][]string {
			var rows [][]string
			for _, b := range boards {
				rows = append(rows, []string{fmt.Sprintf("%d", b.ID), b.Name, b.Type})
			}
			return rows
		}
		return output.Print(cmdutil.OutputFormat(cmd), boards, headers, rowFn)
	},
}

func init() {
	boardCmd.AddCommand(boardListCmd)
	boardListCmd.Flags().StringVar(&boardListOpts.project, "project", "", "filter by project key")
}

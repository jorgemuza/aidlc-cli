package bitbucket

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [subcommand]",
	Short: "Manage projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var projectViewCmd = &cobra.Command{
	Use:   "view [project-key]",
	Short: "View a project",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc bitbucket project view L3SUP
  aidlc bb project view MYPROJ`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		project, err := client.GetProject(args[0])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(project, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Key:         %s\n", project.Key)
		fmt.Printf("Name:        %s\n", project.Name)
		fmt.Printf("Description: %s\n", project.Description)
		fmt.Printf("Public:      %v\n", project.Public)
		fmt.Printf("Type:        %s\n", project.Type)
		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List projects",
	Example: `  aidlc bitbucket project list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		projects, err := client.ListProjects(limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(projects, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-12s %-30s %s\n", "KEY", "NAME", "TYPE")
		fmt.Printf("%-12s %-30s %s\n", "---", "----", "----")
		for _, p := range projects {
			fmt.Printf("%-12s %-30s %s\n", p.Key, p.Name, p.Type)
		}
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectViewCmd)
	projectCmd.AddCommand(projectListCmd)
	projectListCmd.Flags().Int("limit", 50, "max results")
}

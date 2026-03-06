package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [id-or-path]",
	Short: "View a project",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab project 595
  orbit gitlab project schools/frontend/schools-frontend-backoffice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
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

		fmt.Printf("ID:          %d\n", project.ID)
		fmt.Printf("Name:        %s\n", project.Name)
		fmt.Printf("Path:        %s\n", project.PathWithNamespace)
		fmt.Printf("Description: %s\n", project.Description)
		fmt.Printf("Default:     %s\n", project.DefaultBranch)
		fmt.Printf("Visibility:  %s\n", project.Visibility)
		fmt.Printf("URL:         %s\n", project.WebURL)
		fmt.Printf("Activity:    %s\n", project.LastActivityAt)
		return nil
	},
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List projects",
	Example: `  orbit gitlab projects
  orbit gitlab projects --search frontend
  orbit gitlab projects --group schools/frontend`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		search, _ := cmd.Flags().GetString("search")
		group, _ := cmd.Flags().GetString("group")
		limit, _ := cmd.Flags().GetInt("limit")

		format, _ := cmd.Flags().GetString("output")

		if group != "" {
			projects, err := client.ListGroupProjects(group, limit)
			if err != nil {
				return err
			}

			if format == "json" {
				data, _ := json.MarshalIndent(projects, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			fmt.Printf("%-8s %-50s %s\n", "ID", "PATH", "LAST ACTIVITY")
			fmt.Printf("%-8s %-50s %s\n", "--", "----", "-------------")
			for _, p := range projects {
				fmt.Printf("%-8d %-50s %s\n", p.ID, p.PathWithNamespace, p.LastActivityAt[:10])
			}
			return nil
		}

		projects, err := client.ListProjects(search, limit)
		if err != nil {
			return err
		}

		if format == "json" {
			data, _ := json.MarshalIndent(projects, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-8s %-50s %s\n", "ID", "PATH", "LAST ACTIVITY")
		fmt.Printf("%-8s %-50s %s\n", "--", "----", "-------------")
		for _, p := range projects {
			activity := ""
			if len(p.LastActivityAt) >= 10 {
				activity = p.LastActivityAt[:10]
			}
			fmt.Printf("%-8d %-50s %s\n", p.ID, p.PathWithNamespace, activity)
		}
		return nil
	},
}

func init() {
	projectsCmd.Flags().String("search", "", "search projects by name")
	projectsCmd.Flags().String("group", "", "list projects in a group (id or path)")
	projectsCmd.Flags().Int("limit", 50, "max results")
}

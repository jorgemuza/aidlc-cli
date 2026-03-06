package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group [subcommand]",
	Short: "Manage groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var groupViewCmd = &cobra.Command{
	Use:   "view [id-or-path]",
	Short: "View a group",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab group view schools/frontend`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		group, err := client.GetGroup(args[0])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(group, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:          %d\n", group.ID)
		fmt.Printf("Name:        %s\n", group.Name)
		fmt.Printf("Path:        %s\n", group.FullPath)
		fmt.Printf("Description: %s\n", group.Description)
		fmt.Printf("Visibility:  %s\n", group.Visibility)
		fmt.Printf("URL:         %s\n", group.WebURL)
		return nil
	},
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups",
	Example: `  orbit gitlab group list
  orbit gitlab group list --search schools`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		search, _ := cmd.Flags().GetString("search")
		limit, _ := cmd.Flags().GetInt("limit")

		groups, err := client.ListGroups(search, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(groups, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-8s %-40s %s\n", "ID", "PATH", "VISIBILITY")
		fmt.Printf("%-8s %-40s %s\n", "--", "----", "----------")
		for _, g := range groups {
			fmt.Printf("%-8d %-40s %s\n", g.ID, g.FullPath, g.Visibility)
		}
		return nil
	},
}

var groupSubgroupsCmd = &cobra.Command{
	Use:   "subgroups [id-or-path]",
	Short: "List subgroups of a group",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab group subgroups schools`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		groups, err := client.ListSubgroups(args[0], limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(groups, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-8s %-40s %s\n", "ID", "PATH", "VISIBILITY")
		fmt.Printf("%-8s %-40s %s\n", "--", "----", "----------")
		for _, g := range groups {
			fmt.Printf("%-8d %-40s %s\n", g.ID, g.FullPath, g.Visibility)
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupViewCmd)
	groupCmd.AddCommand(groupListCmd)
	groupCmd.AddCommand(groupSubgroupsCmd)

	groupListCmd.Flags().String("search", "", "search groups by name")
	groupListCmd.Flags().Int("limit", 50, "max results")
	groupSubgroupsCmd.Flags().Int("limit", 50, "max results")
}

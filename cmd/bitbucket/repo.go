package bitbucket

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo [subcommand]",
	Short: "Manage repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var repoViewCmd = &cobra.Command{
	Use:   "view [project-key] [repo-slug]",
	Short: "View a repository",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit bitbucket repo view L3SUP agents-sre
  orbit bb repo view MYPROJ my-repo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		repo, err := client.GetRepository(args[0], args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(repo, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:          %d\n", repo.ID)
		fmt.Printf("Slug:        %s\n", repo.Slug)
		fmt.Printf("Name:        %s\n", repo.Name)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("State:       %s\n", repo.State)
		fmt.Printf("SCM:         %s\n", repo.ScmID)
		fmt.Printf("Forkable:    %v\n", repo.Forkable)
		if repo.Project != nil {
			fmt.Printf("Project:     %s (%s)\n", repo.Project.Name, repo.Project.Key)
		}
		if len(repo.Links.Clone) > 0 {
			for _, l := range repo.Links.Clone {
				fmt.Printf("Clone (%s): %s\n", l.Name, l.Href)
			}
		}
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list [project-key]",
	Short: "List repositories in a project",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit bitbucket repo list L3SUP
  orbit bb repo list MYPROJ`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		repos, err := client.ListRepositories(args[0], limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(repos, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-8s %-40s %s\n", "ID", "SLUG", "STATE")
		fmt.Printf("%-8s %-40s %s\n", "--", "----", "-----")
		for _, r := range repos {
			fmt.Printf("%-8d %-40s %s\n", r.ID, r.Slug, r.State)
		}
		return nil
	},
}

func init() {
	repoCmd.AddCommand(repoViewCmd)
	repoCmd.AddCommand(repoListCmd)
	repoListCmd.Flags().Int("limit", 50, "max results")
}

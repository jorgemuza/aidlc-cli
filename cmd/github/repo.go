package github

import (
	"encoding/json"
	"fmt"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo [owner/repo]",
	Short: "View a repository",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github repo octocat/hello-world
  orbit gh repo kubernetes/kubernetes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		r, err := client.GetRepo(owner, repo)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Name:        %s\n", r.FullName)
		fmt.Printf("Description: %s\n", r.Description)
		fmt.Printf("Language:    %s\n", r.Language)
		fmt.Printf("Default:     %s\n", r.DefaultBranch)
		fmt.Printf("Private:     %v\n", r.Private)
		fmt.Printf("Stars:       %d\n", r.StargazersCount)
		fmt.Printf("Forks:       %d\n", r.ForksCount)
		fmt.Printf("Issues:      %d\n", r.OpenIssuesCount)
		fmt.Printf("URL:         %s\n", r.HTMLURL)
		fmt.Printf("Pushed:      %s\n", r.PushedAt)
		return nil
	},
}

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List repositories",
	Example: `  orbit github repos
  orbit github repos --org kubernetes
  orbit gh repos --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		org, _ := cmd.Flags().GetString("org")
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("output")

		if org != "" {
			repos, err := client.ListOrgRepos(org, limit)
			if err != nil {
				return err
			}

			if format == "json" {
				data, _ := json.MarshalIndent(repos, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			fmt.Printf("%-50s %-10s %s\n", "REPO", "LANGUAGE", "PUSHED")
			fmt.Printf("%-50s %-10s %s\n", "----", "--------", "------")
			for _, r := range repos {
				pushed := ""
				if len(r.PushedAt) >= 10 {
					pushed = r.PushedAt[:10]
				}
				fmt.Printf("%-50s %-10s %s\n", r.FullName, r.Language, pushed)
			}
			return nil
		}

		repos, err := client.ListUserRepos(limit)
		if err != nil {
			return err
		}

		if format == "json" {
			data, _ := json.MarshalIndent(repos, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-50s %-10s %s\n", "REPO", "LANGUAGE", "PUSHED")
		fmt.Printf("%-50s %-10s %s\n", "----", "--------", "------")
		for _, r := range repos {
			pushed := ""
			if len(r.PushedAt) >= 10 {
				pushed = r.PushedAt[:10]
			}
			fmt.Printf("%-50s %-10s %s\n", r.FullName, r.Language, pushed)
		}
		return nil
	},
}

func init() {
	reposCmd.Flags().String("org", "", "list repos for an organization")
	reposCmd.Flags().Int("limit", 30, "max results")
}

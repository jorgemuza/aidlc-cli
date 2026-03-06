package github

import (
	"encoding/json"
	"fmt"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [subcommand]",
	Short: "Manage branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var branchListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List branches",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github branch list octocat/hello-world
  orbit gh branch list kubernetes/kubernetes --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		branches, err := client.ListBranches(owner, repo, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branches, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-40s %-9s %s\n", "NAME", "PROTECTED", "SHA")
		fmt.Printf("%-40s %-9s %s\n", "----", "---------", "---")
		for _, b := range branches {
			sha := ""
			if b.Commit != nil {
				sha = b.Commit.SHA[:7]
			}
			fmt.Printf("%-40s %-9v %s\n", b.Name, b.Protected, sha)
		}
		return nil
	},
}

var branchViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [branch]",
	Short: "View a branch",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github branch view octocat/hello-world main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		branch, err := client.GetBranch(owner, repo, args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branch, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Name:      %s\n", branch.Name)
		fmt.Printf("Protected: %v\n", branch.Protected)
		if branch.Commit != nil {
			fmt.Printf("SHA:       %s\n", branch.Commit.SHA)
			fmt.Printf("URL:       %s\n", branch.Commit.HTMLURL)
			if branch.Commit.Commit != nil {
				if branch.Commit.Commit.Author != nil {
					fmt.Printf("Author:    %s\n", branch.Commit.Commit.Author.Name)
					fmt.Printf("Date:      %s\n", branch.Commit.Commit.Author.Date)
				}
				fmt.Printf("Message:   %s\n", branch.Commit.Commit.Message)
			}
		}
		return nil
	},
}

func init() {
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchViewCmd)

	branchListCmd.Flags().Int("limit", 50, "max results")
}

package github

import (
	"encoding/json"
	"fmt"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag [subcommand]",
	Short: "Manage tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List tags",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github tag list octocat/hello-world`,
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
		tags, err := client.ListTags(owner, repo, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(tags, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-30s %s\n", "TAG", "COMMIT")
		fmt.Printf("%-30s %s\n", "---", "------")
		for _, t := range tags {
			sha := ""
			if t.Commit != nil {
				sha = t.Commit.SHA[:7]
			}
			fmt.Printf("%-30s %s\n", t.Name, sha)
		}
		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)

	tagListCmd.Flags().Int("limit", 50, "max results")
}

package github

import (
	"encoding/json"
	"fmt"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit [subcommand]",
	Short: "View commits",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var commitListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List commits",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github commit list octocat/hello-world
  orbit gh commit list octocat/hello-world --ref main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		ref, _ := cmd.Flags().GetString("ref")
		limit, _ := cmd.Flags().GetInt("limit")
		commits, err := client.ListCommits(owner, repo, ref, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(commits, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-12s %-20s %s\n", "SHA", "DATE", "AUTHOR", "MESSAGE")
		fmt.Printf("%-10s %-12s %-20s %s\n", "---", "----", "------", "-------")
		for _, c := range commits {
			sha := c.SHA[:7]
			date := ""
			author := ""
			message := ""
			if c.Commit != nil {
				if c.Commit.Author != nil {
					author = c.Commit.Author.Name
					if len(c.Commit.Author.Date) >= 10 {
						date = c.Commit.Author.Date[:10]
					}
				}
				message = c.Commit.Message
				if idx := len(message); idx > 60 {
					message = message[:57] + "..."
				}
				// Take only first line
				for i, ch := range message {
					if ch == '\n' {
						message = message[:i]
						break
					}
				}
			}
			fmt.Printf("%-10s %-12s %-20s %s\n", sha, date, author, message)
		}
		return nil
	},
}

var commitViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [sha]",
	Short: "View a commit",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github commit view octocat/hello-world abc1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		c, err := client.GetCommit(owner, repo, args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(c, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("SHA:     %s\n", c.SHA)
		if c.Commit != nil {
			if c.Commit.Author != nil {
				fmt.Printf("Author:  %s <%s>\n", c.Commit.Author.Name, c.Commit.Author.Email)
				fmt.Printf("Date:    %s\n", c.Commit.Author.Date)
			}
			fmt.Printf("Message: %s\n", c.Commit.Message)
		}
		fmt.Printf("URL:     %s\n", c.HTMLURL)
		return nil
	},
}

func init() {
	commitCmd.AddCommand(commitListCmd)
	commitCmd.AddCommand(commitViewCmd)

	commitListCmd.Flags().String("ref", "", "branch or tag name")
	commitListCmd.Flags().Int("limit", 20, "max results")
}

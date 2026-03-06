package bitbucket

import (
	"encoding/json"
	"fmt"
	"time"

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
	Use:   "list [project-key] [repo-slug]",
	Short: "List recent commits",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc bb commit list L3SUP agents-sre
  aidlc bb commit list L3SUP agents-sre --branch main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		branch, _ := cmd.Flags().GetString("branch")
		limit, _ := cmd.Flags().GetInt("limit")
		commits, err := client.ListCommits(args[0], args[1], branch, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(commits, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-12s %-20s %-12s %s\n", "ID", "AUTHOR", "DATE", "MESSAGE")
		fmt.Printf("%-12s %-20s %-12s %s\n", "--", "------", "----", "-------")
		for _, c := range commits {
			id := c.DisplayID
			if id == "" && len(c.ID) > 12 {
				id = c.ID[:12]
			}
			author := ""
			if c.Author != nil {
				author = c.Author.Name
			}
			date := ""
			if c.AuthorTS > 0 {
				date = time.UnixMilli(c.AuthorTS).Format("2006-01-02")
			}
			msg := c.Message
			if len(msg) > 60 {
				msg = msg[:57] + "..."
			}
			fmt.Printf("%-12s %-20s %-12s %s\n", id, truncate(author, 20), date, firstLine(msg))
		}
		return nil
	},
}

var commitViewCmd = &cobra.Command{
	Use:   "view [project-key] [repo-slug] [commit-id]",
	Short: "View a commit",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb commit view L3SUP agents-sre abc1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		commit, err := client.GetCommit(args[0], args[1], args[2])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(commit, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:      %s\n", commit.ID)
		if commit.Author != nil {
			fmt.Printf("Author:  %s <%s>\n", commit.Author.Name, commit.Author.EmailAddress)
		}
		if commit.AuthorTS > 0 {
			fmt.Printf("Date:    %s\n", time.UnixMilli(commit.AuthorTS).Format(time.RFC3339))
		}
		fmt.Printf("Message: %s\n", commit.Message)
		return nil
	},
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func firstLine(s string) string {
	for i, c := range s {
		if c == '\n' || c == '\r' {
			return s[:i]
		}
	}
	return s
}

func init() {
	commitCmd.AddCommand(commitListCmd)
	commitCmd.AddCommand(commitViewCmd)

	commitListCmd.Flags().String("branch", "", "branch to list commits from")
	commitListCmd.Flags().Int("limit", 20, "max results")
}

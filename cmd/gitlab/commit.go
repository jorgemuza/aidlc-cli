package gitlab

import (
	"encoding/json"
	"fmt"

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
	Use:   "list [project]",
	Short: "List commits",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab commit list 595
  orbit gitlab commit list 595 --ref main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		ref, _ := cmd.Flags().GetString("ref")
		limit, _ := cmd.Flags().GetInt("limit")
		commits, err := client.ListCommits(args[0], ref, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(commits, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-12s %-20s %s\n", "SHA", "DATE", "AUTHOR", "TITLE")
		fmt.Printf("%-10s %-12s %-20s %s\n", "---", "----", "------", "-----")
		for _, c := range commits {
			date := ""
			if len(c.CommittedDate) >= 10 {
				date = c.CommittedDate[:10]
			}
			title := c.Title
			if len(title) > 60 {
				title = title[:57] + "..."
			}
			fmt.Printf("%-10s %-12s %-20s %s\n", c.ShortID, date, c.AuthorName, title)
		}
		return nil
	},
}

var commitViewCmd = &cobra.Command{
	Use:   "view [project] [sha]",
	Short: "View a commit",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab commit view 595 abc1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		c, err := client.GetCommit(args[0], args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(c, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("SHA:       %s\n", c.ID)
		fmt.Printf("Author:    %s <%s>\n", c.AuthorName, c.AuthorEmail)
		fmt.Printf("Date:      %s\n", c.CommittedDate)
		fmt.Printf("Title:     %s\n", c.Title)
		if c.Message != c.Title {
			fmt.Printf("Message:\n%s\n", c.Message)
		}
		fmt.Printf("URL:       %s\n", c.WebURL)
		return nil
	},
}

func init() {
	commitCmd.AddCommand(commitListCmd)
	commitCmd.AddCommand(commitViewCmd)

	commitListCmd.Flags().String("ref", "", "branch or tag name")
	commitListCmd.Flags().Int("limit", 20, "max results")
}

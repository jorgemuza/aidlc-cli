package github

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue [subcommand]",
	Short: "Manage issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var issueListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List issues",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github issue list octocat/hello-world
  orbit gh issue list octocat/hello-world --state closed --labels bug`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		state, _ := cmd.Flags().GetString("state")
		labelsStr, _ := cmd.Flags().GetString("labels")
		limit, _ := cmd.Flags().GetInt("limit")

		var labels []string
		if labelsStr != "" {
			labels = strings.Split(labelsStr, ",")
		}

		issues, err := client.ListIssues(owner, repo, state, labels, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(issues, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-6s %-8s %-50s %s\n", "#", "STATE", "TITLE", "LABELS")
		fmt.Printf("%-6s %-8s %-50s %s\n", "-", "-----", "-----", "------")
		for _, i := range issues {
			title := i.Title
			if len(title) > 48 {
				title = title[:45] + "..."
			}
			labelNames := ""
			if len(i.Labels) > 0 {
				names := make([]string, len(i.Labels))
				for j, l := range i.Labels {
					names[j] = l.Name
				}
				labelNames = strings.Join(names, ", ")
			}
			fmt.Printf("#%-5d %-8s %-50s %s\n", i.Number, i.State, title, labelNames)
		}
		return nil
	},
}

var issueViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [number]",
	Short: "View an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github issue view octocat/hello-world 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		number, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid issue number: %s", args[1])
		}

		issue, err := client.GetIssue(owner, repo, number)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(issue, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Issue:  #%d\n", issue.Number)
		fmt.Printf("Title:  %s\n", issue.Title)
		fmt.Printf("State:  %s\n", issue.State)
		if issue.User != nil {
			fmt.Printf("Author: %s\n", issue.User.Login)
		}
		if len(issue.Assignees) > 0 {
			names := make([]string, len(issue.Assignees))
			for i, a := range issue.Assignees {
				names[i] = a.Login
			}
			fmt.Printf("Assign: %s\n", strings.Join(names, ", "))
		}
		if len(issue.Labels) > 0 {
			names := make([]string, len(issue.Labels))
			for i, l := range issue.Labels {
				names[i] = l.Name
			}
			fmt.Printf("Labels: %s\n", strings.Join(names, ", "))
		}
		if issue.Milestone != nil {
			fmt.Printf("Mile:   %s\n", issue.Milestone.Title)
		}
		fmt.Printf("URL:    %s\n", issue.HTMLURL)
		if issue.Body != "" {
			fmt.Printf("\n%s\n", issue.Body)
		}
		return nil
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create [owner/repo]",
	Short: "Create an issue",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github issue create octocat/hello-world --title "Fix login bug" --labels bug,urgent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		labelsStr, _ := cmd.Flags().GetString("labels")

		var labels []string
		if labelsStr != "" {
			labels = strings.Split(labelsStr, ",")
		}

		issue, err := client.CreateIssue(owner, repo, title, body, labels, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Created issue #%d: %s\n", issue.Number, issue.Title)
		fmt.Printf("URL: %s\n", issue.HTMLURL)
		return nil
	},
}

var issueCloseCmd = &cobra.Command{
	Use:   "close [owner/repo] [number]",
	Short: "Close an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github issue close octocat/hello-world 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		number, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid issue number: %s", args[1])
		}

		issue, err := client.UpdateIssue(owner, repo, number, map[string]any{"state": "closed"})
		if err != nil {
			return err
		}

		fmt.Printf("Closed issue #%d: %s\n", issue.Number, issue.Title)
		return nil
	},
}

var issueCommentCmd = &cobra.Command{
	Use:   "comment [owner/repo] [number]",
	Short: "Add a comment to an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github issue comment octocat/hello-world 1 --body "Working on this"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		number, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid issue number: %s", args[1])
		}

		body, _ := cmd.Flags().GetString("body")
		comment, err := client.CreateIssueComment(owner, repo, number, body)
		if err != nil {
			return err
		}

		fmt.Printf("Added comment #%d to issue #%d\n", comment.ID, number)
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueViewCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueCloseCmd)
	issueCmd.AddCommand(issueCommentCmd)

	issueListCmd.Flags().String("state", "", "filter by state: open, closed, all")
	issueListCmd.Flags().String("labels", "", "filter by labels (comma-separated)")
	issueListCmd.Flags().Int("limit", 20, "max results")

	issueCreateCmd.Flags().String("title", "", "issue title (required)")
	issueCreateCmd.Flags().String("body", "", "issue body")
	issueCreateCmd.Flags().String("labels", "", "labels (comma-separated)")
	_ = issueCreateCmd.MarkFlagRequired("title")

	issueCommentCmd.Flags().String("body", "", "comment body (required)")
	_ = issueCommentCmd.MarkFlagRequired("body")
}

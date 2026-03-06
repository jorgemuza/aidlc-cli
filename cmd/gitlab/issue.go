package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	Use:   "list [project]",
	Short: "List issues",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab issue list 595
  aidlc gitlab issue list 595 --state opened --labels bug`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
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

		issues, err := client.ListIssues(args[0], state, labels, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(issues, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-6s %-8s %-50s %s\n", "IID", "STATE", "TITLE", "LABELS")
		fmt.Printf("%-6s %-8s %-50s %s\n", "---", "-----", "-----", "------")
		for _, i := range issues {
			title := i.Title
			if len(title) > 48 {
				title = title[:45] + "..."
			}
			labels := ""
			if len(i.Labels) > 0 {
				labels = strings.Join(i.Labels, ", ")
			}
			fmt.Printf("#%-5d %-8s %-50s %s\n", i.IID, i.State, title, labels)
		}
		return nil
	},
}

var issueViewCmd = &cobra.Command{
	Use:   "view [project] [issue-iid]",
	Short: "View an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab issue view 595 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid issue IID: %s", args[1])
		}

		issue, err := client.GetIssue(args[0], iid)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(issue, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Issue:    #%d\n", issue.IID)
		fmt.Printf("Title:    %s\n", issue.Title)
		fmt.Printf("State:    %s\n", issue.State)
		if issue.Author != nil {
			fmt.Printf("Author:   %s\n", issue.Author.Username)
		}
		if len(issue.Assignees) > 0 {
			names := make([]string, len(issue.Assignees))
			for i, a := range issue.Assignees {
				names[i] = a.Username
			}
			fmt.Printf("Assignee: %s\n", strings.Join(names, ", "))
		}
		if len(issue.Labels) > 0 {
			fmt.Printf("Labels:   %s\n", strings.Join(issue.Labels, ", "))
		}
		if issue.DueDate != "" {
			fmt.Printf("Due:      %s\n", issue.DueDate)
		}
		if issue.Milestone != nil {
			fmt.Printf("Mile:     %s\n", issue.Milestone.Title)
		}
		fmt.Printf("URL:      %s\n", issue.WebURL)
		if issue.Description != "" {
			fmt.Printf("\n%s\n", issue.Description)
		}
		return nil
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "Create an issue",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab issue create 595 --title "Fix login bug" --labels bug,urgent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("description")
		labelsStr, _ := cmd.Flags().GetString("labels")

		var labels []string
		if labelsStr != "" {
			labels = strings.Split(labelsStr, ",")
		}

		issue, err := client.CreateIssue(args[0], title, desc, labels, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Created issue #%d: %s\n", issue.IID, issue.Title)
		fmt.Printf("URL: %s\n", issue.WebURL)
		return nil
	},
}

var issueCloseCmd = &cobra.Command{
	Use:   "close [project] [issue-iid]",
	Short: "Close an issue",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab issue close 595 1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid issue IID: %s", args[1])
		}

		issue, err := client.UpdateIssue(args[0], iid, map[string]any{"state_event": "close"})
		if err != nil {
			return err
		}

		fmt.Printf("Closed issue #%d: %s\n", issue.IID, issue.Title)
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueViewCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueCloseCmd)

	issueListCmd.Flags().String("state", "", "filter by state: opened, closed, all")
	issueListCmd.Flags().String("labels", "", "filter by labels (comma-separated)")
	issueListCmd.Flags().Int("limit", 20, "max results")

	issueCreateCmd.Flags().String("title", "", "issue title (required)")
	issueCreateCmd.Flags().String("description", "", "issue description")
	issueCreateCmd.Flags().String("labels", "", "labels (comma-separated)")
	_ = issueCreateCmd.MarkFlagRequired("title")
}

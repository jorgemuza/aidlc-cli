package jira

import (
	"fmt"
	"strings"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var issueViewOpts struct {
	comments int
}

var issueViewCmd = &cobra.Command{
	Use:     "view [issue-key]",
	Aliases: []string{"show"},
	Short:   "View issue details",
	Args:    cobra.ExactArgs(1),
	Example: `  orbit jira issue view PROJ-123
  orbit jira issue view PROJ-123 --comments 5
  orbit jira issue view PROJ-123 -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		issue, err := client.GetIssue(args[0], issueViewOpts.comments)
		if err != nil {
			return err
		}

		format := cmdutil.OutputFormat(cmd)
		if format != output.FormatTable {
			return output.Print(format, issue, nil, nil)
		}

		printIssueDetail(issue)
		return nil
	},
}

func init() {
	issueViewCmd.Flags().IntVar(&issueViewOpts.comments, "comments", 1, "number of recent comments to show")
}

func printIssueDetail(issue *jirasvc.Issue) {
	f := issue.Fields
	fmt.Printf("%-14s %s\n", "Key:", issue.Key)
	fmt.Printf("%-14s %s\n", "Summary:", f.Summary)
	fmt.Printf("%-14s %s\n", "Type:", f.IssueType.Name)
	fmt.Printf("%-14s %s\n", "Status:", f.Status.Name)
	fmt.Printf("%-14s %s\n", "Priority:", f.Priority.Name)
	if f.Resolution != nil {
		fmt.Printf("%-14s %s\n", "Resolution:", f.Resolution.Name)
	}
	if f.Assignee != nil {
		fmt.Printf("%-14s %s\n", "Assignee:", f.Assignee.DisplayName)
	}
	if f.Reporter != nil {
		fmt.Printf("%-14s %s\n", "Reporter:", f.Reporter.DisplayName)
	}
	if f.Parent != nil {
		fmt.Printf("%-14s %s\n", "Parent:", f.Parent.Key)
	}
	if len(f.Labels) > 0 {
		fmt.Printf("%-14s %s\n", "Labels:", strings.Join(f.Labels, ", "))
	}
	if len(f.Components) > 0 {
		names := make([]string, len(f.Components))
		for i, c := range f.Components {
			names[i] = c.Name
		}
		fmt.Printf("%-14s %s\n", "Components:", strings.Join(names, ", "))
	}
	if len(f.FixVersions) > 0 {
		names := make([]string, len(f.FixVersions))
		for i, v := range f.FixVersions {
			names[i] = v.Name
		}
		fmt.Printf("%-14s %s\n", "Fix Versions:", strings.Join(names, ", "))
	}
	if f.TimeTracking != nil {
		if f.TimeTracking.OriginalEstimate != "" {
			fmt.Printf("%-14s %s\n", "Estimate:", f.TimeTracking.OriginalEstimate)
		}
		if f.TimeTracking.TimeSpent != "" {
			fmt.Printf("%-14s %s\n", "Time Spent:", f.TimeTracking.TimeSpent)
		}
	}
	fmt.Printf("%-14s %s\n", "Created:", f.Created)
	fmt.Printf("%-14s %s\n", "Updated:", f.Updated)

	if f.Description != "" {
		fmt.Println("\n--- Description ---")
		fmt.Println(f.Description)
	}

	if len(f.Subtasks) > 0 {
		fmt.Println("\n--- Subtasks ---")
		for _, st := range f.Subtasks {
			fmt.Printf("  %s  %-12s  %s\n", st.Key, st.Fields.Status.Name, st.Fields.Summary)
		}
	}

	if len(f.IssueLinks) > 0 {
		fmt.Println("\n--- Links ---")
		for _, link := range f.IssueLinks {
			if link.OutwardIssue != nil {
				fmt.Printf("  %s %s (%s)\n", link.Type.Outward, link.OutwardIssue.Key, link.OutwardIssue.Fields.Summary)
			}
			if link.InwardIssue != nil {
				fmt.Printf("  %s %s (%s)\n", link.Type.Inward, link.InwardIssue.Key, link.InwardIssue.Fields.Summary)
			}
		}
	}

	if f.Comment != nil && len(f.Comment.Comments) > 0 {
		fmt.Printf("\n--- Comments (%d total) ---\n", f.Comment.Total)
		for _, c := range f.Comment.Comments {
			author := "Unknown"
			if c.Author != nil {
				author = c.Author.DisplayName
			}
			fmt.Printf("\n  %s (%s):\n  %s\n", author, c.Created, c.Body)
		}
	}
}

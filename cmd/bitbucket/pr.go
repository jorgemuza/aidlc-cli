package bitbucket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:     "pr [subcommand]",
	Short:   "Manage pull requests",
	Aliases: []string{"pull-request"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var prListCmd = &cobra.Command{
	Use:   "list [project-key] [repo-slug]",
	Short: "List pull requests",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc bb pr list L3SUP agents-sre
  aidlc bb pr list L3SUP agents-sre --state merged`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		prs, err := client.ListPullRequests(args[0], args[1], state, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(prs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-6s %-10s %-20s %s\n", "ID", "STATE", "AUTHOR", "TITLE")
		fmt.Printf("%-6s %-10s %-20s %s\n", "--", "-----", "------", "-----")
		for _, pr := range prs {
			author := ""
			if pr.Author != nil && pr.Author.User != nil {
				author = pr.Author.User.DisplayName
			}
			title := pr.Title
			if len(title) > 60 {
				title = title[:57] + "..."
			}
			fmt.Printf("%-6d %-10s %-20s %s\n", pr.ID, pr.State, truncate(author, 20), title)
		}
		return nil
	},
}

var prViewCmd = &cobra.Command{
	Use:   "view [project-key] [repo-slug] [pr-id]",
	Short: "View a pull request",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb pr view L3SUP agents-sre 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		prID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid PR ID: %s", args[2])
		}

		pr, err := client.GetPullRequest(args[0], args[1], prID)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(pr, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:          %d\n", pr.ID)
		fmt.Printf("Title:       %s\n", pr.Title)
		fmt.Printf("State:       %s\n", pr.State)
		fmt.Printf("From:        %s\n", pr.FromRef.DisplayID)
		fmt.Printf("To:          %s\n", pr.ToRef.DisplayID)
		if pr.Author != nil && pr.Author.User != nil {
			fmt.Printf("Author:      %s\n", pr.Author.User.DisplayName)
		}
		if len(pr.Reviewers) > 0 {
			names := make([]string, 0, len(pr.Reviewers))
			for _, r := range pr.Reviewers {
				status := ""
				if r.Approved {
					status = " (approved)"
				} else if r.Status == "NEEDS_WORK" {
					status = " (needs work)"
				}
				if r.User != nil {
					names = append(names, r.User.DisplayName+status)
				}
			}
			fmt.Printf("Reviewers:   %s\n", strings.Join(names, ", "))
		}
		if pr.CreatedDate > 0 {
			fmt.Printf("Created:     %s\n", time.UnixMilli(pr.CreatedDate).Format(time.RFC3339))
		}
		if pr.UpdatedDate > 0 {
			fmt.Printf("Updated:     %s\n", time.UnixMilli(pr.UpdatedDate).Format(time.RFC3339))
		}
		if len(pr.Links.Self) > 0 {
			fmt.Printf("URL:         %s\n", pr.Links.Self[0].Href)
		}
		if pr.Description != "" {
			fmt.Printf("\n%s\n", pr.Description)
		}
		return nil
	},
}

var prCreateCmd = &cobra.Command{
	Use:   "create [project-key] [repo-slug]",
	Short: "Create a pull request",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc bb pr create L3SUP agents-sre \
    --from feature/new --to main --title "Add new feature"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("description")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		reviewersStr, _ := cmd.Flags().GetString("reviewers")

		var reviewers []string
		if reviewersStr != "" {
			reviewers = strings.Split(reviewersStr, ",")
		}

		pr, err := client.CreatePullRequest(args[0], args[1], title, desc, from, to, reviewers)
		if err != nil {
			return err
		}

		fmt.Printf("Created PR #%d: %s\n", pr.ID, pr.Title)
		if len(pr.Links.Self) > 0 {
			fmt.Printf("URL: %s\n", pr.Links.Self[0].Href)
		}
		return nil
	},
}

var prMergeCmd = &cobra.Command{
	Use:   "merge [project-key] [repo-slug] [pr-id]",
	Short: "Merge a pull request",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb pr merge L3SUP agents-sre 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		prID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid PR ID: %s", args[2])
		}

		// Fetch current version for optimistic locking
		pr, err := client.GetPullRequest(args[0], args[1], prID)
		if err != nil {
			return err
		}

		merged, err := client.MergePullRequest(args[0], args[1], prID, int(pr.UpdatedDate))
		if err != nil {
			return err
		}

		fmt.Printf("Merged PR #%d: %s\n", merged.ID, merged.Title)
		return nil
	},
}

var prDeclineCmd = &cobra.Command{
	Use:   "decline [project-key] [repo-slug] [pr-id]",
	Short: "Decline a pull request",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb pr decline L3SUP agents-sre 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		prID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid PR ID: %s", args[2])
		}

		pr, err := client.GetPullRequest(args[0], args[1], prID)
		if err != nil {
			return err
		}

		declined, err := client.DeclinePullRequest(args[0], args[1], prID, int(pr.UpdatedDate))
		if err != nil {
			return err
		}

		fmt.Printf("Declined PR #%d: %s\n", declined.ID, declined.Title)
		return nil
	},
}

var prCommentCmd = &cobra.Command{
	Use:   "comment [project-key] [repo-slug] [pr-id]",
	Short: "Add a comment to a pull request",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb pr comment L3SUP agents-sre 42 --body "LGTM!"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		prID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid PR ID: %s", args[2])
		}

		body, _ := cmd.Flags().GetString("body")
		comment, err := client.CommentPullRequest(args[0], args[1], prID, body)
		if err != nil {
			return err
		}

		fmt.Printf("Comment #%d added\n", comment.ID)
		return nil
	},
}

var prActivityCmd = &cobra.Command{
	Use:   "activity [project-key] [repo-slug] [pr-id]",
	Short: "List pull request activity (comments, approvals, etc.)",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb pr activity L3SUP agents-sre 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		prID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid PR ID: %s", args[2])
		}

		limit, _ := cmd.Flags().GetInt("limit")
		activities, err := client.ListPRActivities(args[0], args[1], prID, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(activities, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, a := range activities {
			date := ""
			if a.CreatedDate > 0 {
				date = time.UnixMilli(a.CreatedDate).Format("2006-01-02 15:04")
			}
			user := ""
			if a.User != nil {
				user = a.User.DisplayName
			}

			switch a.Action {
			case "COMMENTED":
				text := ""
				if a.Comment != nil {
					text = firstLine(a.Comment.Text)
					if len(text) > 80 {
						text = text[:77] + "..."
					}
				}
				fmt.Printf("[%s] %s commented: %s\n", date, user, text)
			case "APPROVED":
				fmt.Printf("[%s] %s approved\n", date, user)
			case "REVIEWED":
				fmt.Printf("[%s] %s reviewed (needs work)\n", date, user)
			case "MERGED":
				fmt.Printf("[%s] %s merged\n", date, user)
			case "DECLINED":
				fmt.Printf("[%s] %s declined\n", date, user)
			case "OPENED":
				fmt.Printf("[%s] %s opened\n", date, user)
			default:
				fmt.Printf("[%s] %s %s\n", date, user, a.Action)
			}
		}
		return nil
	},
}

func init() {
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prViewCmd)
	prCmd.AddCommand(prCreateCmd)
	prCmd.AddCommand(prMergeCmd)
	prCmd.AddCommand(prDeclineCmd)
	prCmd.AddCommand(prCommentCmd)
	prCmd.AddCommand(prActivityCmd)

	prListCmd.Flags().String("state", "", "filter: OPEN, MERGED, DECLINED, ALL")
	prListCmd.Flags().Int("limit", 25, "max results")

	prCreateCmd.Flags().String("title", "", "PR title (required)")
	prCreateCmd.Flags().String("description", "", "PR description")
	prCreateCmd.Flags().String("from", "", "source branch (required)")
	prCreateCmd.Flags().String("to", "", "target branch (required)")
	prCreateCmd.Flags().String("reviewers", "", "comma-separated reviewer usernames")
	prCreateCmd.MarkFlagRequired("title")
	prCreateCmd.MarkFlagRequired("from")
	prCreateCmd.MarkFlagRequired("to")

	prCommentCmd.Flags().String("body", "", "comment text (required)")
	prCommentCmd.MarkFlagRequired("body")

	prActivityCmd.Flags().Int("limit", 50, "max results")
}

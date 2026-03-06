package github

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
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
	Use:   "list [owner/repo]",
	Short: "List pull requests",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github pr list octocat/hello-world
  orbit gh pr list octocat/hello-world --state closed`,
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
		limit, _ := cmd.Flags().GetInt("limit")
		prs, err := client.ListPullRequests(owner, repo, state, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(prs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-6s %-8s %-50s %s\n", "#", "STATE", "TITLE", "AUTHOR")
		fmt.Printf("%-6s %-8s %-50s %s\n", "-", "-----", "-----", "------")
		for _, pr := range prs {
			title := pr.Title
			if len(title) > 48 {
				title = title[:45] + "..."
			}
			author := ""
			if pr.User != nil {
				author = pr.User.Login
			}
			state := pr.State
			if pr.Draft {
				state = "draft"
			}
			fmt.Printf("#%-5d %-8s %-50s %s\n", pr.Number, state, title, author)
		}
		return nil
	},
}

var prViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [number]",
	Short: "View a pull request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github pr view octocat/hello-world 42`,
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
			return fmt.Errorf("invalid PR number: %s", args[1])
		}

		pr, err := client.GetPullRequest(owner, repo, number)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(pr, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("PR:       #%d\n", pr.Number)
		fmt.Printf("Title:    %s\n", pr.Title)
		fmt.Printf("State:    %s\n", pr.State)
		if pr.Head != nil {
			fmt.Printf("Head:     %s\n", pr.Head.Ref)
		}
		if pr.Base != nil {
			fmt.Printf("Base:     %s\n", pr.Base.Ref)
		}
		if pr.User != nil {
			fmt.Printf("Author:   %s\n", pr.User.Login)
		}
		if len(pr.Labels) > 0 {
			names := make([]string, len(pr.Labels))
			for i, l := range pr.Labels {
				names[i] = l.Name
			}
			fmt.Printf("Labels:   %s\n", strings.Join(names, ", "))
		}
		fmt.Printf("Comments: %d\n", pr.Comments)
		fmt.Printf("URL:      %s\n", pr.HTMLURL)
		if pr.Body != "" {
			fmt.Printf("\n%s\n", pr.Body)
		}
		return nil
	},
}

var prCreateCmd = &cobra.Command{
	Use:   "create [owner/repo]",
	Short: "Create a pull request",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github pr create octocat/hello-world --head feature/x --base main --title "Add feature X"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		head, _ := cmd.Flags().GetString("head")
		base, _ := cmd.Flags().GetString("base")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")

		pr, err := client.CreatePullRequest(owner, repo, head, base, title, body)
		if err != nil {
			return err
		}

		fmt.Printf("Created PR #%d: %s\n", pr.Number, pr.Title)
		fmt.Printf("URL: %s\n", pr.HTMLURL)
		return nil
	},
}

var prMergeCmd = &cobra.Command{
	Use:   "merge [owner/repo] [number]",
	Short: "Merge a pull request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github pr merge octocat/hello-world 42
  orbit gh pr merge octocat/hello-world 42 --method squash`,
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
			return fmt.Errorf("invalid PR number: %s", args[1])
		}

		method, _ := cmd.Flags().GetString("method")
		if err := client.MergePullRequest(owner, repo, number, method); err != nil {
			return err
		}

		fmt.Printf("Merged PR #%d\n", number)
		return nil
	},
}

var prCommentCmd = &cobra.Command{
	Use:   "comment [owner/repo] [number]",
	Short: "Add a comment to a pull request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github pr comment octocat/hello-world 42 --body "LGTM!"`,
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
			return fmt.Errorf("invalid PR number: %s", args[1])
		}

		body, _ := cmd.Flags().GetString("body")
		comment, err := client.CreatePRComment(owner, repo, number, body)
		if err != nil {
			return err
		}

		fmt.Printf("Added comment #%d to PR #%d\n", comment.ID, number)
		return nil
	},
}

var prCommentsCmd = &cobra.Command{
	Use:   "comments [owner/repo] [number]",
	Short: "List comments on a pull request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github pr comments octocat/hello-world 42`,
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
			return fmt.Errorf("invalid PR number: %s", args[1])
		}

		limit, _ := cmd.Flags().GetInt("limit")
		comments, err := client.ListPRComments(owner, repo, number, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(comments, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, c := range comments {
			author := ""
			if c.User != nil {
				author = c.User.Login
			}
			date := ""
			if len(c.CreatedAt) >= 10 {
				date = c.CreatedAt[:10]
			}
			fmt.Printf("--- #%d by %s on %s ---\n%s\n\n", c.ID, author, date, c.Body)
		}
		return nil
	},
}

func init() {
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prViewCmd)
	prCmd.AddCommand(prCreateCmd)
	prCmd.AddCommand(prMergeCmd)
	prCmd.AddCommand(prCommentCmd)
	prCmd.AddCommand(prCommentsCmd)

	prListCmd.Flags().String("state", "", "filter by state: open, closed, all")
	prListCmd.Flags().Int("limit", 20, "max results")

	prCreateCmd.Flags().String("head", "", "head branch (required)")
	prCreateCmd.Flags().String("base", "", "base branch (required)")
	prCreateCmd.Flags().String("title", "", "PR title (required)")
	prCreateCmd.Flags().String("body", "", "PR body")
	_ = prCreateCmd.MarkFlagRequired("head")
	_ = prCreateCmd.MarkFlagRequired("base")
	_ = prCreateCmd.MarkFlagRequired("title")

	prMergeCmd.Flags().String("method", "", "merge method: merge, squash, rebase")

	prCommentCmd.Flags().String("body", "", "comment body (required)")
	_ = prCommentCmd.MarkFlagRequired("body")

	prCommentsCmd.Flags().Int("limit", 50, "max results")
}

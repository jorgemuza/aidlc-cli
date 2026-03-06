package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var mrCmd = &cobra.Command{
	Use:   "mr [subcommand]",
	Short: "Manage merge requests",
	Aliases: []string{"merge-request"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var mrListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List merge requests",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab mr list 595
  orbit gitlab mr list 595 --state opened
  orbit gitlab mr list 595 --state merged`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		mrs, err := client.ListMergeRequests(args[0], state, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(mrs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-6s %-8s %-50s %s\n", "IID", "STATE", "TITLE", "AUTHOR")
		fmt.Printf("%-6s %-8s %-50s %s\n", "---", "-----", "-----", "------")
		for _, mr := range mrs {
			title := mr.Title
			if len(title) > 48 {
				title = title[:45] + "..."
			}
			author := ""
			if mr.Author != nil {
				author = mr.Author.Username
			}
			state := mr.State
			if mr.Draft {
				state = "draft"
			}
			fmt.Printf("!%-5d %-8s %-50s %s\n", mr.IID, state, title, author)
		}
		return nil
	},
}

var mrViewCmd = &cobra.Command{
	Use:   "view [project] [mr-iid]",
	Short: "View a merge request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab mr view 595 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid MR IID: %s", args[1])
		}

		mr, err := client.GetMergeRequest(args[0], iid)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(mr, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("MR:          !%d\n", mr.IID)
		fmt.Printf("Title:       %s\n", mr.Title)
		fmt.Printf("State:       %s\n", mr.State)
		fmt.Printf("Source:      %s\n", mr.SourceBranch)
		fmt.Printf("Target:      %s\n", mr.TargetBranch)
		if mr.Author != nil {
			fmt.Printf("Author:      %s\n", mr.Author.Username)
		}
		if mr.Assignee != nil {
			fmt.Printf("Assignee:    %s\n", mr.Assignee.Username)
		}
		if len(mr.Labels) > 0 {
			fmt.Printf("Labels:      %s\n", strings.Join(mr.Labels, ", "))
		}
		fmt.Printf("Merge:       %s\n", mr.MergeStatus)
		fmt.Printf("Conflicts:   %v\n", mr.HasConflicts)
		fmt.Printf("Comments:    %d\n", mr.UserNotesCount)
		fmt.Printf("URL:         %s\n", mr.WebURL)
		if mr.Description != "" {
			fmt.Printf("\n%s\n", mr.Description)
		}
		return nil
	},
}

var mrCreateCmd = &cobra.Command{
	Use:   "create [project]",
	Short: "Create a merge request",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab mr create 595 --source feature/x --target main --title "Add feature X"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		source, _ := cmd.Flags().GetString("source")
		target, _ := cmd.Flags().GetString("target")
		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("description")

		mr, err := client.CreateMergeRequest(args[0], source, target, title, desc)
		if err != nil {
			return err
		}

		fmt.Printf("Created MR !%d: %s\n", mr.IID, mr.Title)
		fmt.Printf("URL: %s\n", mr.WebURL)
		return nil
	},
}

var mrMergeCmd = &cobra.Command{
	Use:   "merge [project] [mr-iid]",
	Short: "Merge a merge request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab mr merge 595 42
  orbit gitlab mr merge 595 42 --squash`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid MR IID: %s", args[1])
		}

		squash, _ := cmd.Flags().GetBool("squash")
		mr, err := client.MergeMergeRequest(args[0], iid, squash)
		if err != nil {
			return err
		}

		fmt.Printf("Merged MR !%d: %s\n", mr.IID, mr.Title)
		return nil
	},
}

var mrCommentCmd = &cobra.Command{
	Use:   "comment [project] [mr-iid]",
	Short: "Add a comment to a merge request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab mr comment 595 42 --body "LGTM!"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid MR IID: %s", args[1])
		}

		body, _ := cmd.Flags().GetString("body")
		note, err := client.CreateMRNote(args[0], iid, body)
		if err != nil {
			return err
		}

		fmt.Printf("Added comment #%d to MR !%d\n", note.ID, iid)
		return nil
	},
}

var mrNotesCmd = &cobra.Command{
	Use:   "notes [project] [mr-iid]",
	Short: "List comments on a merge request",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab mr notes 595 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		iid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid MR IID: %s", args[1])
		}

		limit, _ := cmd.Flags().GetInt("limit")
		notes, err := client.ListMRNotes(args[0], iid, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(notes, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, n := range notes {
			if n.System {
				continue
			}
			author := ""
			if n.Author != nil {
				author = n.Author.Username
			}
			date := ""
			if len(n.CreatedAt) >= 10 {
				date = n.CreatedAt[:10]
			}
			fmt.Printf("--- #%d by %s on %s ---\n%s\n\n", n.ID, author, date, n.Body)
		}
		return nil
	},
}

func init() {
	mrCmd.AddCommand(mrListCmd)
	mrCmd.AddCommand(mrViewCmd)
	mrCmd.AddCommand(mrCreateCmd)
	mrCmd.AddCommand(mrMergeCmd)
	mrCmd.AddCommand(mrCommentCmd)
	mrCmd.AddCommand(mrNotesCmd)

	mrListCmd.Flags().String("state", "", "filter by state: opened, closed, merged, all")
	mrListCmd.Flags().Int("limit", 20, "max results")

	mrCreateCmd.Flags().String("source", "", "source branch (required)")
	mrCreateCmd.Flags().String("target", "", "target branch (required)")
	mrCreateCmd.Flags().String("title", "", "MR title (required)")
	mrCreateCmd.Flags().String("description", "", "MR description")
	_ = mrCreateCmd.MarkFlagRequired("source")
	_ = mrCreateCmd.MarkFlagRequired("target")
	_ = mrCreateCmd.MarkFlagRequired("title")

	mrMergeCmd.Flags().Bool("squash", false, "squash commits")

	mrCommentCmd.Flags().String("body", "", "comment body (required)")
	_ = mrCommentCmd.MarkFlagRequired("body")

	mrNotesCmd.Flags().Int("limit", 50, "max results")
}

package github

import (
	"encoding/json"
	"fmt"
	"strconv"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run [subcommand]",
	Short:   "Manage GitHub Actions workflow runs",
	Aliases: []string{"actions"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var runListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List workflow runs",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github run list octocat/hello-world
  orbit gh run list octocat/hello-world --branch main --status success`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		branch, _ := cmd.Flags().GetString("branch")
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		runs, err := client.ListWorkflowRuns(owner, repo, branch, status, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(runs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-12s %-12s %-12s %-20s %-10s %s\n", "ID", "STATUS", "CONCLUSION", "BRANCH", "EVENT", "CREATED")
		fmt.Printf("%-12s %-12s %-12s %-20s %-10s %s\n", "--", "------", "----------", "------", "-----", "-------")
		for _, r := range runs {
			created := ""
			if len(r.CreatedAt) >= 10 {
				created = r.CreatedAt[:10]
			}
			branch := r.HeadBranch
			if len(branch) > 18 {
				branch = branch[:15] + "..."
			}
			fmt.Printf("%-12d %-12s %-12s %-20s %-10s %s\n", r.ID, r.Status, r.Conclusion, branch, r.Event, created)
		}
		return nil
	},
}

var runViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [run-id]",
	Short: "View a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run view octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		r, err := client.GetWorkflowRun(owner, repo, id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:         %d\n", r.ID)
		fmt.Printf("Name:       %s\n", r.Name)
		fmt.Printf("Status:     %s\n", r.Status)
		fmt.Printf("Conclusion: %s\n", r.Conclusion)
		fmt.Printf("Branch:     %s\n", r.HeadBranch)
		fmt.Printf("SHA:        %s\n", r.HeadSHA)
		fmt.Printf("Event:      %s\n", r.Event)
		fmt.Printf("Created:    %s\n", r.CreatedAt)
		fmt.Printf("URL:        %s\n", r.HTMLURL)
		return nil
	},
}

var runCancelCmd = &cobra.Command{
	Use:   "cancel [owner/repo] [run-id]",
	Short: "Cancel a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run cancel octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		if err := client.CancelWorkflowRun(owner, repo, id); err != nil {
			return err
		}

		fmt.Printf("Canceled workflow run %d\n", id)
		return nil
	},
}

var runRerunCmd = &cobra.Command{
	Use:   "rerun [owner/repo] [run-id]",
	Short: "Re-run a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run rerun octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		if err := client.RerunWorkflowRun(owner, repo, id); err != nil {
			return err
		}

		fmt.Printf("Re-running workflow run %d\n", id)
		return nil
	},
}

func init() {
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runViewCmd)
	runCmd.AddCommand(runCancelCmd)
	runCmd.AddCommand(runRerunCmd)

	runListCmd.Flags().String("branch", "", "filter by branch")
	runListCmd.Flags().String("status", "", "filter by status: completed, in_progress, queued")
	runListCmd.Flags().Int("limit", 20, "max results")
}

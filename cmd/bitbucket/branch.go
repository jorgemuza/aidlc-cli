package bitbucket

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [subcommand]",
	Short: "Manage branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var branchListCmd = &cobra.Command{
	Use:   "list [project-key] [repo-slug]",
	Short: "List branches",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc bitbucket branch list L3SUP agents-sre
  aidlc bb branch list L3SUP agents-sre --filter main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		filter, _ := cmd.Flags().GetString("filter")
		limit, _ := cmd.Flags().GetInt("limit")
		branches, err := client.ListBranches(args[0], args[1], filter, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branches, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-40s %-7s %s\n", "NAME", "DEFAULT", "LATEST COMMIT")
		fmt.Printf("%-40s %-7s %s\n", "----", "-------", "-------------")
		for _, b := range branches {
			commit := b.LatestCommit
			if len(commit) > 12 {
				commit = commit[:12]
			}
			fmt.Printf("%-40s %-7v %s\n", b.DisplayID, b.IsDefault, commit)
		}
		return nil
	},
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [project-key] [repo-slug] [name] [start-point]",
	Short: "Create a branch",
	Args:  cobra.ExactArgs(4),
	Example: `  aidlc bb branch create L3SUP agents-sre feature/new-thing main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		if err := client.CreateBranch(args[0], args[1], args[2], args[3]); err != nil {
			return err
		}

		fmt.Printf("Created branch: %s\n", args[2])
		return nil
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete [project-key] [repo-slug] [branch-name]",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc bb branch delete L3SUP agents-sre feature/old-thing`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		if err := client.DeleteBranch(args[0], args[1], args[2]); err != nil {
			return err
		}

		fmt.Printf("Deleted branch: %s\n", args[2])
		return nil
	},
}

var branchDefaultCmd = &cobra.Command{
	Use:   "default [project-key] [repo-slug]",
	Short: "Show default branch",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc bb branch default L3SUP agents-sre`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		branch, err := client.GetDefaultBranch(args[0], args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branch, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Name:   %s\n", branch.DisplayID)
		fmt.Printf("Commit: %s\n", branch.LatestCommit)
		return nil
	},
}

func init() {
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchCreateCmd)
	branchCmd.AddCommand(branchDeleteCmd)
	branchCmd.AddCommand(branchDefaultCmd)

	branchListCmd.Flags().String("filter", "", "filter branches by name")
	branchListCmd.Flags().Int("limit", 50, "max results")
}

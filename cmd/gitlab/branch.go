package gitlab

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
	Use:   "list [project]",
	Short: "List branches",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab branch list 595
  aidlc gitlab branch list schools/frontend/schools-frontend-backoffice --search main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		search, _ := cmd.Flags().GetString("search")
		limit, _ := cmd.Flags().GetInt("limit")
		branches, err := client.ListBranches(args[0], search, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branches, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-40s %-9s %-7s %s\n", "NAME", "PROTECTED", "DEFAULT", "LAST COMMIT")
		fmt.Printf("%-40s %-9s %-7s %s\n", "----", "---------", "-------", "-----------")
		for _, b := range branches {
			lastCommit := ""
			if b.Commit != nil && len(b.Commit.CommittedDate) >= 10 {
				lastCommit = b.Commit.ShortID + " " + b.Commit.CommittedDate[:10]
			}
			fmt.Printf("%-40s %-9v %-7v %s\n", b.Name, b.Protected, b.Default, lastCommit)
		}
		return nil
	},
}

var branchViewCmd = &cobra.Command{
	Use:   "view [project] [branch]",
	Short: "View a branch",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab branch view 595 main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		branch, err := client.GetBranch(args[0], args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(branch, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Name:      %s\n", branch.Name)
		fmt.Printf("Protected: %v\n", branch.Protected)
		fmt.Printf("Default:   %v\n", branch.Default)
		if branch.Commit != nil {
			fmt.Printf("Commit:    %s\n", branch.Commit.ShortID)
			fmt.Printf("Author:    %s\n", branch.Commit.AuthorName)
			fmt.Printf("Date:      %s\n", branch.Commit.CommittedDate)
			fmt.Printf("Message:   %s\n", branch.Commit.Title)
		}
		return nil
	},
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [project] [name] [ref]",
	Short: "Create a branch",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc gitlab branch create 595 feature/new-feature main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		branch, err := client.CreateBranch(args[0], args[1], args[2])
		if err != nil {
			return err
		}

		fmt.Printf("Created branch: %s\n", branch.Name)
		return nil
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete [project] [branch]",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab branch delete 595 feature/old-feature`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		if err := client.DeleteBranch(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Deleted branch: %s\n", args[1])
		return nil
	},
}

func init() {
	branchCmd.AddCommand(branchListCmd)
	branchCmd.AddCommand(branchViewCmd)
	branchCmd.AddCommand(branchCreateCmd)
	branchCmd.AddCommand(branchDeleteCmd)

	branchListCmd.Flags().String("search", "", "search branches by name")
	branchListCmd.Flags().Int("limit", 50, "max results")
}

package gitlab

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file [subcommand]",
	Short: "View repository files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var fileReadCmd = &cobra.Command{
	Use:   "read [project] [file-path]",
	Short: "Read a file from the repository",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab file read foundation/ai .gitlab-ci.yml
  orbit gitlab file read foundation/ai src/main.go --ref develop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		ref, _ := cmd.Flags().GetString("ref")
		if ref == "" {
			ref = "main"
		}

		content, err := client.GetFileRaw(args[0], args[1], ref)
		if err != nil {
			return err
		}

		fmt.Print(content)
		return nil
	},
}

var fileUpdateCmd = &cobra.Command{
	Use:   "update [project] [file-path]",
	Short: "Update a file in the repository",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab file update foundation/ai .gitlab-ci.yml --body-file ci.yml --message "fix: add runner tags"
  orbit gitlab file update foundation/ai README.md --body "# Hello" --branch main --message "update readme"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		branch, _ := cmd.Flags().GetString("branch")
		message, _ := cmd.Flags().GetString("message")
		body, _ := cmd.Flags().GetString("body")
		bodyFile, _ := cmd.Flags().GetString("body-file")

		if body == "" && bodyFile == "" {
			return fmt.Errorf("either --body or --body-file is required")
		}
		if bodyFile != "" {
			data, err := os.ReadFile(bodyFile)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			body = string(data)
		}

		if err := client.UpdateFile(args[0], args[1], branch, body, message); err != nil {
			return err
		}

		fmt.Printf("Updated %s on branch %s\n", args[1], branch)
		return nil
	},
}

func init() {
	fileCmd.AddCommand(fileReadCmd)
	fileCmd.AddCommand(fileUpdateCmd)

	fileReadCmd.Flags().String("ref", "main", "branch, tag, or commit SHA")

	fileUpdateCmd.Flags().String("branch", "main", "target branch")
	fileUpdateCmd.Flags().String("message", "", "commit message (required)")
	fileUpdateCmd.Flags().String("body", "", "file content as string")
	fileUpdateCmd.Flags().String("body-file", "", "read content from local file")
	fileUpdateCmd.MarkFlagRequired("message")
}

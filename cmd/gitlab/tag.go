package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag [subcommand]",
	Short: "Manage tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List tags",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab tag list 595`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		tags, err := client.ListTags(args[0], limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(tags, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-30s %-10s %s\n", "TAG", "COMMIT", "MESSAGE")
		fmt.Printf("%-30s %-10s %s\n", "---", "------", "-------")
		for _, t := range tags {
			commitID := ""
			if t.Commit != nil {
				commitID = t.Commit.ShortID
			}
			fmt.Printf("%-30s %-10s %s\n", t.Name, commitID, t.Message)
		}
		return nil
	},
}

var tagCreateCmd = &cobra.Command{
	Use:   "create [project] [tag-name] [ref]",
	Short: "Create a tag",
	Args:  cobra.ExactArgs(3),
	Example: `  aidlc gitlab tag create 595 v1.0.0 main
  aidlc gitlab tag create 595 v1.0.0 main --message "Release v1.0.0"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		message, _ := cmd.Flags().GetString("message")
		tag, err := client.CreateTag(args[0], args[1], args[2], message)
		if err != nil {
			return err
		}

		fmt.Printf("Created tag: %s\n", tag.Name)
		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagCreateCmd)

	tagListCmd.Flags().Int("limit", 50, "max results")
	tagCreateCmd.Flags().StringP("message", "m", "", "tag message")
}

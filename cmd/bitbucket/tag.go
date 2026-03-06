package bitbucket

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
	Use:   "list [project-key] [repo-slug]",
	Short: "List tags",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit bb tag list L3SUP agents-sre`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		filter, _ := cmd.Flags().GetString("filter")
		limit, _ := cmd.Flags().GetInt("limit")
		tags, err := client.ListTags(args[0], args[1], filter, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(tags, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-30s %s\n", "NAME", "LATEST COMMIT")
		fmt.Printf("%-30s %s\n", "----", "-------------")
		for _, t := range tags {
			commit := t.LatestCommit
			if len(commit) > 12 {
				commit = commit[:12]
			}
			fmt.Printf("%-30s %s\n", t.DisplayID, commit)
		}
		return nil
	},
}

var tagCreateCmd = &cobra.Command{
	Use:   "create [project-key] [repo-slug] [tag-name] [start-point]",
	Short: "Create a tag",
	Args:  cobra.ExactArgs(4),
	Example: `  orbit bb tag create L3SUP agents-sre v1.0.0 main`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		message, _ := cmd.Flags().GetString("message")
		if err := client.CreateTag(args[0], args[1], args[2], args[3], message); err != nil {
			return err
		}

		fmt.Printf("Created tag: %s\n", args[2])
		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagCreateCmd)

	tagListCmd.Flags().String("filter", "", "filter tags by name")
	tagListCmd.Flags().Int("limit", 50, "max results")
	tagCreateCmd.Flags().StringP("message", "m", "", "tag message")
}

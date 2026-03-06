package confluence

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var pageViewCmd = &cobra.Command{
	Use:   "page [page-id]",
	Short: "View a Confluence page",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc confluence page 12345
  aidlc confluence page 12345 -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		page, err := client.GetPage(args[0])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(page, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:      %s\n", page.ID)
		fmt.Printf("Title:   %s\n", page.Title)
		fmt.Printf("Version: %d\n", page.Version.Number)
		if page.Space != nil {
			fmt.Printf("Space:   %s\n", page.Space.Key)
		}
		if page.Links != nil && page.Links.Base != "" {
			fmt.Printf("URL:     %s%s\n", page.Links.Base, page.Links.WebUI)
		}
		return nil
	},
}

var pageChildrenCmd = &cobra.Command{
	Use:   "children [page-id]",
	Short: "List child pages of a Confluence page",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc confluence children 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		children, err := client.GetChildPages(args[0])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(children, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-15s %-8s %s\n", "ID", "VERSION", "TITLE")
		fmt.Printf("%-15s %-8s %s\n", "--", "-------", "-----")
		for _, p := range children {
			ver := 0
			if p.Version != nil {
				ver = p.Version.Number
			}
			fmt.Printf("%-15s %-8d %s\n", p.ID, ver, p.Title)
		}
		return nil
	},
}

var pageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Confluence page",
	Example: `  aidlc confluence create --space FO --parent 12345 --title "My Page" --body "<p>Hello</p>"
  aidlc confluence create --space FO --parent 12345 --title "My Page" --file doc.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		space, _ := cmd.Flags().GetString("space")
		parent, _ := cmd.Flags().GetString("parent")
		body, _ := cmd.Flags().GetString("body")
		file, _ := cmd.Flags().GetString("file")

		if file != "" {
			content, err := readAndConvertMarkdown(file)
			if err != nil {
				return err
			}
			body = content
		}

		page, err := client.CreatePage(space, parent, title, body)
		if err != nil {
			return err
		}

		fmt.Printf("Created page %s: %s\n", page.ID, page.Title)
		if page.Links != nil && page.Links.Base != "" {
			fmt.Printf("URL: %s%s\n", page.Links.Base, page.Links.WebUI)
		}
		return nil
	},
}

var pageUpdateCmd = &cobra.Command{
	Use:   "update [page-id]",
	Short: "Update an existing Confluence page",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc confluence update 12345 --title "New Title" --body "<p>Updated</p>"
  aidlc confluence update 12345 --file doc.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		// Get current page for version
		current, err := client.GetPage(args[0])
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			title = current.Title
		}

		body, _ := cmd.Flags().GetString("body")
		file, _ := cmd.Flags().GetString("file")

		if file != "" {
			content, err := readAndConvertMarkdown(file)
			if err != nil {
				return err
			}
			body = content
		}

		page, err := client.UpdatePage(args[0], title, body, current.Version.Number+1)
		if err != nil {
			return err
		}

		fmt.Printf("Updated page %s (v%d): %s\n", page.ID, page.Version.Number, page.Title)
		return nil
	},
}

var pageDeleteCmd = &cobra.Command{
	Use:   "delete [page-id]",
	Short: "Delete a Confluence page (moves to trash)",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc confluence delete 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		// Fetch page title for confirmation message
		page, _ := client.GetPage(args[0])
		title := args[0]
		if page != nil {
			title = page.Title
		}

		if err := client.DeletePage(args[0]); err != nil {
			return err
		}

		fmt.Printf("Deleted page %s: %s\n", args[0], title)
		return nil
	},
}

func init() {
	// create flags
	pageCreateCmd.Flags().String("title", "", "page title (required)")
	pageCreateCmd.Flags().String("space", "", "space key (required)")
	pageCreateCmd.Flags().String("parent", "", "parent page ID")
	pageCreateCmd.Flags().StringP("body", "b", "", "page body in storage format")
	pageCreateCmd.Flags().StringP("file", "f", "", "markdown file to convert and upload")
	_ = pageCreateCmd.MarkFlagRequired("title")
	_ = pageCreateCmd.MarkFlagRequired("space")

	// update flags
	pageUpdateCmd.Flags().String("title", "", "new title (keeps current if empty)")
	pageUpdateCmd.Flags().StringP("body", "b", "", "new body in storage format")
	pageUpdateCmd.Flags().StringP("file", "f", "", "markdown file to convert and upload")
}

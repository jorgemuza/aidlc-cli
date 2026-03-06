package confluence

import (
	"fmt"

	conflsvc "github.com/jorgemuza/orbit/internal/service/confluence"
	"github.com/spf13/cobra"
)

var setWidthCmd = &cobra.Command{
	Use:   "set-width [page-id...]",
	Short: "Set page width (wide or fixed)",
	Args:  cobra.MinimumNArgs(1),
	Example: `  orbit confluence set-width 12345
  orbit confluence set-width 12345 67890 --width fixed
  orbit confluence set-width 12345 --recursive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveConfluenceClient(cmd)
		if err != nil {
			return err
		}

		width, _ := cmd.Flags().GetString("width")
		recursive, _ := cmd.Flags().GetBool("recursive")

		appearance := "full-width"
		if width == "fixed" {
			appearance = "fixed"
		}

		for _, pageID := range args {
			if err := setWidthOnPage(client, pageID, appearance, recursive); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	setWidthCmd.Flags().String("width", "wide", "page width: wide or fixed")
	setWidthCmd.Flags().Bool("recursive", false, "apply to all child pages recursively")
}

func setWidthOnPage(client *conflsvc.Client, pageID, appearance string, recursive bool) error {
	page, err := client.GetPage(pageID)
	if err != nil {
		return err
	}

	if err := client.SetPageWidth(pageID, appearance); err != nil {
		return fmt.Errorf("setting width on %s: %w", pageID, err)
	}
	fmt.Printf("  ✅ %s — %s\n", pageID, page.Title)

	if recursive {
		children, err := client.GetChildPages(pageID)
		if err != nil {
			return err
		}
		for _, child := range children {
			if err := setWidthOnPage(client, child.ID, appearance, true); err != nil {
				return err
			}
		}
	}
	return nil
}

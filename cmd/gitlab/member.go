package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{
	Use:   "member [subcommand]",
	Short: "Manage project members",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var memberListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List project members",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab member list 595`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		members, err := client.ListProjectMembers(args[0], limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(members, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-20s %-30s %s\n", "ID", "USERNAME", "NAME", "ACCESS")
		fmt.Printf("%-10s %-20s %-30s %s\n", "--", "--------", "----", "------")
		for _, m := range members {
			access := accessLevelName(m.AccessLevel)
			fmt.Printf("%-10d %-20s %-30s %s\n", m.ID, m.Username, m.Name, access)
		}
		return nil
	},
}

func accessLevelName(level int) string {
	switch level {
	case 10:
		return "Guest"
	case 20:
		return "Reporter"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	default:
		return fmt.Sprintf("Level %d", level)
	}
}

func init() {
	memberCmd.AddCommand(memberListCmd)
	memberListCmd.Flags().Int("limit", 50, "max results")
}

package bitbucket

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user [subcommand]",
	Short: "User operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var userListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List users",
	Example: `  aidlc bb user list --filter john`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveBBClient(cmd)
		if err != nil {
			return err
		}

		filter, _ := cmd.Flags().GetString("filter")
		limit, _ := cmd.Flags().GetInt("limit")
		users, err := client.ListUsers(filter, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(users, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-20s %-30s %s\n", "ID", "USERNAME", "DISPLAY NAME", "ACTIVE")
		fmt.Printf("%-10s %-20s %-30s %s\n", "--", "--------", "------------", "------")
		for _, u := range users {
			fmt.Printf("%-10d %-20s %-30s %v\n", u.ID, u.Name, u.DisplayName, u.Active)
		}
		return nil
	},
}

func init() {
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().String("filter", "", "filter by username or display name")
	userListCmd.Flags().Int("limit", 25, "max results")
}

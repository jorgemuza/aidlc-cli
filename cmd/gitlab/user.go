package gitlab

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

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		user, err := client.CurrentUser()
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(user, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:       %d\n", user.ID)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Name:     %s\n", user.Name)
		fmt.Printf("Email:    %s\n", user.Email)
		fmt.Printf("State:    %s\n", user.State)
		fmt.Printf("URL:      %s\n", user.WebURL)
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Example: `  aidlc gitlab user list --search john`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		search, _ := cmd.Flags().GetString("search")
		limit, _ := cmd.Flags().GetInt("limit")
		users, err := client.ListUsers(search, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(users, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-20s %-30s %s\n", "ID", "USERNAME", "NAME", "STATE")
		fmt.Printf("%-10s %-20s %-30s %s\n", "--", "--------", "----", "-----")
		for _, u := range users {
			fmt.Printf("%-10d %-20s %-30s %s\n", u.ID, u.Username, u.Name, u.State)
		}
		return nil
	},
}

func init() {
	userCmd.AddCommand(userMeCmd)
	userCmd.AddCommand(userListCmd)

	userListCmd.Flags().String("search", "", "search by username or name")
	userListCmd.Flags().Int("limit", 20, "max results")
}

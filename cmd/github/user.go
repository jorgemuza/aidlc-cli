package github

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
		client, err := resolveGitHubClient(cmd)
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
		fmt.Printf("Login:    %s\n", user.Login)
		fmt.Printf("Name:     %s\n", user.Name)
		fmt.Printf("Email:    %s\n", user.Email)
		fmt.Printf("Type:     %s\n", user.Type)
		fmt.Printf("URL:      %s\n", user.HTMLURL)
		return nil
	},
}

var userViewCmd = &cobra.Command{
	Use:   "view [username]",
	Short: "View a user profile",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github user view octocat`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		user, err := client.GetUser(args[0])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(user, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:    %d\n", user.ID)
		fmt.Printf("Login: %s\n", user.Login)
		fmt.Printf("Name:  %s\n", user.Name)
		fmt.Printf("Type:  %s\n", user.Type)
		fmt.Printf("URL:   %s\n", user.HTMLURL)
		return nil
	},
}

func init() {
	userCmd.AddCommand(userMeCmd)
	userCmd.AddCommand(userViewCmd)
}

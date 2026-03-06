package github

import (
	"encoding/json"
	"fmt"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:   "secret [subcommand]",
	Short: "Manage GitHub Actions secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var secretListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List repository secrets",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github secret list octocat/hello-world`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		secrets, err := client.ListRepoSecrets(owner, repo, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(secrets, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-40s %-12s %s\n", "NAME", "CREATED", "UPDATED")
		fmt.Printf("%-40s %-12s %s\n", "----", "-------", "-------")
		for _, s := range secrets {
			created := ""
			if len(s.CreatedAt) >= 10 {
				created = s.CreatedAt[:10]
			}
			updated := ""
			if len(s.UpdatedAt) >= 10 {
				updated = s.UpdatedAt[:10]
			}
			fmt.Printf("%-40s %-12s %s\n", s.Name, created, updated)
		}
		return nil
	},
}

var secretSetCmd = &cobra.Command{
	Use:   "set [owner/repo] [name] [value]",
	Short: "Create or update a repository secret",
	Args:  cobra.ExactArgs(3),
	Example: `  orbit github secret set octocat/hello-world MY_SECRET "secret-value"
  orbit gh secret set octocat/hello-world DEPLOY_KEY "$(cat key.pem)"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		if err := client.SetRepoSecret(owner, repo, args[1], args[2]); err != nil {
			return err
		}

		fmt.Printf("Secret %s set for %s/%s\n", args[1], owner, repo)
		return nil
	},
}

var secretDeleteCmd = &cobra.Command{
	Use:   "delete [owner/repo] [name]",
	Short: "Delete a repository secret",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github secret delete octocat/hello-world MY_SECRET`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		if err := client.DeleteRepoSecret(owner, repo, args[1]); err != nil {
			return err
		}

		fmt.Printf("Deleted secret %s from %s/%s\n", args[1], owner, repo)
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretSetCmd)
	secretCmd.AddCommand(secretDeleteCmd)

	secretListCmd.Flags().Int("limit", 30, "max results")
}

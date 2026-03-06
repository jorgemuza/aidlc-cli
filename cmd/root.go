package cmd

import (
	cmdbitbucket "github.com/jorgemuza/orbit/cmd/bitbucket"
	cmdconfluence "github.com/jorgemuza/orbit/cmd/confluence"
	cmdgithub "github.com/jorgemuza/orbit/cmd/github"
	cmdgitlab "github.com/jorgemuza/orbit/cmd/gitlab"
	cmdjira "github.com/jorgemuza/orbit/cmd/jira"
	"github.com/jorgemuza/orbit/cmd/profile"
	cmdservice "github.com/jorgemuza/orbit/cmd/service"
	"github.com/jorgemuza/orbit/cmd/version"

	// Register all service types
	_ "github.com/jorgemuza/orbit/internal/service/bitbucket"
	_ "github.com/jorgemuza/orbit/internal/service/confluence"
	_ "github.com/jorgemuza/orbit/internal/service/github"
	_ "github.com/jorgemuza/orbit/internal/service/gitlab"
	_ "github.com/jorgemuza/orbit/internal/service/jira"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "orbit",
	Short: "AI Life Development Cycle CLI",
	Long: `orbit is a unified CLI for managing connections to development lifecycle services.

Supports Jira, Confluence, GitLab, GitHub, and Bitbucket (cloud and self-hosted).
Organize connections into profiles to switch between projects seamlessly.

Secrets can be stored as 1Password references (op://vault/item/field) and
are resolved at runtime using the 1Password CLI.`,
}

func SetVersion(ver, commit, date string) {
	version.Set(ver, commit, date)
	rootCmd.Version = ver
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default ~/.config/orbit/config.yaml)")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "profile to use (overrides default)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table, json, yaml")

	rootCmd.AddCommand(profile.Command)
	rootCmd.AddCommand(cmdservice.Command)
	rootCmd.AddCommand(cmdjira.Command)
	rootCmd.AddCommand(cmdconfluence.Command)
	rootCmd.AddCommand(cmdgitlab.Command)
	rootCmd.AddCommand(cmdgithub.Command)
	rootCmd.AddCommand(cmdbitbucket.Command)
	rootCmd.AddCommand(version.Command)

	rootCmd.SilenceUsage = true
}

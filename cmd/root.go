package cmd

import (
	cmdbitbucket "github.com/jorgemuza/aidlc-cli/cmd/bitbucket"
	cmdconfluence "github.com/jorgemuza/aidlc-cli/cmd/confluence"
	cmdgitlab "github.com/jorgemuza/aidlc-cli/cmd/gitlab"
	cmdjira "github.com/jorgemuza/aidlc-cli/cmd/jira"
	"github.com/jorgemuza/aidlc-cli/cmd/profile"
	cmdservice "github.com/jorgemuza/aidlc-cli/cmd/service"
	"github.com/jorgemuza/aidlc-cli/cmd/version"

	// Register all service types
	_ "github.com/jorgemuza/aidlc-cli/internal/service/bitbucket"
	_ "github.com/jorgemuza/aidlc-cli/internal/service/confluence"
	_ "github.com/jorgemuza/aidlc-cli/internal/service/gitlab"
	_ "github.com/jorgemuza/aidlc-cli/internal/service/jira"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aidlc",
	Short: "AI Life Development Cycle CLI",
	Long: `aidlc is a unified CLI for managing connections to development lifecycle services.

Supports Jira, Confluence, GitLab, and Bitbucket (cloud and self-hosted).
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
	rootCmd.PersistentFlags().String("config", "", "config file (default ~/.config/aidlc/config.yaml)")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "profile to use (overrides default)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table, json, yaml")

	rootCmd.AddCommand(profile.Command)
	rootCmd.AddCommand(cmdservice.Command)
	rootCmd.AddCommand(cmdjira.Command)
	rootCmd.AddCommand(cmdconfluence.Command)
	rootCmd.AddCommand(cmdgitlab.Command)
	rootCmd.AddCommand(cmdbitbucket.Command)
	rootCmd.AddCommand(version.Command)

	rootCmd.SilenceUsage = true
}

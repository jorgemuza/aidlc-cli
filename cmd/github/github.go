package github

import (
	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var serviceName string

// Command is the top-level github command.
var Command = &cobra.Command{
	Use:     "github",
	Short:   "Manage GitHub repositories, pull requests, issues, and more",
	Aliases: []string{"gh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.PersistentFlags().StringVar(&serviceName, "service", "", "github service name (if profile has multiple)")
	Command.AddCommand(repoCmd)
	Command.AddCommand(reposCmd)
	Command.AddCommand(branchCmd)
	Command.AddCommand(tagCmd)
	Command.AddCommand(commitCmd)
	Command.AddCommand(prCmd)
	Command.AddCommand(issueCmd)
	Command.AddCommand(releaseCmd)
	Command.AddCommand(runCmd)
	Command.AddCommand(secretCmd)
	Command.AddCommand(userCmd)
}

func resolveGitHubClient(cmd *cobra.Command) (*ghsvc.Client, error) {
	_, p, err := cmdutil.ResolveProfile(cmd)
	if err != nil {
		return nil, err
	}

	conn, err := cmdutil.FindServiceByTypeOrName(p, config.ServiceTypeGitHub, serviceName)
	if err != nil {
		return nil, err
	}

	svc, err := service.Create(*conn)
	if err != nil {
		return nil, err
	}

	return ghsvc.ClientFromService(svc)
}

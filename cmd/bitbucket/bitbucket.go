package bitbucket

import (
	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
	bbsvc "github.com/jorgemuza/orbit/internal/service/bitbucket"
	"github.com/spf13/cobra"
)

var serviceName string

// Command is the top-level bitbucket command.
var Command = &cobra.Command{
	Use:     "bitbucket",
	Short:   "Manage Bitbucket repositories, pull requests, branches, and more",
	Aliases: []string{"bb"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.PersistentFlags().StringVar(&serviceName, "service", "", "bitbucket service name (if profile has multiple)")
	Command.AddCommand(projectCmd)
	Command.AddCommand(repoCmd)
	Command.AddCommand(branchCmd)
	Command.AddCommand(tagCmd)
	Command.AddCommand(commitCmd)
	Command.AddCommand(prCmd)
	Command.AddCommand(userCmd)
}

func resolveBBClient(cmd *cobra.Command) (*bbsvc.Client, error) {
	_, p, err := cmdutil.ResolveProfile(cmd)
	if err != nil {
		return nil, err
	}

	conn, err := cmdutil.FindServiceByTypeOrName(p, config.ServiceTypeBitbucket, serviceName)
	if err != nil {
		return nil, err
	}

	svc, err := service.Create(*conn)
	if err != nil {
		return nil, err
	}

	return bbsvc.ClientFromService(svc)
}

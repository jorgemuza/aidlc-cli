package confluence

import (
	"github.com/jorgemuza/aidlc-cli/cmd/cmdutil"
	"github.com/jorgemuza/aidlc-cli/internal/config"
	"github.com/jorgemuza/aidlc-cli/internal/service"
	conflsvc "github.com/jorgemuza/aidlc-cli/internal/service/confluence"
	"github.com/spf13/cobra"
)

var serviceName string

// Command is the top-level confluence command.
var Command = &cobra.Command{
	Use:   "confluence",
	Short: "Manage Confluence pages and content",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.PersistentFlags().StringVar(&serviceName, "service", "", "confluence service name (if profile has multiple)")
	Command.AddCommand(pageViewCmd)
	Command.AddCommand(pageCreateCmd)
	Command.AddCommand(pageUpdateCmd)
	Command.AddCommand(pageChildrenCmd)
	Command.AddCommand(pageDeleteCmd)
	Command.AddCommand(publishCmd)
	Command.AddCommand(setWidthCmd)
}

func resolveConfluenceClient(cmd *cobra.Command) (*conflsvc.Client, error) {
	_, p, err := cmdutil.ResolveProfile(cmd)
	if err != nil {
		return nil, err
	}

	conn, err := cmdutil.FindServiceByTypeOrName(p, config.ServiceTypeConfluence, serviceName)
	if err != nil {
		return nil, err
	}

	svc, err := service.Create(*conn)
	if err != nil {
		return nil, err
	}

	return conflsvc.ClientFromService(svc)
}

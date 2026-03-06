package gitlab

import (
	"github.com/jorgemuza/aidlc-cli/cmd/cmdutil"
	"github.com/jorgemuza/aidlc-cli/internal/config"
	"github.com/jorgemuza/aidlc-cli/internal/service"
	glsvc "github.com/jorgemuza/aidlc-cli/internal/service/gitlab"
	"github.com/spf13/cobra"
)

var serviceName string

// Command is the top-level gitlab command.
var Command = &cobra.Command{
	Use:   "gitlab",
	Short: "Manage GitLab projects, merge requests, pipelines, and more",
	Aliases: []string{"gl"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.PersistentFlags().StringVar(&serviceName, "service", "", "gitlab service name (if profile has multiple)")
	Command.AddCommand(projectCmd)
	Command.AddCommand(projectsCmd)
	Command.AddCommand(groupCmd)
	Command.AddCommand(branchCmd)
	Command.AddCommand(tagCmd)
	Command.AddCommand(commitCmd)
	Command.AddCommand(mrCmd)
	Command.AddCommand(pipelineCmd)
	Command.AddCommand(jobCmd)
	Command.AddCommand(runnerCmd)
	Command.AddCommand(issueCmd)
	Command.AddCommand(fileCmd)
	Command.AddCommand(memberCmd)
	Command.AddCommand(userCmd)
	Command.AddCommand(variableCmd)
}

func resolveGitLabClient(cmd *cobra.Command) (*glsvc.Client, error) {
	_, p, err := cmdutil.ResolveProfile(cmd)
	if err != nil {
		return nil, err
	}

	conn, err := cmdutil.FindServiceByTypeOrName(p, config.ServiceTypeGitLab, serviceName)
	if err != nil {
		return nil, err
	}

	svc, err := service.Create(*conn)
	if err != nil {
		return nil, err
	}

	return glsvc.ClientFromService(svc)
}

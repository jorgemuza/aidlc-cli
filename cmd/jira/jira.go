package jira

import (
	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var serviceName string

// Command is the top-level jira command.
var Command = &cobra.Command{
	Use:   "jira",
	Short: "Manage Jira issues, epics, sprints, boards, projects, and releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.PersistentFlags().StringVar(&serviceName, "service", "", "jira service name (if profile has multiple)")
	Command.AddCommand(issueCmd)
	Command.AddCommand(epicCmd)
	Command.AddCommand(sprintCmd)
	Command.AddCommand(boardCmd)
	Command.AddCommand(projectCmd)
	Command.AddCommand(releaseCmd)
	Command.AddCommand(fieldCmd)
	Command.AddCommand(statusCmd)
	Command.AddCommand(issueTypeListCmd)
}

// resolveJiraClient resolves the Jira client from the active profile.
func resolveJiraClient(cmd *cobra.Command) (*jirasvc.Client, error) {
	_, p, err := cmdutil.ResolveProfile(cmd)
	if err != nil {
		return nil, err
	}

	conn, err := cmdutil.FindServiceByTypeOrName(p, config.ServiceTypeJira, serviceName)
	if err != nil {
		return nil, err
	}

	svc, err := service.Create(*conn)
	if err != nil {
		return nil, err
	}

	return jirasvc.ClientFromService(svc)
}

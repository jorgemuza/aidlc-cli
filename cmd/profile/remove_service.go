package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var removeSvcOpts struct {
	profileName string
}

var removeServiceCmd = &cobra.Command{
	Use:   "remove-service [service-name]",
	Short: "Remove a service connection from a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		p, err := cmdutil.ResolveProfileWithOverride(cmd, cfg, removeSvcOpts.profileName)
		if err != nil {
			return err
		}

		if err := p.RemoveService(args[0]); err != nil {
			return err
		}

		if err := cmdutil.SaveConfig(cmd, cfg); err != nil {
			return err
		}
		fmt.Printf("Service %q removed from profile %q.\n", args[0], p.Name)
		return nil
	},
}

func init() {
	removeServiceCmd.Flags().StringVar(&removeSvcOpts.profileName, "profile-name", "", "profile to remove the service from")
}

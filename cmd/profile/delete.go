package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a profile",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		if err := cfg.RemoveProfile(args[0]); err != nil {
			return err
		}

		if err := cmdutil.SaveConfig(cmd, cfg); err != nil {
			return err
		}
		fmt.Printf("Profile %q deleted.\n", args[0])
		return nil
	},
}

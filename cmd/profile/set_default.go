package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var setDefaultCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set a profile as the default",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		if err := cfg.SetDefault(args[0]); err != nil {
			return err
		}

		if err := cmdutil.SaveConfig(cmd, cfg); err != nil {
			return err
		}
		fmt.Printf("Profile %q is now the default.\n", args[0])
		return nil
	},
}

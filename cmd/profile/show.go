package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show profile details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		p := cfg.FindProfile(args[0])
		if p == nil {
			return fmt.Errorf("profile %q not found", args[0])
		}

		headers, rowFn := output.ServiceTable(p.Services)
		return output.Print(cmdutil.OutputFormat(cmd), p, headers, rowFn)
	},
}

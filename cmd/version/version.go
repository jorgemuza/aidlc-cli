package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ver    = "dev"
	commit = "none"
	date   = "unknown"
)

func Set(v, c, d string) {
	ver = v
	commit = c
	date = d
}

var Command = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("orbit %s (commit: %s, built: %s)\n", ver, commit, date)
	},
}

package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		headers := []string{"NAME", "DESCRIPTION", "DEFAULT", "SERVICES"}
		rows := func() [][]string {
			var r [][]string
			for _, p := range cfg.Profiles {
				def := ""
				if p.Default {
					def = "*"
				}
				svcTypes := map[string]int{}
				for _, s := range p.Services {
					svcTypes[s.Type]++
				}
				svcSummary := ""
				for t, c := range svcTypes {
					if svcSummary != "" {
						svcSummary += ", "
					}
					svcSummary += t
					if c > 1 {
						svcSummary += fmt.Sprintf("(%d)", c)
					}
				}
				r = append(r, []string{p.Name, p.Description, def, svcSummary})
			}
			return r
		}

		return output.Print(cmdutil.OutputFormat(cmd), cfg.Profiles, headers, rows)
	},
}

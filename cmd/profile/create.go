package profile

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/spf13/cobra"
)

var createOpts struct {
	name        string
	description string
	setDefault  bool
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "Create a new profile",
	Example: `  orbit profile create --name project-a --description "Project A services" --default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		p := config.Profile{
			Name:        createOpts.name,
			Description: createOpts.description,
			Default:     createOpts.setDefault,
		}
		if err := cfg.AddProfile(p); err != nil {
			return err
		}

		if err := cmdutil.SaveConfig(cmd, cfg); err != nil {
			return err
		}
		fmt.Printf("Profile %q created.\n", p.Name)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createOpts.name, "name", "n", "", "profile name (required)")
	createCmd.Flags().StringVarP(&createOpts.description, "description", "d", "", "profile description")
	createCmd.Flags().BoolVar(&createOpts.setDefault, "default", false, "set as default profile")
	_ = createCmd.MarkFlagRequired("name")
}

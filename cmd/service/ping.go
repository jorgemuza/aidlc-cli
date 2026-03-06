package service

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping [service-name]",
	Short: "Test connectivity to a service (or all services in the active profile)",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, p, err := cmdutil.ResolveProfile(cmd)
		if err != nil {
			return err
		}

		services := p.Services
		if len(args) > 0 {
			svc := p.FindService(args[0])
			if svc == nil {
				return fmt.Errorf("service %q not found in profile %q", args[0], p.Name)
			}
			services = []config.ServiceConnection{*svc}
		}

		if len(services) == 0 {
			fmt.Println("No services configured in this profile.")
			return nil
		}

		for _, conn := range services {
			svc, err := service.Create(conn)
			if err != nil {
				fmt.Printf("  %-20s FAIL  %s\n", conn.Name, err)
				continue
			}
			info, err := svc.Ping()
			if err != nil {
				fmt.Printf("  %-20s FAIL  %s\n", conn.Name, err)
			} else {
				fmt.Printf("  %-20s OK    %s\n", conn.Name, info)
			}
		}
		return nil
	},
}

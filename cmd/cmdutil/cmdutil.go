package cmdutil

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/output"
	"github.com/spf13/cobra"
)

// ConfigPath reads the --config flag from the root command.
func ConfigPath(cmd *cobra.Command) string {
	v, _ := cmd.Root().PersistentFlags().GetString("config")
	return v
}

// LoadConfig loads the config using the --config flag.
func LoadConfig(cmd *cobra.Command) (*config.Config, error) {
	return config.Load(ConfigPath(cmd))
}

// SaveConfig saves the config using the --config flag.
func SaveConfig(cmd *cobra.Command, cfg *config.Config) error {
	return config.Save(cfg, ConfigPath(cmd))
}

// ResolveProfile loads config and resolves the active profile from --profile flag.
func ResolveProfile(cmd *cobra.Command) (*config.Config, *config.Profile, error) {
	cfg, err := LoadConfig(cmd)
	if err != nil {
		return nil, nil, err
	}
	profileName, _ := cmd.Root().PersistentFlags().GetString("profile")
	p, err := cfg.ResolveProfile(profileName)
	return cfg, p, err
}

// ResolveProfileWithOverride is like ResolveProfile but checks a local override first.
func ResolveProfileWithOverride(cmd *cobra.Command, cfg *config.Config, localOverride string) (*config.Profile, error) {
	name := localOverride
	if name == "" {
		name, _ = cmd.Root().PersistentFlags().GetString("profile")
	}
	return cfg.ResolveProfile(name)
}

// OutputFormat reads the --output flag and parses it.
func OutputFormat(cmd *cobra.Command) output.Format {
	return output.FormatFromCmd(cmd)
}

// FindServiceByType finds the first service of the given type in a profile.
func FindServiceByType(p *config.Profile, serviceType string) (*config.ServiceConnection, error) {
	for i := range p.Services {
		if p.Services[i].Type == serviceType {
			return &p.Services[i], nil
		}
	}
	return nil, fmt.Errorf("no %s service found in profile %q", serviceType, p.Name)
}

// FindServiceByTypeOrName finds a service by explicit name or falls back to the first of the given type.
func FindServiceByTypeOrName(p *config.Profile, serviceType, serviceName string) (*config.ServiceConnection, error) {
	if serviceName != "" {
		s := p.FindService(serviceName)
		if s == nil {
			return nil, fmt.Errorf("service %q not found in profile %q", serviceName, p.Name)
		}
		return s, nil
	}
	return FindServiceByType(p, serviceType)
}

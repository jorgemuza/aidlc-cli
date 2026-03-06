package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Service type constants.
const (
	ServiceTypeJira       = "jira"
	ServiceTypeConfluence = "confluence"
	ServiceTypeGitLab     = "gitlab"
	ServiceTypeBitbucket  = "bitbucket"
	ServiceTypeGitHub     = "github"
)

// Variant constants.
const (
	VariantCloud  = "cloud"
	VariantServer = "server"
)

// Auth method constants.
const (
	AuthMethodToken  = "token"
	AuthMethodBasic  = "basic"
	AuthMethodOAuth2 = "oauth2"
)

var supportedServiceTypes = []string{ServiceTypeJira, ServiceTypeConfluence, ServiceTypeGitLab, ServiceTypeBitbucket, ServiceTypeGitHub}

// SupportedServiceTypes returns the list of service types the CLI supports.
func SupportedServiceTypes() []string {
	return supportedServiceTypes
}

// Config is the top-level configuration containing all profiles.
type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

// ServiceConnection defines a connection to a single service instance.
type ServiceConnection struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"`
	Variant string            `yaml:"variant"`
	BaseURL string            `yaml:"base_url"`
	Auth    AuthConfig        `yaml:"auth"`
	Options map[string]string `yaml:"options,omitempty"`
}

// AuthConfig holds authentication credentials for a service.
// Values can be plain text or 1Password secret references (op://vault/item/field).
type AuthConfig struct {
	Method       string `yaml:"method"`
	Token        string `yaml:"token,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
}

// Profile groups service connections under a named profile.
type Profile struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description,omitempty"`
	Default     bool                `yaml:"default,omitempty"`
	Services    []ServiceConnection `yaml:"services"`
}

// configFilePath returns the default config file path (~/.config/orbit/config.yaml).
// We always use ~/.config to be consistent across platforms and match the documented path.
func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "orbit", "config.yaml"), nil
}

// Load reads the config from the default path or the given override path.
func Load(path string) (*Config, error) {
	if path == "" {
		p, err := configFilePath()
		if err != nil {
			return nil, err
		}
		path = p
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to the default path or the given override path.
func Save(cfg *Config, path string) error {
	if path == "" {
		p, err := configFilePath()
		if err != nil {
			return err
		}
		path = p
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// FindProfile finds a profile by name. Returns nil if not found.
func (c *Config) FindProfile(name string) *Profile {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i]
		}
	}
	return nil
}

// DefaultProfile returns the profile marked as default, or the first profile if none is default.
func (c *Config) DefaultProfile() *Profile {
	for i := range c.Profiles {
		if c.Profiles[i].Default {
			return &c.Profiles[i]
		}
	}
	if len(c.Profiles) > 0 {
		return &c.Profiles[0]
	}
	return nil
}

// AddProfile adds a profile. Returns error if a profile with the same name exists.
func (c *Config) AddProfile(p Profile) error {
	if c.FindProfile(p.Name) != nil {
		return fmt.Errorf("profile %q already exists", p.Name)
	}
	if p.Default {
		c.clearDefaults()
	}
	c.Profiles = append(c.Profiles, p)
	return nil
}

// RemoveProfile removes a profile by name.
func (c *Config) RemoveProfile(name string) error {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			c.Profiles = append(c.Profiles[:i], c.Profiles[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("profile %q not found", name)
}

func (c *Config) clearDefaults() {
	for i := range c.Profiles {
		c.Profiles[i].Default = false
	}
}

// FindService finds a service connection by name within a profile.
func (p *Profile) FindService(name string) *ServiceConnection {
	for i := range p.Services {
		if p.Services[i].Name == name {
			return &p.Services[i]
		}
	}
	return nil
}

// AddService adds a service connection to the profile.
func (p *Profile) AddService(svc ServiceConnection) error {
	if p.FindService(svc.Name) != nil {
		return fmt.Errorf("service %q already exists in profile %q", svc.Name, p.Name)
	}
	p.Services = append(p.Services, svc)
	return nil
}

// RemoveService removes a service connection by name from the profile.
func (p *Profile) RemoveService(name string) error {
	for i := range p.Services {
		if p.Services[i].Name == name {
			p.Services = append(p.Services[:i], p.Services[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("service %q not found in profile %q", name, p.Name)
}

// SetDefault sets this profile as default and clears the default flag on the given config.
func (c *Config) SetDefault(name string) error {
	found := false
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			c.Profiles[i].Default = true
			found = true
		} else {
			c.Profiles[i].Default = false
		}
	}
	if !found {
		return fmt.Errorf("profile %q not found", name)
	}
	return nil
}

// ResolveProfile finds a profile by name, or returns the default if name is empty.
func (c *Config) ResolveProfile(name string) (*Profile, error) {
	if name != "" {
		p := c.FindProfile(name)
		if p == nil {
			return nil, fmt.Errorf("profile %q not found", name)
		}
		return p, nil
	}
	p := c.DefaultProfile()
	if p == nil {
		return nil, fmt.Errorf("no profiles configured; run 'orbit profile create' first")
	}
	return p, nil
}

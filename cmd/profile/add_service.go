package profile

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/config"
	"github.com/spf13/cobra"
)

var addSvcOpts struct {
	profileName  string
	name         string
	serviceType  string
	variant      string
	baseURL      string
	authMethod   string
	token        string
	username     string
	password     string
	clientID     string
	clientSecret string
}

var addServiceCmd = &cobra.Command{
	Use:   "add-service",
	Short: "Add a service connection to a profile",
	Long: `Add a service connection to a profile.

Auth credentials can be plain text or 1Password secret references:
  --token "op://DevVault/gitlab-token/credential"
  --password "op://DevVault/jira/password"

Supported types: jira, confluence, gitlab, bitbucket
Supported variants: cloud, server`,
	Example: `  # Add Jira Cloud with 1Password secret
  orbit profile add-service --profile project-a \
    --name jira-cloud --type jira --variant cloud \
    --base-url https://myco.atlassian.net \
    --auth-method token --token "op://Dev/jira-token/credential"

  # Add self-hosted GitLab with basic auth
  orbit profile add-service --profile project-a \
    --name gitlab-onprem --type gitlab --variant server \
    --base-url https://gitlab.internal.com \
    --auth-method basic --username admin --password "op://Dev/gitlab/password"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		p, err := cmdutil.ResolveProfileWithOverride(cmd, cfg, addSvcOpts.profileName)
		if err != nil {
			return err
		}

		if !slices.Contains(config.SupportedServiceTypes(), addSvcOpts.serviceType) {
			return fmt.Errorf("unsupported service type %q; supported: %s",
				addSvcOpts.serviceType, strings.Join(config.SupportedServiceTypes(), ", "))
		}

		svc := config.ServiceConnection{
			Name:    addSvcOpts.name,
			Type:    addSvcOpts.serviceType,
			Variant: addSvcOpts.variant,
			BaseURL: addSvcOpts.baseURL,
			Auth: config.AuthConfig{
				Method:       addSvcOpts.authMethod,
				Token:        addSvcOpts.token,
				Username:     addSvcOpts.username,
				Password:     addSvcOpts.password,
				ClientID:     addSvcOpts.clientID,
				ClientSecret: addSvcOpts.clientSecret,
			},
		}

		if err := p.AddService(svc); err != nil {
			return err
		}

		if err := cmdutil.SaveConfig(cmd, cfg); err != nil {
			return err
		}
		fmt.Printf("Service %q added to profile %q.\n", svc.Name, p.Name)
		return nil
	},
}

func init() {
	addServiceCmd.Flags().StringVar(&addSvcOpts.profileName, "profile-name", "", "profile to add the service to")
	addServiceCmd.Flags().StringVarP(&addSvcOpts.name, "name", "n", "", "service connection name (required)")
	addServiceCmd.Flags().StringVarP(&addSvcOpts.serviceType, "type", "t", "", "service type: jira, confluence, gitlab, bitbucket (required)")
	addServiceCmd.Flags().StringVar(&addSvcOpts.variant, "variant", "cloud", "variant: cloud, server")
	addServiceCmd.Flags().StringVar(&addSvcOpts.baseURL, "base-url", "", "base URL of the service")
	addServiceCmd.Flags().StringVar(&addSvcOpts.authMethod, "auth-method", "token", "auth method: token, basic, oauth2")
	addServiceCmd.Flags().StringVar(&addSvcOpts.token, "token", "", "API token or PAT (supports op:// references)")
	addServiceCmd.Flags().StringVar(&addSvcOpts.username, "username", "", "username for basic auth")
	addServiceCmd.Flags().StringVar(&addSvcOpts.password, "password", "", "password for basic auth (supports op:// references)")
	addServiceCmd.Flags().StringVar(&addSvcOpts.clientID, "client-id", "", "OAuth2 client ID")
	addServiceCmd.Flags().StringVar(&addSvcOpts.clientSecret, "client-secret", "", "OAuth2 client secret (supports op:// references)")
	_ = addServiceCmd.MarkFlagRequired("name")
	_ = addServiceCmd.MarkFlagRequired("type")
}

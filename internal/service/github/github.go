package github

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

func init() {
	service.Register(config.ServiceTypeGitHub, newService)
}

type svc struct{ service.BaseService }

func newService(conn config.ServiceConnection) (service.Service, error) {
	if conn.BaseURL == "" {
		if conn.Variant == config.VariantCloud || conn.Variant == "" {
			conn.BaseURL = "https://api.github.com"
		} else {
			return nil, fmt.Errorf("github: base_url is required for GitHub Enterprise")
		}
	}
	return &svc{service.NewBaseService(conn)}, nil
}

func (s *svc) Type() string { return config.ServiceTypeGitHub }

func (s *svc) Ping() (string, error) {
	var user struct {
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := s.DoGet("/user", &user); err != nil {
		return "", fmt.Errorf("github: %w", err)
	}
	return fmt.Sprintf("GitHub authenticated as %s (%s)", user.Login, user.Name), nil
}

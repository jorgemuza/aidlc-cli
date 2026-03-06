package gitlab

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

func init() {
	service.Register(config.ServiceTypeGitLab, newService)
}

type svc struct{ service.BaseService }

func newService(conn config.ServiceConnection) (service.Service, error) {
	if conn.BaseURL == "" {
		if conn.Variant == config.VariantCloud {
			conn.BaseURL = "https://gitlab.com"
		} else {
			return nil, fmt.Errorf("gitlab: base_url is required for self-hosted instances")
		}
	}
	return &svc{service.NewBaseService(conn)}, nil
}

func (s *svc) Type() string { return config.ServiceTypeGitLab }

func (s *svc) Ping() (string, error) {
	var info struct {
		Version  string `json:"version"`
		Revision string `json:"revision"`
	}
	if err := s.DoGet("/api/v4/version", &info); err != nil {
		return "", fmt.Errorf("gitlab: %w", err)
	}
	return fmt.Sprintf("GitLab %s (rev %s)", info.Version, info.Revision), nil
}

package confluence

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

func init() {
	service.Register(config.ServiceTypeConfluence, newService)
}

type svc struct{ service.BaseService }

func newService(conn config.ServiceConnection) (service.Service, error) {
	if conn.BaseURL == "" {
		return nil, fmt.Errorf("confluence: base_url is required")
	}
	return &svc{service.NewBaseService(conn)}, nil
}

func (s *svc) Type() string { return config.ServiceTypeConfluence }

func (s *svc) Ping() (string, error) {
	path := "/rest/api/space?limit=1"
	if s.Conn.Variant == config.VariantCloud {
		path = "/wiki/rest/api/space?limit=1"
	}

	var result struct {
		Results []struct {
			Key string `json:"key"`
		} `json:"results"`
	}
	if err := s.DoGet(path, &result); err != nil {
		return "", fmt.Errorf("confluence: %w", err)
	}
	return fmt.Sprintf("Confluence (%s) - OK, %d spaces accessible", s.Conn.Variant, len(result.Results)), nil
}

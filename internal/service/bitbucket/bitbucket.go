package bitbucket

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

func init() {
	service.Register(config.ServiceTypeBitbucket, newService)
}

type svc struct{ service.BaseService }

func newService(conn config.ServiceConnection) (service.Service, error) {
	if conn.BaseURL == "" {
		if conn.Variant == config.VariantCloud {
			conn.BaseURL = "https://api.bitbucket.org/2.0"
		} else {
			return nil, fmt.Errorf("bitbucket: base_url is required for self-hosted instances")
		}
	}
	return &svc{service.NewBaseService(conn)}, nil
}

func (s *svc) Type() string { return config.ServiceTypeBitbucket }

func (s *svc) Ping() (string, error) {
	path := "/user"
	if s.Conn.Variant == config.VariantServer {
		path = "/rest/api/latest/application-properties"
	}
	if err := s.DoGet(path, nil); err != nil {
		return "", fmt.Errorf("bitbucket: %w", err)
	}
	return fmt.Sprintf("Bitbucket (%s) - OK", s.Conn.Variant), nil
}

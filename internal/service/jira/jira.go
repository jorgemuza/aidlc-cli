package jira

import (
	"fmt"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

func init() {
	service.Register(config.ServiceTypeJira, newService)
}

type svc struct {
	*Client
}

func newService(conn config.ServiceConnection) (service.Service, error) {
	if conn.BaseURL == "" {
		return nil, fmt.Errorf("jira: base_url is required")
	}
	return &svc{NewClient(service.NewBaseService(conn))}, nil
}

func (s *svc) Type() string { return config.ServiceTypeJira }

func (s *svc) Ping() (string, error) {
	var info struct {
		Version     string `json:"version"`
		ServerTitle string `json:"serverTitle"`
	}
	if err := s.DoGet("/rest/api/2/serverInfo", &info); err != nil {
		return "", fmt.Errorf("jira: %w", err)
	}
	return fmt.Sprintf("Jira %s (%s)", info.Version, info.ServerTitle), nil
}

// ClientFromService extracts the Jira Client from a Service interface.
func ClientFromService(s service.Service) (*Client, error) {
	js, ok := s.(*svc)
	if !ok {
		return nil, fmt.Errorf("service is not a Jira service")
	}
	return js.Client, nil
}

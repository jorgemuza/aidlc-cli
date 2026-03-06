package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jorgemuza/aidlc-cli/internal/config"
	"github.com/jorgemuza/aidlc-cli/internal/secrets"
)

// Service is the interface that all service integrations must implement.
type Service interface {
	Type() string
	Ping() (string, error)
}

// Factory is a function that creates a Service from a ServiceConnection.
type Factory func(conn config.ServiceConnection) (Service, error)

var registry = map[string]Factory{}

// Register registers a service factory for the given service type.
func Register(serviceType string, factory Factory) {
	registry[serviceType] = factory
}

// Create creates a service instance from a ServiceConnection, resolving any secrets.
func Create(conn config.ServiceConnection) (Service, error) {
	factory, ok := registry[conn.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported service type: %q", conn.Type)
	}

	resolved, err := resolveAuth(conn.Auth)
	if err != nil {
		return nil, fmt.Errorf("resolving auth for %q: %w", conn.Name, err)
	}
	conn.Auth = resolved

	return factory(conn)
}

func resolveAuth(auth config.AuthConfig) (config.AuthConfig, error) {
	var err error
	if auth.Username, err = secrets.Resolve(auth.Username); err != nil {
		return auth, err
	}
	if auth.Token, err = secrets.Resolve(auth.Token); err != nil {
		return auth, err
	}
	if auth.Password, err = secrets.Resolve(auth.Password); err != nil {
		return auth, err
	}
	if auth.ClientSecret, err = secrets.Resolve(auth.ClientSecret); err != nil {
		return auth, err
	}
	return auth, nil
}

// BaseService provides common fields for all service implementations.
type BaseService struct {
	Conn   config.ServiceConnection
	Client *http.Client
}

// NewBaseService creates a BaseService with an authenticated HTTP client.
func NewBaseService(conn config.ServiceConnection) BaseService {
	return BaseService{
		Conn:   conn,
		Client: newHTTPClient(conn.Auth),
	}
}

func newHTTPClient(auth config.AuthConfig) *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &authTransport{
			auth: auth,
			base: http.DefaultTransport,
		},
	}
}

type authTransport struct {
	auth config.AuthConfig
	base http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch t.auth.Method {
	case config.AuthMethodToken:
		req.Header.Set("Authorization", "Bearer "+t.auth.Token)
	case config.AuthMethodBasic:
		req.SetBasicAuth(t.auth.Username, t.auth.Password)
	}
	return t.base.RoundTrip(req)
}

// DoGet performs a GET request, checks for 200, and JSON-decodes into target.
// If target is nil, the body is drained and discarded.
func (b *BaseService) DoGet(path string, target any) error {
	url := strings.TrimRight(b.Conn.BaseURL, "/") + path
	resp, err := b.Client.Get(url)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}
	return nil
}

// DoRequest performs an HTTP request with a JSON body and decodes the response.
func (b *BaseService) DoRequest(method, path string, body any, target any) error {
	url := strings.TrimRight(b.Conn.BaseURL, "/") + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = strings.NewReader(string(data))
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if target != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}
	return nil
}

// DoPost performs a POST request with JSON body.
func (b *BaseService) DoPost(path string, body any, target any) error {
	return b.DoRequest(http.MethodPost, path, body, target)
}

// DoPut performs a PUT request with JSON body.
func (b *BaseService) DoPut(path string, body any, target any) error {
	return b.DoRequest(http.MethodPut, path, body, target)
}

// DoGetRaw performs a GET request and returns the response body as a string.
func (b *BaseService) DoGetRaw(path string) (string, error) {
	url := strings.TrimRight(b.Conn.BaseURL, "/") + path
	resp, err := b.Client.Get(url)
	if err != nil {
		return "", fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// DoDelete performs a DELETE request.
func (b *BaseService) DoDelete(path string) error {
	return b.DoRequest(http.MethodDelete, path, nil, nil)
}

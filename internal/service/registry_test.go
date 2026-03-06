package service

import (
	"testing"

	"github.com/jorgemuza/orbit/internal/config"
)

type mockService struct {
	svcType string
}

func (m *mockService) Type() string        { return m.svcType }
func (m *mockService) Ping() (string, error) { return "ok", nil }

func TestRegisterAndCreate(t *testing.T) {
	Register("mock", func(conn config.ServiceConnection) (Service, error) {
		return &mockService{svcType: conn.Type}, nil
	})
	defer delete(registry, "mock")

	svc, err := Create(config.ServiceConnection{
		Name: "test", Type: "mock",
		Auth: config.AuthConfig{Method: "token", Token: "plain-token"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Type() != "mock" {
		t.Fatalf("expected type 'mock', got %q", svc.Type())
	}
	info, err := svc.Ping()
	if err != nil {
		t.Fatalf("ping error: %v", err)
	}
	if info != "ok" {
		t.Fatalf("expected 'ok', got %q", info)
	}
}

func TestCreateUnsupportedType(t *testing.T) {
	_, err := Create(config.ServiceConnection{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestNewBaseService(t *testing.T) {
	conn := config.ServiceConnection{
		Name:    "test",
		BaseURL: "https://example.com",
		Auth:    config.AuthConfig{Method: "token", Token: "abc"},
	}
	bs := NewBaseService(conn)
	if bs.Client == nil {
		t.Fatal("expected HTTP client to be set")
	}
	if bs.Conn.Name != "test" {
		t.Fatalf("expected conn name 'test', got %q", bs.Conn.Name)
	}
}

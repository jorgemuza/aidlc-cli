package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNonExistentReturnsEmpty(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Profiles) != 0 {
		t.Fatalf("expected 0 profiles, got %d", len(cfg.Profiles))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	cfg := &Config{
		Profiles: []Profile{
			{Name: "test", Description: "desc", Default: true},
		},
	}
	if err := Save(cfg, path); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(loaded.Profiles))
	}
	if loaded.Profiles[0].Name != "test" {
		t.Fatalf("expected name 'test', got %q", loaded.Profiles[0].Name)
	}
	if !loaded.Profiles[0].Default {
		t.Fatal("expected default to be true")
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "dir", "config.yaml")
	cfg := &Config{}
	if err := Save(cfg, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestAddProfile(t *testing.T) {
	cfg := &Config{}
	if err := cfg.AddProfile(Profile{Name: "a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(cfg.Profiles))
	}
}

func TestAddProfileDuplicate(t *testing.T) {
	cfg := &Config{}
	_ = cfg.AddProfile(Profile{Name: "a"})
	err := cfg.AddProfile(Profile{Name: "a"})
	if err == nil {
		t.Fatal("expected error for duplicate profile")
	}
}

func TestAddProfileDefaultClearsOthers(t *testing.T) {
	cfg := &Config{}
	_ = cfg.AddProfile(Profile{Name: "a", Default: true})
	_ = cfg.AddProfile(Profile{Name: "b", Default: true})
	if cfg.Profiles[0].Default {
		t.Fatal("expected first profile default to be cleared")
	}
	if !cfg.Profiles[1].Default {
		t.Fatal("expected second profile to be default")
	}
}

func TestRemoveProfile(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}, {Name: "b"}}}
	if err := cfg.RemoveProfile("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(cfg.Profiles))
	}
	if cfg.Profiles[0].Name != "b" {
		t.Fatalf("expected remaining profile 'b', got %q", cfg.Profiles[0].Name)
	}
}

func TestRemoveProfileNotFound(t *testing.T) {
	cfg := &Config{}
	if err := cfg.RemoveProfile("missing"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestFindProfile(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}, {Name: "b"}}}
	p := cfg.FindProfile("b")
	if p == nil || p.Name != "b" {
		t.Fatal("expected to find profile 'b'")
	}
	if cfg.FindProfile("missing") != nil {
		t.Fatal("expected nil for missing profile")
	}
}

func TestDefaultProfile(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}, {Name: "b", Default: true}}}
	p := cfg.DefaultProfile()
	if p == nil || p.Name != "b" {
		t.Fatal("expected default profile 'b'")
	}
}

func TestDefaultProfileFallsBackToFirst(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}, {Name: "b"}}}
	p := cfg.DefaultProfile()
	if p == nil || p.Name != "a" {
		t.Fatal("expected fallback to first profile 'a'")
	}
}

func TestDefaultProfileEmpty(t *testing.T) {
	cfg := &Config{}
	if cfg.DefaultProfile() != nil {
		t.Fatal("expected nil for empty config")
	}
}

func TestSetDefault(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a", Default: true}, {Name: "b"}}}
	if err := cfg.SetDefault("b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Profiles[0].Default {
		t.Fatal("expected 'a' default to be cleared")
	}
	if !cfg.Profiles[1].Default {
		t.Fatal("expected 'b' to be default")
	}
}

func TestSetDefaultNotFound(t *testing.T) {
	cfg := &Config{}
	if err := cfg.SetDefault("missing"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestResolveProfileByName(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}, {Name: "b"}}}
	p, err := cfg.ResolveProfile("b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "b" {
		t.Fatalf("expected 'b', got %q", p.Name)
	}
}

func TestResolveProfileDefault(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a", Default: true}}}
	p, err := cfg.ResolveProfile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "a" {
		t.Fatalf("expected 'a', got %q", p.Name)
	}
}

func TestResolveProfileNoProfiles(t *testing.T) {
	cfg := &Config{}
	_, err := cfg.ResolveProfile("")
	if err == nil {
		t.Fatal("expected error for no profiles")
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	cfg := &Config{Profiles: []Profile{{Name: "a"}}}
	_, err := cfg.ResolveProfile("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestAddService(t *testing.T) {
	p := &Profile{Name: "test"}
	svc := ServiceConnection{Name: "jira-1", Type: ServiceTypeJira}
	if err := p.AddService(svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(p.Services))
	}
}

func TestAddServiceDuplicate(t *testing.T) {
	p := &Profile{Name: "test"}
	svc := ServiceConnection{Name: "jira-1", Type: ServiceTypeJira}
	_ = p.AddService(svc)
	if err := p.AddService(svc); err == nil {
		t.Fatal("expected error for duplicate service")
	}
}

func TestFindService(t *testing.T) {
	p := &Profile{
		Name:     "test",
		Services: []ServiceConnection{{Name: "jira-1"}, {Name: "gitlab-1"}},
	}
	s := p.FindService("gitlab-1")
	if s == nil || s.Name != "gitlab-1" {
		t.Fatal("expected to find 'gitlab-1'")
	}
	if p.FindService("missing") != nil {
		t.Fatal("expected nil for missing service")
	}
}

func TestRemoveService(t *testing.T) {
	p := &Profile{
		Name:     "test",
		Services: []ServiceConnection{{Name: "a"}, {Name: "b"}},
	}
	if err := p.RemoveService("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Services) != 1 || p.Services[0].Name != "b" {
		t.Fatal("expected only 'b' to remain")
	}
}

func TestRemoveServiceNotFound(t *testing.T) {
	p := &Profile{Name: "test"}
	if err := p.RemoveService("missing"); err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestSupportedServiceTypes(t *testing.T) {
	types := SupportedServiceTypes()
	if len(types) != 5 {
		t.Fatalf("expected 5 service types, got %d", len(types))
	}
}

func TestSaveAndLoadWithServices(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	cfg := &Config{
		Profiles: []Profile{
			{
				Name:    "prod",
				Default: true,
				Services: []ServiceConnection{
					{
						Name:    "jira-cloud",
						Type:    ServiceTypeJira,
						Variant: VariantCloud,
						BaseURL: "https://myco.atlassian.net",
						Auth: AuthConfig{
							Method: AuthMethodToken,
							Token:  "op://Vault/jira/token",
						},
					},
				},
			},
		},
	}
	if err := Save(cfg, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	svc := loaded.Profiles[0].Services[0]
	if svc.Name != "jira-cloud" || svc.Type != ServiceTypeJira {
		t.Fatalf("unexpected service: %+v", svc)
	}
	if svc.Auth.Token != "op://Vault/jira/token" {
		t.Fatalf("unexpected token: %q", svc.Auth.Token)
	}
}

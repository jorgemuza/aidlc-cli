package output

import (
	"testing"

	"github.com/jorgemuza/orbit/internal/config"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input string
		want  Format
	}{
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"yaml", FormatYAML},
		{"YAML", FormatYAML},
		{"table", FormatTable},
		{"", FormatTable},
		{"unknown", FormatTable},
	}
	for _, tt := range tests {
		if got := ParseFormat(tt.input); got != tt.want {
			t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestServiceTableHeaders(t *testing.T) {
	headers, _ := ServiceTable(nil)
	if len(headers) != 6 {
		t.Fatalf("expected 6 headers, got %d", len(headers))
	}
}

func TestServiceTableRows(t *testing.T) {
	services := []config.ServiceConnection{
		{
			Name: "jira-1", Type: "jira", Variant: "cloud",
			BaseURL: "https://example.com",
			Auth:    config.AuthConfig{Method: "token", Token: "op://vault/item"},
		},
		{
			Name: "gitlab-1", Type: "gitlab", Variant: "server",
			BaseURL: "https://gitlab.local",
			Auth:    config.AuthConfig{Method: "basic", Username: "user", Password: "pass"},
		},
	}
	_, rowFn := ServiceTable(services)
	rows := rowFn()
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	// First row should detect 1Password reference
	if rows[0][5] != "yes" {
		t.Errorf("expected '1PASSWORD' column to be 'yes' for op:// token, got %q", rows[0][5])
	}
	// Second row has no op:// refs
	if rows[1][5] != "" {
		t.Errorf("expected '1PASSWORD' column to be empty for plain auth, got %q", rows[1][5])
	}
}

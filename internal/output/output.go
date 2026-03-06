package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/secrets"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Format represents the output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// ParseFormat parses a string into a Format, defaulting to table.
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatTable
	}
}

// FormatFromCmd reads the --output flag from the root command and parses it.
func FormatFromCmd(cmd *cobra.Command) Format {
	outFmt, _ := cmd.Root().PersistentFlags().GetString("output")
	return ParseFormat(outFmt)
}

// Print prints data in the given format.
func Print(format Format, data any, headers []string, rowFn func() [][]string) error {
	switch format {
	case FormatJSON:
		return printJSON(data)
	case FormatYAML:
		return printYAML(data)
	default:
		printTable(headers, rowFn())
		return nil
	}
}

// ServiceTable builds headers and a row function for displaying service connections.
func ServiceTable(services []config.ServiceConnection) ([]string, func() [][]string) {
	headers := []string{"NAME", "TYPE", "VARIANT", "BASE URL", "AUTH", "1PASSWORD"}
	rowFn := func() [][]string {
		var rows [][]string
		for _, svc := range services {
			opRef := ""
			if secrets.IsSecretReference(svc.Auth.Token) ||
				secrets.IsSecretReference(svc.Auth.Password) ||
				secrets.IsSecretReference(svc.Auth.ClientSecret) {
				opRef = "yes"
			}
			rows = append(rows, []string{svc.Name, svc.Type, svc.Variant, svc.BaseURL, svc.Auth.Method, opRef})
		}
		return rows
	}
	return headers, rowFn
}

func printTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	fmt.Fprintln(w, strings.Join(makeSeparators(headers), "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

func makeSeparators(headers []string) []string {
	seps := make([]string, len(headers))
	for i, h := range headers {
		seps[i] = strings.Repeat("-", len(h))
	}
	return seps
}

func printJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func printYAML(data any) error {
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	return enc.Encode(data)
}

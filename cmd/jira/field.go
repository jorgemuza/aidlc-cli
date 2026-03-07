package jira

import (
	"fmt"
	"strings"

	jirasvc "github.com/jorgemuza/orbit/internal/service/jira"
	"github.com/spf13/cobra"
)

var fieldCmd = &cobra.Command{
	Use:   "field",
	Short: "Manage Jira fields",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// field list — list all fields, optionally filtered
var fieldListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "List Jira fields (system and custom)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		fields, err := client.ListFields()
		if err != nil {
			return err
		}

		filter, _ := cmd.Flags().GetString("filter")
		filter = strings.ToLower(filter)
		customOnly, _ := cmd.Flags().GetBool("custom")

		for _, f := range fields {
			if customOnly && !f.Custom {
				continue
			}
			if filter != "" && !strings.Contains(strings.ToLower(f.Name), filter) && !strings.Contains(strings.ToLower(f.ID), filter) {
				continue
			}
			customLabel := ""
			if f.Custom {
				customLabel = " [custom]"
			}
			desc := ""
			if f.Description != "" {
				desc = "  — " + f.Description
			}
			fmt.Printf("%-30s %s%s%s\n", f.ID, f.Name, customLabel, desc)
		}
		return nil
	},
}

// field create — create a custom field (Cloud only)
var fieldCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom field (Jira Cloud only)",
	Long: `Create a custom Jira field. Supported types:
  com.atlassian.jira.plugin.system.customfieldtypes:select       — Single select
  com.atlassian.jira.plugin.system.customfieldtypes:multiselect  — Multi select
  com.atlassian.jira.plugin.system.customfieldtypes:float        — Number (float)
  com.atlassian.jira.plugin.system.customfieldtypes:multicheckboxes — Checkboxes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		fieldType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		searcherKey, _ := cmd.Flags().GetString("searcher")

		if name == "" || fieldType == "" {
			return fmt.Errorf("--name and --type are required")
		}

		// Resolve shorthand types
		fieldType = resolveFieldType(fieldType)
		if searcherKey == "" {
			searcherKey = resolveSearcherKey(fieldType)
		}

		req := &jirasvc.CreateFieldRequest{
			Name:        name,
			Description: description,
			Type:        fieldType,
			SearcherKey: searcherKey,
		}
		result, err := client.CreateField(req)
		if err != nil {
			return err
		}
		fmt.Printf("Created field: %s (ID: %s)\n", result.Name, result.ID)
		return nil
	},
}

// field context list — list contexts for a field
var fieldContextListCmd = &cobra.Command{
	Use:   "context-list [field-id]",
	Short: "List contexts for a custom field (Cloud only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		contexts, err := client.ListFieldContexts(args[0])
		if err != nil {
			return err
		}

		for _, ctx := range contexts {
			global := ""
			if ctx.IsGlobal {
				global = " [global]"
			}
			anyIssue := ""
			if ctx.IsAnyIssue {
				anyIssue = " [any-issue-type]"
			}
			fmt.Printf("%-10s %s%s%s\n", ctx.ID, ctx.Name, global, anyIssue)
		}
		return nil
	},
}

// field option list — list options for a select field context
var fieldOptionListCmd = &cobra.Command{
	Use:   "option-list [field-id] [context-id]",
	Short: "List options for a select/multi-select field context (Cloud only)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		options, err := client.ListFieldOptions(args[0], args[1])
		if err != nil {
			return err
		}

		for _, opt := range options {
			disabled := ""
			if opt.Disabled {
				disabled = " [disabled]"
			}
			fmt.Printf("%-10s %s%s\n", opt.ID, opt.Value, disabled)
		}
		return nil
	},
}

// field option add — add options to a select field context
var fieldOptionAddCmd = &cobra.Command{
	Use:   "option-add [field-id] [context-id]",
	Short: "Add options to a select/multi-select field context (Cloud only)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		values, _ := cmd.Flags().GetStringSlice("values")
		if len(values) == 0 {
			return fmt.Errorf("--values is required (comma-separated)")
		}

		options, err := client.AddFieldOptions(args[0], args[1], values)
		if err != nil {
			return err
		}

		for _, opt := range options {
			fmt.Printf("Added option: %s (ID: %s)\n", opt.Value, opt.ID)
		}
		return nil
	},
}

func init() {
	// field list flags
	fieldListSubCmd.Flags().String("filter", "", "filter fields by name or ID (case-insensitive)")
	fieldListSubCmd.Flags().Bool("custom", false, "show only custom fields")

	// field create flags
	fieldCreateCmd.Flags().String("name", "", "field name (required)")
	fieldCreateCmd.Flags().String("type", "", "field type or shorthand: select, multiselect, number, checkbox (required)")
	fieldCreateCmd.Flags().String("description", "", "field description")
	fieldCreateCmd.Flags().String("searcher", "", "searcher key (auto-resolved if omitted)")

	// field option add flags
	fieldOptionAddCmd.Flags().StringSlice("values", nil, "option values to add (comma-separated)")

	// register subcommands
	fieldCmd.AddCommand(fieldListSubCmd)
	fieldCmd.AddCommand(fieldCreateCmd)
	fieldCmd.AddCommand(fieldContextListCmd)
	fieldCmd.AddCommand(fieldOptionListCmd)
	fieldCmd.AddCommand(fieldOptionAddCmd)
}

// resolveFieldType maps shorthand names to full Jira custom field type keys.
func resolveFieldType(t string) string {
	prefix := "com.atlassian.jira.plugin.system.customfieldtypes:"
	shorthands := map[string]string{
		"select":      prefix + "select",
		"multiselect": prefix + "multiselect",
		"number":      prefix + "float",
		"float":       prefix + "float",
		"checkbox":    prefix + "multicheckboxes",
		"checkboxes":  prefix + "multicheckboxes",
		"text":        prefix + "textfield",
		"textarea":    prefix + "textarea",
	}
	if full, ok := shorthands[strings.ToLower(t)]; ok {
		return full
	}
	return t
}

// resolveSearcherKey returns the appropriate searcher key for a field type.
func resolveSearcherKey(fieldType string) string {
	prefix := "com.atlassian.jira.plugin.system.customfieldtypes:"
	searchers := map[string]string{
		prefix + "select":          prefix + "multiselectsearcher",
		prefix + "multiselect":     prefix + "multiselectsearcher",
		prefix + "float":           prefix + "exactnumber",
		prefix + "multicheckboxes": prefix + "multiselectsearcher",
		prefix + "textfield":       prefix + "textsearcher",
		prefix + "textarea":        prefix + "textsearcher",
	}
	if s, ok := searchers[fieldType]; ok {
		return s
	}
	return ""
}

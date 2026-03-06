package gitlab

import (
	"encoding/json"
	"fmt"

	glsvc "github.com/jorgemuza/orbit/internal/service/gitlab"
	"github.com/spf13/cobra"
)

var variableCmd = &cobra.Command{
	Use:     "variable [subcommand]",
	Short:   "Manage CI/CD variables",
	Aliases: []string{"var"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var variableListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List CI/CD variables",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab variable list foundation/ai
  orbit gitlab variable list 650`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		vars, err := client.ListVariables(args[0], limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(vars, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-30s %-10s %-10s %-10s %s\n", "KEY", "PROTECTED", "MASKED", "TYPE", "SCOPE")
		fmt.Printf("%-30s %-10s %-10s %-10s %s\n", "---", "---------", "------", "----", "-----")
		for _, v := range vars {
			fmt.Printf("%-30s %-10v %-10v %-10s %s\n", v.Key, v.Protected, v.Masked, v.VariableType, v.EnvironmentScope)
		}
		return nil
	},
}

var variableGetCmd = &cobra.Command{
	Use:   "get [project] [key]",
	Short: "Get a CI/CD variable",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab variable get foundation/ai CONFLUENCE_USERNAME`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		v, err := client.GetVariable(args[0], args[1])
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(v, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Key:       %s\n", v.Key)
		fmt.Printf("Value:     %s\n", v.Value)
		fmt.Printf("Type:      %s\n", v.VariableType)
		fmt.Printf("Protected: %v\n", v.Protected)
		fmt.Printf("Masked:    %v\n", v.Masked)
		fmt.Printf("Scope:     %s\n", v.EnvironmentScope)
		return nil
	},
}

var variableSetCmd = &cobra.Command{
	Use:   "set [project] [key] [value]",
	Short: "Create or update a CI/CD variable",
	Args:  cobra.ExactArgs(3),
	Example: `  orbit gitlab variable set foundation/ai MY_VAR "my-value"
  orbit gitlab variable set foundation/ai MY_VAR "secret" --masked --protected`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		protected, _ := cmd.Flags().GetBool("protected")
		masked, _ := cmd.Flags().GetBool("masked")
		scope, _ := cmd.Flags().GetString("scope")

		v := glsvc.Variable{
			Key:              args[1],
			Value:            args[2],
			VariableType:     "env_var",
			Protected:        protected,
			Masked:           masked,
			EnvironmentScope: scope,
		}

		// Try update first; if it fails, create.
		result, err := client.UpdateVariable(args[0], v)
		if err != nil {
			result, err = client.CreateVariable(args[0], v)
			if err != nil {
				return err
			}
			fmt.Printf("Created variable %s (protected=%v, masked=%v)\n", result.Key, result.Protected, result.Masked)
			return nil
		}

		fmt.Printf("Updated variable %s (protected=%v, masked=%v)\n", result.Key, result.Protected, result.Masked)
		return nil
	},
}

var variableDeleteCmd = &cobra.Command{
	Use:   "delete [project] [key]",
	Short: "Delete a CI/CD variable",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab variable delete foundation/ai MY_VAR`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		if err := client.DeleteVariable(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Deleted variable %s\n", args[1])
		return nil
	},
}

func init() {
	variableCmd.AddCommand(variableListCmd)
	variableCmd.AddCommand(variableGetCmd)
	variableCmd.AddCommand(variableSetCmd)
	variableCmd.AddCommand(variableDeleteCmd)

	variableListCmd.Flags().Int("limit", 50, "max results")

	variableSetCmd.Flags().Bool("protected", false, "only expose to protected branches")
	variableSetCmd.Flags().Bool("masked", false, "mask variable in job logs")
	variableSetCmd.Flags().String("scope", "*", "environment scope")
}

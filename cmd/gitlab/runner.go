package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var runnerCmd = &cobra.Command{
	Use:   "runner [subcommand]",
	Short: "Manage CI/CD runners",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var runnerListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List runners assigned to a project",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab runner list foundation/ai
  orbit gitlab runner list foundation/ai --status online
  orbit gitlab runner list foundation/ai --tag docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		status, _ := cmd.Flags().GetString("status")
		tag, _ := cmd.Flags().GetString("tag")
		limit, _ := cmd.Flags().GetInt("limit")

		runners, err := client.ListProjectRunners(args[0], status, tag, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(runners, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(runners) == 0 {
			fmt.Println("No runners assigned to this project.")
			fmt.Println("Use 'orbit gitlab runner list-all' to see available runners,")
			fmt.Println("then 'orbit gitlab runner enable <project> <runner-id>' to assign one.")
			return nil
		}

		fmt.Printf("%-8s %-30s %-10s %-8s %-8s %s\n", "ID", "DESCRIPTION", "STATUS", "SHARED", "ACTIVE", "TAGS")
		fmt.Printf("%-8s %-30s %-10s %-8s %-8s %s\n", "--", "-----------", "------", "------", "------", "----")
		for _, r := range runners {
			desc := r.Description
			if len(desc) > 28 {
				desc = desc[:25] + "..."
			}
			tags := strings.Join(r.TagList, ", ")
			status := r.Status
			if r.Online {
				status = "online"
			}
			fmt.Printf("%-8d %-30s %-10s %-8v %-8v %s\n", r.ID, desc, status, r.IsShared, r.Active, tags)
		}
		return nil
	},
}

var runnerListAllCmd = &cobra.Command{
	Use:   "list-all",
	Short: "List all runners visible to you (admin)",
	Example: `  orbit gitlab runner list-all
  orbit gitlab runner list-all --scope shared --status online
  orbit gitlab runner list-all --tag docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		scope, _ := cmd.Flags().GetString("scope")
		status, _ := cmd.Flags().GetString("status")
		tag, _ := cmd.Flags().GetString("tag")
		limit, _ := cmd.Flags().GetInt("limit")

		runners, err := client.ListAllRunners(scope, status, tag, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(runners, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(runners) == 0 {
			fmt.Println("No runners found.")
			return nil
		}

		fmt.Printf("%-8s %-30s %-12s %-10s %-8s %-8s %s\n", "ID", "DESCRIPTION", "TYPE", "STATUS", "SHARED", "ACTIVE", "TAGS")
		fmt.Printf("%-8s %-30s %-12s %-10s %-8s %-8s %s\n", "--", "-----------", "----", "------", "------", "------", "----")
		for _, r := range runners {
			desc := r.Description
			if len(desc) > 28 {
				desc = desc[:25] + "..."
			}
			tags := strings.Join(r.TagList, ", ")
			status := r.Status
			if r.Online {
				status = "online"
			}
			fmt.Printf("%-8d %-30s %-12s %-10s %-8v %-8v %s\n", r.ID, desc, r.RunnerType, status, r.IsShared, r.Active, tags)
		}
		return nil
	},
}

var runnerViewCmd = &cobra.Command{
	Use:   "view [runner-id]",
	Short: "View runner details",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit gitlab runner view 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid runner ID: %s", args[0])
		}

		runner, err := client.GetRunner(id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(runner, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:          %d\n", runner.ID)
		fmt.Printf("Description: %s\n", runner.Description)
		fmt.Printf("Type:        %s\n", runner.RunnerType)
		fmt.Printf("Status:      %s\n", runner.Status)
		fmt.Printf("Active:      %v\n", runner.Active)
		fmt.Printf("Shared:      %v\n", runner.IsShared)
		fmt.Printf("Online:      %v\n", runner.Online)
		if runner.IPAddress != "" {
			fmt.Printf("IP:          %s\n", runner.IPAddress)
		}
		if len(runner.TagList) > 0 {
			fmt.Printf("Tags:        %s\n", strings.Join(runner.TagList, ", "))
		}
		return nil
	},
}

var runnerEnableCmd = &cobra.Command{
	Use:   "enable [project] [runner-id]",
	Short: "Assign a runner to a project",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab runner enable foundation/ai 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid runner ID: %s", args[1])
		}

		runner, err := client.EnableProjectRunner(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Enabled runner #%d '%s' for project %s\n", runner.ID, runner.Description, args[0])
		return nil
	},
}

var runnerDisableCmd = &cobra.Command{
	Use:   "disable [project] [runner-id]",
	Short: "Remove a runner from a project",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit gitlab runner disable foundation/ai 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid runner ID: %s", args[1])
		}

		if err := client.DisableProjectRunner(args[0], id); err != nil {
			return err
		}

		fmt.Printf("Disabled runner #%d for project %s\n", id, args[0])
		return nil
	},
}

func init() {
	runnerCmd.AddCommand(runnerListCmd)
	runnerCmd.AddCommand(runnerListAllCmd)
	runnerCmd.AddCommand(runnerViewCmd)
	runnerCmd.AddCommand(runnerEnableCmd)
	runnerCmd.AddCommand(runnerDisableCmd)

	runnerListCmd.Flags().String("status", "", "filter: online, offline, stale, never_contacted")
	runnerListCmd.Flags().String("tag", "", "filter by tag")
	runnerListCmd.Flags().Int("limit", 20, "max results")

	runnerListAllCmd.Flags().String("scope", "", "scope: shared, specific, group_type, project_type")
	runnerListAllCmd.Flags().String("status", "", "filter: online, offline, stale, never_contacted")
	runnerListAllCmd.Flags().String("tag", "", "filter by tag")
	runnerListAllCmd.Flags().Int("limit", 20, "max results")
}

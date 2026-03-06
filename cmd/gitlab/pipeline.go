package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var pipelineCmd = &cobra.Command{
	Use:     "pipeline [subcommand]",
	Short:   "Manage pipelines",
	Aliases: []string{"pipe", "ci"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var pipelineListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List pipelines",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab pipeline list 595
  aidlc gitlab pipeline list 595 --ref main --status success`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		ref, _ := cmd.Flags().GetString("ref")
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		pipelines, err := client.ListPipelines(args[0], ref, status, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(pipelines, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-12s %-20s %-10s %s\n", "ID", "STATUS", "REF", "SOURCE", "CREATED")
		fmt.Printf("%-10s %-12s %-20s %-10s %s\n", "--", "------", "---", "------", "-------")
		for _, p := range pipelines {
			created := ""
			if len(p.CreatedAt) >= 10 {
				created = p.CreatedAt[:10]
			}
			ref := p.Ref
			if len(ref) > 18 {
				ref = ref[:15] + "..."
			}
			fmt.Printf("%-10d %-12s %-20s %-10s %s\n", p.ID, p.Status, ref, p.Source, created)
		}
		return nil
	},
}

var pipelineViewCmd = &cobra.Command{
	Use:   "view [project] [pipeline-id]",
	Short: "View a pipeline",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab pipeline view 595 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid pipeline ID: %s", args[1])
		}

		pipeline, err := client.GetPipeline(args[0], id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(pipeline, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:      %d\n", pipeline.ID)
		fmt.Printf("Status:  %s\n", pipeline.Status)
		fmt.Printf("Ref:     %s\n", pipeline.Ref)
		fmt.Printf("SHA:     %s\n", pipeline.SHA)
		fmt.Printf("Source:  %s\n", pipeline.Source)
		fmt.Printf("Created: %s\n", pipeline.CreatedAt)
		fmt.Printf("URL:     %s\n", pipeline.WebURL)
		return nil
	},
}

var pipelineJobsCmd = &cobra.Command{
	Use:   "jobs [project] [pipeline-id]",
	Short: "List jobs in a pipeline",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab pipeline jobs 595 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid pipeline ID: %s", args[1])
		}

		limit, _ := cmd.Flags().GetInt("limit")
		jobs, err := client.ListPipelineJobs(args[0], id, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(jobs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-20s %-15s %-12s %s\n", "ID", "NAME", "STAGE", "STATUS", "DURATION")
		fmt.Printf("%-10s %-20s %-15s %-12s %s\n", "--", "----", "-----", "------", "--------")
		for _, j := range jobs {
			dur := ""
			if j.Duration > 0 {
				dur = fmt.Sprintf("%.0fs", j.Duration)
			}
			fmt.Printf("%-10d %-20s %-15s %-12s %s\n", j.ID, j.Name, j.Stage, j.Status, dur)
		}
		return nil
	},
}

var pipelineRetryCmd = &cobra.Command{
	Use:   "retry [project] [pipeline-id]",
	Short: "Retry a pipeline",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab pipeline retry 595 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid pipeline ID: %s", args[1])
		}

		pipeline, err := client.RetryPipeline(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Retried pipeline %d (status: %s)\n", pipeline.ID, pipeline.Status)
		return nil
	},
}

var pipelineCancelCmd = &cobra.Command{
	Use:   "cancel [project] [pipeline-id]",
	Short: "Cancel a pipeline",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab pipeline cancel 595 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid pipeline ID: %s", args[1])
		}

		pipeline, err := client.CancelPipeline(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Canceled pipeline %d (status: %s)\n", pipeline.ID, pipeline.Status)
		return nil
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineListCmd)
	pipelineCmd.AddCommand(pipelineViewCmd)
	pipelineCmd.AddCommand(pipelineJobsCmd)
	pipelineCmd.AddCommand(pipelineRetryCmd)
	pipelineCmd.AddCommand(pipelineCancelCmd)

	pipelineListCmd.Flags().String("ref", "", "filter by branch/tag")
	pipelineListCmd.Flags().String("status", "", "filter by status: running, pending, success, failed, canceled")
	pipelineListCmd.Flags().Int("limit", 20, "max results")

	pipelineJobsCmd.Flags().Int("limit", 50, "max results")
}

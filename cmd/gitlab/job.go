package gitlab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job [subcommand]",
	Short: "Manage CI/CD jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var jobListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List jobs in a project",
	Args:  cobra.ExactArgs(1),
	Example: `  aidlc gitlab job list foundation/ai
  aidlc gitlab job list foundation/ai --scope running,pending
  aidlc gitlab job list foundation/ai --scope failed --limit 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		scopeStr, _ := cmd.Flags().GetString("scope")
		limit, _ := cmd.Flags().GetInt("limit")

		var scopes []string
		if scopeStr != "" {
			scopes = strings.Split(scopeStr, ",")
		}

		jobs, err := client.ListProjectJobs(args[0], scopes, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(jobs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-10s %-25s %-15s %-12s %-10s %-10s %s\n", "ID", "NAME", "STAGE", "STATUS", "DURATION", "QUEUED", "PIPELINE")
		fmt.Printf("%-10s %-25s %-15s %-12s %-10s %-10s %s\n", "--", "----", "-----", "------", "--------", "------", "--------")
		for _, j := range jobs {
			dur := ""
			if j.Duration > 0 {
				dur = formatDuration(j.Duration)
			}
			queued := ""
			if j.QueuedDuration > 0 {
				queued = formatDuration(j.QueuedDuration)
			}
			pipeID := ""
			if j.Pipeline != nil {
				pipeID = strconv.Itoa(j.Pipeline.ID)
			}
			name := j.Name
			if len(name) > 23 {
				name = name[:20] + "..."
			}
			fmt.Printf("%-10d %-25s %-15s %-12s %-10s %-10s %s\n", j.ID, name, j.Stage, j.Status, dur, queued, pipeID)
		}
		return nil
	},
}

var jobViewCmd = &cobra.Command{
	Use:   "view [project] [job-id]",
	Short: "View job details",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab job view foundation/ai 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", args[1])
		}

		job, err := client.GetJob(args[0], id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(job, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:             %d\n", job.ID)
		fmt.Printf("Name:           %s\n", job.Name)
		fmt.Printf("Stage:          %s\n", job.Stage)
		fmt.Printf("Status:         %s\n", job.Status)
		if job.FailureReason != "" {
			fmt.Printf("Failure Reason: %s\n", job.FailureReason)
		}
		fmt.Printf("Ref:            %s\n", job.Ref)
		fmt.Printf("Created:        %s\n", job.CreatedAt)
		if job.StartedAt != "" {
			fmt.Printf("Started:        %s\n", job.StartedAt)
		}
		if job.FinishedAt != "" {
			fmt.Printf("Finished:       %s\n", job.FinishedAt)
		}
		if job.Duration > 0 {
			fmt.Printf("Duration:       %s\n", formatDuration(job.Duration))
		}
		if job.QueuedDuration > 0 {
			fmt.Printf("Queued:         %s\n", formatDuration(job.QueuedDuration))
		}
		if job.AllowFailure {
			fmt.Printf("Allow Failure:  true\n")
		}
		if job.Runner != nil {
			fmt.Printf("Runner:         #%d %s (shared=%v, active=%v)\n", job.Runner.ID, job.Runner.Description, job.Runner.IsShared, job.Runner.Active)
		} else if job.Status == "pending" || job.Status == "created" {
			fmt.Printf("Runner:         (none assigned — job may be stuck)\n")
		}
		if job.User != nil {
			fmt.Printf("User:           %s\n", job.User.Username)
		}
		if job.Pipeline != nil {
			fmt.Printf("Pipeline:       %d\n", job.Pipeline.ID)
		}
		fmt.Printf("URL:            %s\n", job.WebURL)
		return nil
	},
}

var jobLogCmd = &cobra.Command{
	Use:     "log [project] [job-id]",
	Short:   "View job log/trace output",
	Aliases: []string{"trace"},
	Args:    cobra.ExactArgs(2),
	Example: `  aidlc gitlab job log foundation/ai 12345
  aidlc gitlab job log foundation/ai 12345 --tail 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", args[1])
		}

		log, err := client.GetJobLog(args[0], id)
		if err != nil {
			return err
		}

		tail, _ := cmd.Flags().GetInt("tail")
		if tail > 0 {
			lines := strings.Split(log, "\n")
			if len(lines) > tail {
				lines = lines[len(lines)-tail:]
			}
			fmt.Println(strings.Join(lines, "\n"))
			return nil
		}

		fmt.Println(log)
		return nil
	},
}

var jobRetryCmd = &cobra.Command{
	Use:   "retry [project] [job-id]",
	Short: "Retry a job",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab job retry foundation/ai 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", args[1])
		}

		job, err := client.RetryJob(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Retried job %d '%s' (status: %s)\n", job.ID, job.Name, job.Status)
		return nil
	},
}

var jobCancelCmd = &cobra.Command{
	Use:   "cancel [project] [job-id]",
	Short: "Cancel a running or pending job",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab job cancel foundation/ai 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", args[1])
		}

		job, err := client.CancelJob(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Canceled job %d '%s' (status: %s)\n", job.ID, job.Name, job.Status)
		return nil
	},
}

var jobPlayCmd = &cobra.Command{
	Use:   "play [project] [job-id]",
	Short: "Trigger a manual job",
	Args:  cobra.ExactArgs(2),
	Example: `  aidlc gitlab job play foundation/ai 253805`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitLabClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", args[1])
		}

		job, err := client.PlayJob(args[0], id)
		if err != nil {
			return err
		}

		fmt.Printf("Triggered job %d '%s' (status: %s)\n", job.ID, job.Name, job.Status)
		return nil
	},
}

func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}
	m := int(seconds) / 60
	s := int(seconds) % 60
	if m < 60 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	h := m / 60
	m = m % 60
	return fmt.Sprintf("%dh%dm", h, m)
}

func init() {
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobViewCmd)
	jobCmd.AddCommand(jobLogCmd)
	jobCmd.AddCommand(jobRetryCmd)
	jobCmd.AddCommand(jobCancelCmd)
	jobCmd.AddCommand(jobPlayCmd)

	jobListCmd.Flags().String("scope", "", "filter by status: created, pending, running, failed, success, canceled, skipped, manual")
	jobListCmd.Flags().Int("limit", 20, "max results")

	jobLogCmd.Flags().Int("tail", 0, "show only last N lines")
}

package github

import (
	"encoding/json"
	"fmt"
	"strconv"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release [subcommand]",
	Short: "Manage releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var releaseListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List releases",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github release list octocat/hello-world`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		releases, err := client.ListReleases(owner, repo, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(releases, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-20s %-30s %-12s %s\n", "TAG", "NAME", "PUBLISHED", "FLAGS")
		fmt.Printf("%-20s %-30s %-12s %s\n", "---", "----", "---------", "-----")
		for _, r := range releases {
			published := ""
			if len(r.PublishedAt) >= 10 {
				published = r.PublishedAt[:10]
			}
			flags := ""
			if r.Draft {
				flags = "draft"
			}
			if r.Prerelease {
				if flags != "" {
					flags += ", "
				}
				flags += "prerelease"
			}
			name := r.Name
			if len(name) > 28 {
				name = name[:25] + "..."
			}
			fmt.Printf("%-20s %-30s %-12s %s\n", r.TagName, name, published, flags)
		}
		return nil
	},
}

var releaseViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [id]",
	Short: "View a release",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github release view octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid release ID: %s", args[1])
		}

		r, err := client.GetRelease(owner, repo, id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Tag:        %s\n", r.TagName)
		fmt.Printf("Name:       %s\n", r.Name)
		fmt.Printf("Draft:      %v\n", r.Draft)
		fmt.Printf("Prerelease: %v\n", r.Prerelease)
		fmt.Printf("Published:  %s\n", r.PublishedAt)
		if r.Author != nil {
			fmt.Printf("Author:     %s\n", r.Author.Login)
		}
		fmt.Printf("URL:        %s\n", r.HTMLURL)
		if r.Body != "" {
			fmt.Printf("\n%s\n", r.Body)
		}
		return nil
	},
}

var releaseLatestCmd = &cobra.Command{
	Use:   "latest [owner/repo]",
	Short: "View the latest release",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github release latest octocat/hello-world`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		r, err := client.GetLatestRelease(owner, repo)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Tag:        %s\n", r.TagName)
		fmt.Printf("Name:       %s\n", r.Name)
		fmt.Printf("Published:  %s\n", r.PublishedAt)
		if r.Author != nil {
			fmt.Printf("Author:     %s\n", r.Author.Login)
		}
		fmt.Printf("URL:        %s\n", r.HTMLURL)
		if r.Body != "" {
			fmt.Printf("\n%s\n", r.Body)
		}
		return nil
	},
}

func init() {
	releaseCmd.AddCommand(releaseListCmd)
	releaseCmd.AddCommand(releaseViewCmd)
	releaseCmd.AddCommand(releaseLatestCmd)

	releaseListCmd.Flags().Int("limit", 20, "max results")
}

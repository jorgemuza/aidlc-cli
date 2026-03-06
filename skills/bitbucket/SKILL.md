---
name: bitbucket
description: "Manage Bitbucket repositories, pull requests, branches, tags, commits, and projects using the aidlc CLI. Use this skill whenever the user asks about Bitbucket repos, PRs (pull requests), branches, tags, commits, code review, or project management on Bitbucket Server/Data Center or Bitbucket Cloud. Trigger on phrases like 'list PRs', 'show pull requests', 'create a branch', 'open a PR', 'view the latest commits', 'list repos in project X', 'merge the PR', 'decline the PR', 'check PR activity', or any Bitbucket-related task — even casual references like 'what PRs are open', 'show me the repos', 'tag a release', 'check if it merged', 'who approved it', or 'list branches'. Also trigger when the user provides a Bitbucket Server URL (e.g., https://git.example.com/projects/PROJ/repos/my-repo/) or mentions Bitbucket Data Center. The aidlc CLI alias is `bb`."
---

# Bitbucket with aidlc CLI

Manage Bitbucket repositories, pull requests, branches, tags, commits, projects, and users through the `aidlc` CLI. Works with both Bitbucket Server/Data Center and Cloud, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `aidlc` binary built and accessible
2. A profile with a `bitbucket` service configured in `~/.config/aidlc/config.yaml`
3. Valid credentials (Personal Access Token or Bearer token for Server; app password for Cloud) — can be stored in 1Password with `op://` prefix

## Quick Reference

All commands follow the pattern: `aidlc -p <profile> bitbucket <command> [flags]`

Alias: `aidlc -p <profile> bb <command> [flags]`

All commands support `-o json` for JSON output. For full command details and all flags, see `references/commands.md`.

## Addressing Convention

Bitbucket Server/Data Center uses **project-key + repo-slug** to identify repositories. These are the two positional arguments most commands require:

- **Project key**: uppercase short code (e.g., `L3SUP`, `MYPROJ`)
- **Repo slug**: lowercase hyphenated name (e.g., `agents-sre`, `my-service`)

You can extract these from a Bitbucket Server URL:
```
https://git.example.com/projects/L3SUP/repos/agents-sre/
                                  ^^^^^       ^^^^^^^^^^
                              project-key     repo-slug
```

## Core Workflows

### Exploring Projects and Repos

```bash
# List all projects
aidlc -p myprofile bb project list

# View project details
aidlc -p myprofile bb project view L3SUP

# List repositories in a project
aidlc -p myprofile bb repo list L3SUP

# View repo details (includes clone URLs)
aidlc -p myprofile bb repo view L3SUP agents-sre
```

### Working with Pull Requests

Bitbucket uses "pull requests" (PR), same as GitHub.

```bash
# List open PRs (default)
aidlc -p myprofile bb pr list L3SUP agents-sre

# List merged PRs
aidlc -p myprofile bb pr list L3SUP agents-sre --state merged

# List all PRs
aidlc -p myprofile bb pr list L3SUP agents-sre --state all

# View PR details (shows from/to branches, reviewers, approval status)
aidlc -p myprofile bb pr view L3SUP agents-sre 42

# Create a PR
aidlc -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add new feature"

# Create a PR with reviewers
aidlc -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add feature" \
  --reviewers john.doe,jane.smith

# Merge a PR
aidlc -p myprofile bb pr merge L3SUP agents-sre 42

# Decline a PR
aidlc -p myprofile bb pr decline L3SUP agents-sre 42

# Add a comment to a PR
aidlc -p myprofile bb pr comment L3SUP agents-sre 42 --body "LGTM!"

# View PR activity (comments, approvals, status changes)
aidlc -p myprofile bb pr activity L3SUP agents-sre 42
```

**PR states:** `OPEN`, `MERGED`, `DECLINED`, `ALL`

### Branches and Tags

```bash
# List branches
aidlc -p myprofile bb branch list L3SUP agents-sre

# Filter branches by name
aidlc -p myprofile bb branch list L3SUP agents-sre --filter feature

# Show default branch
aidlc -p myprofile bb branch default L3SUP agents-sre

# Create a branch from a ref
aidlc -p myprofile bb branch create L3SUP agents-sre feature/new-thing main

# Delete a branch
aidlc -p myprofile bb branch delete L3SUP agents-sre feature/old-thing

# List tags
aidlc -p myprofile bb tag list L3SUP agents-sre

# Create a tag
aidlc -p myprofile bb tag create L3SUP agents-sre v1.0.0 main -m "Release v1.0.0"
```

### Commits

```bash
# List recent commits (default branch)
aidlc -p myprofile bb commit list L3SUP agents-sre

# List commits on a specific branch
aidlc -p myprofile bb commit list L3SUP agents-sre --branch feature/new

# View commit details
aidlc -p myprofile bb commit view L3SUP agents-sre abc1234def5678
```

### Users

```bash
# List users
aidlc -p myprofile bb user list

# Filter users by name
aidlc -p myprofile bb user list --filter john
```

## Common Patterns

**Get JSON for scripting:**
```bash
aidlc -p myprofile bb pr list L3SUP agents-sre -o json | jq '.[].title'
```

**Review a PR end-to-end:**
```bash
# View PR details and reviewers
aidlc -p myprofile bb pr view L3SUP agents-sre 42
# Check activity (approvals, comments)
aidlc -p myprofile bb pr activity L3SUP agents-sre 42
# Add your review comment
aidlc -p myprofile bb pr comment L3SUP agents-sre 42 --body "Approved, looks good"
```

**Extract project key and repo slug from a URL:**
Given `https://git.cnvrmedia.net/projects/L3SUP/repos/agents-sre/pull-requests`:
- Project key: `L3SUP`
- Repo slug: `agents-sre`
Then run: `aidlc -p myprofile bb pr list L3SUP agents-sre`

**Check what repos exist in a project:**
```bash
aidlc -p myprofile bb repo list L3SUP
```

## Important Notes

- **Profile required** — Always pass `-p <profile>` to select the Bitbucket connection. The profile must have a service of type `bitbucket` configured.
- **Service flag** — If a profile has multiple Bitbucket services, use `--service <name>` to disambiguate.
- **Server vs Cloud** — The service variant (`server` or `cloud`) in config determines the API prefix. Server uses `/rest/api/latest/`, Cloud uses `/2.0/`.
- **1Password integration** — Auth tokens in config can use `op://vault/item/field` and are resolved at runtime.
- **PR states are uppercase** — Use `OPEN`, `MERGED`, `DECLINED`, or `ALL` (case-insensitive input is accepted).
- **Pagination** — Most list commands default to 25-50 results. Use `--limit N` to adjust.
- **URL parsing** — When a user provides a Bitbucket Server URL like `https://host/projects/KEY/repos/SLUG/...`, extract the project key and repo slug from the URL path to use with aidlc commands.

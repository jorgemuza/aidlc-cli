---
name: github
description: "Create and manage GitHub repositories, pull requests, issues, releases, branches, secrets, and more using the orbit CLI. Use this skill whenever the user asks about GitHub repositories, PRs (pull requests), GitHub Actions workflow runs, branches, tags, commits, issues, releases, secrets, or organization repos. Trigger on phrases like 'list PRs', 'check the actions', 'watch the workflow', 'create a secret', 'open a pull request', 'view the latest commits', 'list repos in org X', 'rerun the workflow', 'close the issue', 'latest release', 'set a GitHub secret', or any GitHub-related task — even casual references like 'what's running in CI', 'show me the PRs', 'tag a release', 'check if it merged', 'list repos', 'is the build passing', or 'add a deploy key secret'. Also trigger when the user wants to monitor CI/CD progress, manage Actions secrets for deployments, or debug failing workflows. The orbit CLI alias is `gh`."
---

# GitHub with orbit CLI

Manage GitHub repositories, pull requests, issues, releases, branches, tags, commits, workflow runs, secrets, and users through the `orbit` CLI. Works with both GitHub.com and GitHub Enterprise via REST API, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. A profile with a `github` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials (Personal Access Token) - can be stored as 1Password (`op://`) or Infisical (`infisical://`) references

## Quick Reference

All commands follow the pattern: `orbit -p <profile> github <command> [flags]`

Alias: `orbit -p <profile> gh <command> [flags]`

All commands support `-o json` for JSON output. For full command details and all flags, see `references/commands.md`.

## Repository Identification

Repositories are always referenced as `owner/repo`:
- `orbit -p myprofile gh repo octocat/hello-world`
- `orbit -p myprofile gh repo kubernetes/kubernetes`

## Core Workflows

### Exploring Repositories

```bash
# View repo details
orbit -p myprofile gh repo octocat/hello-world

# List your repos (sorted by most recently pushed)
orbit -p myprofile gh repos

# List repos in an organization
orbit -p myprofile gh repos --org kubernetes

# List repos with limit
orbit -p myprofile gh repos --limit 10

# Edit repo settings
orbit -p myprofile gh repo edit Paybook/ai --description "Updated"

# Archive a repo
orbit -p myprofile gh repo edit Paybook/ai --archived

# List collaborators
orbit -p myprofile gh repo collab list Paybook/ai

# Add collaborator
orbit -p myprofile gh repo collab add Paybook/ai jorgemuza --permission admin

# Remove collaborator
orbit -p myprofile gh repo collab remove Paybook/ai jorgemuza
```

### Working with Pull Requests

```bash
# List open PRs
orbit -p myprofile gh pr list octocat/hello-world

# List closed PRs
orbit -p myprofile gh pr list octocat/hello-world --state closed

# View PR details (shows head/base branch, labels, comments)
orbit -p myprofile gh pr view octocat/hello-world 42

# Create a PR
orbit -p myprofile gh pr create octocat/hello-world \
  --head feature/login --base main --title "Add login page"

# Merge a PR (with optional method: merge, squash, rebase)
orbit -p myprofile gh pr merge octocat/hello-world 42 --method squash

# Add a comment
orbit -p myprofile gh pr comment octocat/hello-world 42 --body "LGTM!"

# List comments
orbit -p myprofile gh pr comments octocat/hello-world 42
```

### GitHub Actions Workflow Runs

```bash
# List recent workflow runs
orbit -p myprofile gh run list octocat/hello-world

# Filter by branch and status
orbit -p myprofile gh run list octocat/hello-world --branch main --status completed

# View workflow run details
orbit -p myprofile gh run view octocat/hello-world 12345

# Watch a run in real-time (polls and shows job/step progress until completion)
orbit -p myprofile gh run watch octocat/hello-world
orbit -p myprofile gh run watch octocat/hello-world 12345 --interval 10

# Cancel a running workflow
orbit -p myprofile gh run cancel octocat/hello-world 12345

# Re-run a workflow
orbit -p myprofile gh run rerun octocat/hello-world 12345
```

Run aliases: `run`, `actions` — so `orbit gh actions list octocat/hello-world` works too.

The `watch` command auto-discovers the most recent in-progress run if no run-id is given. It shows live job and step status with elapsed time, and exits with an error if the run fails.

### Workflow Management

```bash
# List workflows (get the numeric workflow ID)
orbit -p myprofile gh workflow list octocat/hello-world

# Trigger a workflow dispatch (--ref is required)
orbit -p myprofile gh workflow run octocat/hello-world 245836153 --ref main

# Trigger with inputs
orbit -p myprofile gh workflow run octocat/hello-world 245836153 --ref main --input env=staging

# Enable/disable a workflow
orbit -p myprofile gh workflow enable octocat/hello-world 245836153
orbit -p myprofile gh workflow disable octocat/hello-world 245836153
```

**Important:** `workflow run` requires the numeric workflow ID (from `workflow list`), not the filename. The `--ref` flag is mandatory.

### GitHub Actions Secrets

```bash
# List repository secrets
orbit -p myprofile gh secret list octocat/hello-world

# Create or update a secret
orbit -p myprofile gh secret set octocat/hello-world MY_SECRET "secret-value"

# Delete a secret
orbit -p myprofile gh secret delete octocat/hello-world MY_SECRET
```

Secrets are encrypted client-side using the repository's public key before being sent to the API.

### Branches and Tags

```bash
# List branches
orbit -p myprofile gh branch list octocat/hello-world

# View branch details (includes latest commit)
orbit -p myprofile gh branch view octocat/hello-world main

# List tags
orbit -p myprofile gh tag list octocat/hello-world
```

### Commits

```bash
# List recent commits (default branch)
orbit -p myprofile gh commit list octocat/hello-world

# List commits on a specific branch
orbit -p myprofile gh commit list octocat/hello-world --ref feature/login

# View commit details
orbit -p myprofile gh commit view octocat/hello-world abc1234
```

### Issues

```bash
# List open issues
orbit -p myprofile gh issue list octocat/hello-world --state open

# Filter by labels
orbit -p myprofile gh issue list octocat/hello-world --labels bug,urgent

# View issue details
orbit -p myprofile gh issue view octocat/hello-world 1

# Create an issue
orbit -p myprofile gh issue create octocat/hello-world --title "Fix login bug" --labels bug,urgent

# Close an issue
orbit -p myprofile gh issue close octocat/hello-world 1

# Add a comment to an issue
orbit -p myprofile gh issue comment octocat/hello-world 1 --body "Working on this"
```

### Releases

```bash
# List releases
orbit -p myprofile gh release list octocat/hello-world

# View a specific release
orbit -p myprofile gh release view octocat/hello-world 12345

# View the latest release
orbit -p myprofile gh release latest octocat/hello-world
```

### Users

```bash
# Show current authenticated user
orbit -p myprofile gh user me

# View a user profile
orbit -p myprofile gh user view octocat
```

## Common Patterns

**Get JSON for scripting:**
Any command supports `-o json` for machine-readable output:
```bash
orbit -p myprofile gh pr list octocat/hello-world -o json | jq '.[].title'
```

**Check CI status for a branch:**
```bash
orbit -p myprofile gh run list octocat/hello-world --branch main --limit 1
```

**Monitor a release pipeline:**
```bash
orbit -p myprofile gh run watch octocat/hello-world
```

**Set a deployment secret:**
```bash
orbit -p myprofile gh secret set octocat/hello-world DEPLOY_TOKEN "ghp_xxxxx"
```

**Review a PR end-to-end:**
```bash
# View PR details
orbit -p myprofile gh pr view octocat/hello-world 42
# Check its workflow runs
orbit -p myprofile gh run list octocat/hello-world --branch feature/login --limit 1
# Read discussion
orbit -p myprofile gh pr comments octocat/hello-world 42
# Approve with comment
orbit -p myprofile gh pr comment octocat/hello-world 42 --body "Approved, looks good"
```

## Important Notes

- **Profile required** — Always pass `-p <profile>` to select the GitHub connection. The profile must have a service of type `github` configured.
- **Service flag** — If a profile has multiple GitHub services, use `--service <name>` to disambiguate.
- **Cloud vs Enterprise** — Works with both. For GitHub.com the base_url defaults to `https://api.github.com`. For GitHub Enterprise, set the base_url in your profile config.
- **Secret references** - Credentials in config can use 1Password (`op://vault/item/field`) or Infisical (`infisical://<env>/<path>/<KEY>`) references, resolved at runtime. Run `orbit auth` once to resolve and cache all secrets (a single biometric prompt for 1Password). Use `orbit auth clear` to wipe the cache. See [Secrets](../../docs/secrets.md).
- **Pagination** — Most list commands default to 20-50 results. Use `--limit N` to adjust.

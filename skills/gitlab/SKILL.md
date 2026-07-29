---
name: gitlab
description: "Create and manage GitLab projects, merge requests, pipelines, issues, branches, and more using the orbit CLI. Use this skill whenever the user asks about GitLab repositories, MRs (merge requests), CI/CD pipelines, branches, tags, commits, issues, groups, or project members. Trigger on phrases like 'list MRs', 'check the pipeline', 'create a branch', 'open a merge request', 'view the latest commits', 'list projects in group X', 'retry the CI', 'close the issue', 'who are the members', or any GitLab-related task — even casual references like 'what's running in CI', 'show me the MRs', 'tag a release', 'check if it merged', or 'list repos'. Also trigger when the user mentions PR/pull request in a GitLab context (GitLab calls them merge requests). The orbit CLI alias is `gl`."
---

# GitLab with orbit CLI

Manage GitLab projects, merge requests, pipelines, pipeline schedules, issues, branches, tags, commits, members, and users through the `orbit` CLI. Works with both GitLab Cloud and self-hosted instances via REST API v4, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. A profile with a `gitlab` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials (Personal Access Token or Bearer token) - can be stored as 1Password (`op://`) or Infisical (`infisical://`) references

## Quick Reference

All commands follow the pattern: `orbit -p <profile> gitlab <command> [flags]`

Alias: `orbit -p <profile> gl <command> [flags]`

All commands support `-o json` for JSON output. For full command details and all flags, see `references/commands.md`.

For self-hosted instances with self-signed certificates, add `tls_skip_verify: true` to the service config. For proxy access, add `proxy: socks5://host:port` (also supports `http://` and `https://`).

## Project Identification

Projects can be referenced by numeric ID or full path:
- `orbit -p myprofile gl project 595`
- `orbit -p myprofile gl project schools/frontend/my-app`

Groups also accept ID or full path: `schools/frontend`

## Core Workflows

### Exploring Projects and Groups

```bash
# View project details
orbit -p myprofile gl project schools/frontend/my-app

# List your projects (membership-based)
orbit -p myprofile gl project list --search frontend

# Edit project settings (default branch, visibility, description)
orbit -p myprofile gl project edit 595 --default-branch main
orbit -p myprofile gl project edit 595 --visibility private
orbit -p myprofile gl project edit 595 --archived

# List all projects in a group (includes subgroups)
orbit -p myprofile gl project list --group schools/frontend

# Projects with activity today (shows branches)
orbit -p myprofile gl project activity

# Activity in the last 7 days, filtered by branch
orbit -p myprofile gl project activity --days 7 --branch development,master

# View group info
orbit -p myprofile gl group view schools/frontend

# List subgroups
orbit -p myprofile gl group subgroups schools
```

### Working with Merge Requests

GitLab calls them "merge requests" (MR), equivalent to GitHub's "pull requests" (PR).

```bash
# List open MRs
orbit -p myprofile gl mr list 595

# List merged MRs
orbit -p myprofile gl mr list 595 --state merged

# View MR details (shows source/target branch, conflicts, review status)
orbit -p myprofile gl mr view 595 42

# Create an MR
orbit -p myprofile gl mr create 595 \
  --source feature/login --target main --title "Add login page"

# Merge an MR (with optional squash)
orbit -p myprofile gl mr merge 595 42 --squash

# Approve / unapprove
orbit -p myprofile gl mr approve 595 42
orbit -p myprofile gl mr unapprove 595 42

# Check approval status
orbit -p myprofile gl mr approvals 595 42

# Rebase onto target branch
orbit -p myprofile gl mr rebase 595 42

# Close / reopen
orbit -p myprofile gl mr close 595 42
orbit -p myprofile gl mr reopen 595 42

# Edit MR (title, description, assignee, labels, draft)
orbit -p myprofile gl mr edit 595 42 --title "Updated title"
orbit -p myprofile gl mr edit 595 42 --assignee 15 --labels "bug,urgent"

# Add a comment
orbit -p myprofile gl mr comment 595 42 --body "LGTM!"

# List discussion comments (excludes system notes)
orbit -p myprofile gl mr notes 595 42
```

### CI/CD Pipelines

```bash
# List recent pipelines
orbit -p myprofile gl pipeline list 595

# Filter by branch and status
orbit -p myprofile gl pipeline list 595 --ref main --status failed

# View pipeline details
orbit -p myprofile gl pipeline view 595 12345

# List jobs in a pipeline (shows stage, status, duration)
orbit -p myprofile gl pipeline jobs 595 12345

# Retry a failed pipeline
orbit -p myprofile gl pipeline retry 595 12345

# Cancel a running pipeline
orbit -p myprofile gl pipeline cancel 595 12345
```

Pipeline aliases: `pipeline`, `pipe`, `ci` — so `orbit gl ci list 595` works too.

# List pipeline schedules
orbit -p myprofile gl schedule list 595

# Create a nightly schedule
orbit -p myprofile gl schedule create 595 --desc "Nightly build" --ref main --cron "0 2 * * *"

# Trigger a schedule immediately
orbit -p myprofile gl schedule run 595 42

# Update schedule cron or disable it
orbit -p myprofile gl schedule update 595 42 --cron "0 3 * * *"
orbit -p myprofile gl schedule update 595 42 --active=false

# Add/remove schedule variables
orbit -p myprofile gl schedule var 595 42 DEPLOY_ENV production
orbit -p myprofile gl schedule var-delete 595 42 DEPLOY_ENV

# Delete a schedule
orbit -p myprofile gl schedule delete 595 42

Schedule aliases: `schedule`, `sched`.

### Branches and Tags

```bash
# List branches
orbit -p myprofile gl branch list 595 --search feature

# View branch details (includes latest commit)
orbit -p myprofile gl branch view 595 main

# Create a branch from a ref
orbit -p myprofile gl branch create 595 feature/new-thing main

# Delete a branch
orbit -p myprofile gl branch delete 595 feature/old-thing

# Protect a branch (only maintainers can merge, no direct push)
orbit -p myprofile gl branch protect 650 main --push no-access --merge maintainer

# List protected branches
orbit -p myprofile gl branch protections 650

# Remove protection
orbit -p myprofile gl branch unprotect 650 main

# List tags
orbit -p myprofile gl tag list 595

# Create an annotated tag
orbit -p myprofile gl tag create 595 v1.0.0 main -m "Release v1.0.0"
```

### Commits

```bash
# List recent commits (default branch)
orbit -p myprofile gl commit list 595

# List commits on a specific branch
orbit -p myprofile gl commit list 595 --ref feature/login

# View commit details
orbit -p myprofile gl commit view 595 abc1234
```

### Issues

```bash
# List open issues
orbit -p myprofile gl issue list 595 --state opened

# Filter by labels
orbit -p myprofile gl issue list 595 --labels bug,urgent

# View issue details
orbit -p myprofile gl issue view 595 1

# Create an issue
orbit -p myprofile gl issue create 595 --title "Fix login bug" --labels bug,urgent

# Close an issue
orbit -p myprofile gl issue close 595 1
```

### Members and Users

```bash
# List project members (shows access level: Guest/Reporter/Developer/Maintainer/Owner)
orbit -p myprofile gl member list 595

# Show current authenticated user
orbit -p myprofile gl user me

# Search users
orbit -p myprofile gl user list --search john
```

## Common Patterns

**Get JSON for scripting:**
Any command supports `-o json` for machine-readable output:
```bash
orbit -p myprofile gl mr list 595 -o json | jq '.[].title'
```

**Check CI status for a branch:**
```bash
orbit -p myprofile gl pipeline list 595 --ref main --limit 1
```

**Find who's working on a project:**
```bash
orbit -p myprofile gl member list 595
```

**Review an MR end-to-end:**
```bash
# View MR details
orbit -p myprofile gl mr view 595 42
# Check its pipeline
orbit -p myprofile gl pipeline list 595 --ref feature/login --limit 1
# Read discussion
orbit -p myprofile gl mr notes 595 42
# Approve with comment
orbit -p myprofile gl mr comment 595 42 --body "Approved, looks good"
```

## Important Notes

- **Profile required** — Always pass `-p <profile>` to select the GitLab connection. The profile must have a service of type `gitlab` configured.
- **Service flag** — If a profile has multiple GitLab services, use `--service <name>` to disambiguate.
- **Cloud vs Self-hosted** — Works with both. The base URL in your profile config determines the GitLab instance.
- **Secret references** - Credentials in config can use 1Password (`op://vault/item/field`) or Infisical (`infisical://<env>/<path>/<KEY>`) references, resolved at runtime. Run `orbit auth` once to resolve and cache all secrets (a single biometric prompt for 1Password). Use `orbit auth clear` to wipe the cache. See [Secrets](../../docs/secrets.md).
- **MR = PR** — If a user says "pull request" or "PR" in a GitLab context, they mean merge request.
- **Pagination** — Most list commands default to 20-50 results. Use `--limit N` to adjust.

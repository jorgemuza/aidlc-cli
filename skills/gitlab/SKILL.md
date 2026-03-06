---
name: gitlab
description: "Create and manage GitLab projects, merge requests, pipelines, issues, branches, and more using the aidlc CLI. Use this skill whenever the user asks about GitLab repositories, MRs (merge requests), CI/CD pipelines, branches, tags, commits, issues, groups, or project members. Trigger on phrases like 'list MRs', 'check the pipeline', 'create a branch', 'open a merge request', 'view the latest commits', 'list projects in group X', 'retry the CI', 'close the issue', 'who are the members', or any GitLab-related task — even casual references like 'what's running in CI', 'show me the MRs', 'tag a release', 'check if it merged', or 'list repos'. Also trigger when the user mentions PR/pull request in a GitLab context (GitLab calls them merge requests). The aidlc CLI alias is `gl`."
---

# GitLab with aidlc CLI

Manage GitLab projects, merge requests, pipelines, issues, branches, tags, commits, members, and users through the `aidlc` CLI. Works with both GitLab Cloud and self-hosted instances via REST API v4, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `aidlc` binary built and accessible
2. A profile with a `gitlab` service configured in `~/.config/aidlc/config.yaml`
3. Valid credentials (Personal Access Token or Bearer token) — can be stored in 1Password with `op://` prefix

## Quick Reference

All commands follow the pattern: `aidlc -p <profile> gitlab <command> [flags]`

Alias: `aidlc -p <profile> gl <command> [flags]`

All commands support `-o json` for JSON output. For full command details and all flags, see `references/commands.md`.

## Project Identification

Projects can be referenced by numeric ID or full path:
- `aidlc -p myprofile gl project 595`
- `aidlc -p myprofile gl project schools/frontend/my-app`

Groups also accept ID or full path: `schools/frontend`

## Core Workflows

### Exploring Projects and Groups

```bash
# View project details
aidlc -p myprofile gl project schools/frontend/my-app

# List your projects (membership-based)
aidlc -p myprofile gl projects --search frontend

# List all projects in a group (includes subgroups)
aidlc -p myprofile gl projects --group schools/frontend

# View group info
aidlc -p myprofile gl group view schools/frontend

# List subgroups
aidlc -p myprofile gl group subgroups schools
```

### Working with Merge Requests

GitLab calls them "merge requests" (MR), equivalent to GitHub's "pull requests" (PR).

```bash
# List open MRs
aidlc -p myprofile gl mr list 595

# List merged MRs
aidlc -p myprofile gl mr list 595 --state merged

# View MR details (shows source/target branch, conflicts, review status)
aidlc -p myprofile gl mr view 595 42

# Create an MR
aidlc -p myprofile gl mr create 595 \
  --source feature/login --target main --title "Add login page"

# Merge an MR (with optional squash)
aidlc -p myprofile gl mr merge 595 42 --squash

# Add a comment
aidlc -p myprofile gl mr comment 595 42 --body "LGTM!"

# List discussion comments (excludes system notes)
aidlc -p myprofile gl mr notes 595 42
```

### CI/CD Pipelines

```bash
# List recent pipelines
aidlc -p myprofile gl pipeline list 595

# Filter by branch and status
aidlc -p myprofile gl pipeline list 595 --ref main --status failed

# View pipeline details
aidlc -p myprofile gl pipeline view 595 12345

# List jobs in a pipeline (shows stage, status, duration)
aidlc -p myprofile gl pipeline jobs 595 12345

# Retry a failed pipeline
aidlc -p myprofile gl pipeline retry 595 12345

# Cancel a running pipeline
aidlc -p myprofile gl pipeline cancel 595 12345
```

Pipeline aliases: `pipeline`, `pipe`, `ci` — so `aidlc gl ci list 595` works too.

### Branches and Tags

```bash
# List branches
aidlc -p myprofile gl branch list 595 --search feature

# View branch details (includes latest commit)
aidlc -p myprofile gl branch view 595 main

# Create a branch from a ref
aidlc -p myprofile gl branch create 595 feature/new-thing main

# Delete a branch
aidlc -p myprofile gl branch delete 595 feature/old-thing

# List tags
aidlc -p myprofile gl tag list 595

# Create an annotated tag
aidlc -p myprofile gl tag create 595 v1.0.0 main -m "Release v1.0.0"
```

### Commits

```bash
# List recent commits (default branch)
aidlc -p myprofile gl commit list 595

# List commits on a specific branch
aidlc -p myprofile gl commit list 595 --ref feature/login

# View commit details
aidlc -p myprofile gl commit view 595 abc1234
```

### Issues

```bash
# List open issues
aidlc -p myprofile gl issue list 595 --state opened

# Filter by labels
aidlc -p myprofile gl issue list 595 --labels bug,urgent

# View issue details
aidlc -p myprofile gl issue view 595 1

# Create an issue
aidlc -p myprofile gl issue create 595 --title "Fix login bug" --labels bug,urgent

# Close an issue
aidlc -p myprofile gl issue close 595 1
```

### Members and Users

```bash
# List project members (shows access level: Guest/Reporter/Developer/Maintainer/Owner)
aidlc -p myprofile gl member list 595

# Show current authenticated user
aidlc -p myprofile gl user me

# Search users
aidlc -p myprofile gl user list --search john
```

## Common Patterns

**Get JSON for scripting:**
Any command supports `-o json` for machine-readable output:
```bash
aidlc -p myprofile gl mr list 595 -o json | jq '.[].title'
```

**Check CI status for a branch:**
```bash
aidlc -p myprofile gl pipeline list 595 --ref main --limit 1
```

**Find who's working on a project:**
```bash
aidlc -p myprofile gl member list 595
```

**Review an MR end-to-end:**
```bash
# View MR details
aidlc -p myprofile gl mr view 595 42
# Check its pipeline
aidlc -p myprofile gl pipeline list 595 --ref feature/login --limit 1
# Read discussion
aidlc -p myprofile gl mr notes 595 42
# Approve with comment
aidlc -p myprofile gl mr comment 595 42 --body "Approved, looks good"
```

## Important Notes

- **Profile required** — Always pass `-p <profile>` to select the GitLab connection. The profile must have a service of type `gitlab` configured.
- **Service flag** — If a profile has multiple GitLab services, use `--service <name>` to disambiguate.
- **Cloud vs Self-hosted** — Works with both. The base URL in your profile config determines the GitLab instance.
- **1Password integration** — Auth tokens in config can use `op://vault/item/field` and are resolved at runtime.
- **MR = PR** — If a user says "pull request" or "PR" in a GitLab context, they mean merge request.
- **Pagination** — Most list commands default to 20-50 results. Use `--limit N` to adjust.

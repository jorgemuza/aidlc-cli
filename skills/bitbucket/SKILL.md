---
name: bitbucket
description: "Manage Bitbucket repositories, pull requests, branches, tags, commits, projects, and admin settings using the orbit CLI. Use this skill whenever the user asks about Bitbucket repos, PRs (pull requests), branches, tags, commits, code review, project management, default reviewer conditions, required approvals, merge restrictions, or PR approvals on Bitbucket Server/Data Center or Bitbucket Cloud. Trigger on phrases like 'list PRs', 'show pull requests', 'create a branch', 'open a PR', 'view the latest commits', 'list repos in project X', 'merge the PR', 'decline the PR', 'approve the PR', 'unapprove', 'request changes', 'needs work', 'mark as needs work', 'reject the PR', 'block the merge', 'check PR activity', 'bypass merge check', 'required approvals', 'reviewer conditions', 'who needs to approve', or any Bitbucket-related task — even casual references like 'what PRs are open', 'show me the repos', 'tag a release', 'check if it merged', 'who approved it', 'list branches', or 'why can't I merge'. Also trigger when the user provides a Bitbucket Server URL (e.g., https://git.example.com/projects/PROJ/repos/my-repo/) or mentions Bitbucket Data Center. The orbit CLI alias is `bb`."
---

# Bitbucket with orbit CLI

Manage Bitbucket repositories, pull requests, branches, tags, commits, projects, and users through the `orbit` CLI. Works with both Bitbucket Server/Data Center and Cloud, with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. A profile with a `bitbucket` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials (Personal Access Token or Bearer token for Server; app password for Cloud) - can be stored as 1Password (`op://`) or Infisical (`infisical://`) references

## Quick Reference

All commands follow the pattern: `orbit -p <profile> bitbucket <command> [flags]`

Alias: `orbit -p <profile> bb <command> [flags]`

All commands support `-o json` for JSON output. For full command details and all flags, see `references/commands.md`.

For self-hosted instances with self-signed certificates, add `tls_skip_verify: true` to the service config. For proxy access, add `proxy: socks5://host:port`.

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
orbit -p myprofile bb project list

# View project details
orbit -p myprofile bb project view L3SUP

# List repositories in a project
orbit -p myprofile bb repo list L3SUP

# View repo details (includes clone URLs)
orbit -p myprofile bb repo view L3SUP agents-sre
```

### Managing Repositories

```bash
# Create a repo
orbit -p myprofile bb repo create L3SUP --name my-service --description "Payment API"

# Update a repo
orbit -p myprofile bb repo edit L3SUP my-service --description "New description"

# Delete a repo (irreversible)
orbit -p myprofile bb repo delete L3SUP my-service

# Fork a repo into another project
orbit -p myprofile bb repo fork L3SUP agents-sre --target-project MYPROJ

# List forks
orbit -p myprofile bb repo forks L3SUP agents-sre

# Manage user permissions (REPO_READ, REPO_WRITE, REPO_ADMIN)
orbit -p myprofile bb repo permissions L3SUP agents-sre
orbit -p myprofile bb repo grant L3SUP agents-sre john.doe --permission REPO_WRITE
orbit -p myprofile bb repo revoke L3SUP agents-sre john.doe
```

### Working with Pull Requests

Bitbucket uses "pull requests" (PR), same as GitHub.

```bash
# List open PRs (default)
orbit -p myprofile bb pr list L3SUP agents-sre

# List merged PRs
orbit -p myprofile bb pr list L3SUP agents-sre --state merged

# List all PRs
orbit -p myprofile bb pr list L3SUP agents-sre --state all

# View PR details (shows from/to branches, reviewers, approval status)
orbit -p myprofile bb pr view L3SUP agents-sre 42

# Create a PR
orbit -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add new feature"

# Create a PR with reviewers
orbit -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add feature" \
  --reviewers john.doe,jane.smith

# View PR diff (for code review)
orbit -p myprofile bb pr diff L3SUP agents-sre 42

# View PR diff with more context lines
orbit -p myprofile bb pr diff L3SUP agents-sre 42 --context 10

# Merge a PR
orbit -p myprofile bb pr merge L3SUP agents-sre 42

# Merge a PR bypassing review checks (requires repo admin)
orbit -p myprofile bb pr merge L3SUP agents-sre 42 --bypass-review

# Approve a PR
orbit -p myprofile bb pr approve L3SUP agents-sre 42

# Remove approval from a PR
orbit -p myprofile bb pr unapprove L3SUP agents-sre 42

# Request changes on a PR (Bitbucket Server NEEDS_WORK). Aliases: request-changes.
# Pair with 'pr comment' to leave the actionable feedback first, then flip the status.
orbit -p myprofile bb pr comment L3SUP agents-sre 42 -m "Please address the items above"
orbit -p myprofile bb pr needs-work L3SUP agents-sre 42

# Decline a PR (closes it — use needs-work for "request changes" instead)
orbit -p myprofile bb pr decline L3SUP agents-sre 42

# Add a comment to a PR
orbit -p myprofile bb pr comment L3SUP agents-sre 42 --body "LGTM!"

# View PR activity (comments, approvals, status changes)
orbit -p myprofile bb pr activity L3SUP agents-sre 42
```

**PR states:** `OPEN`, `MERGED`, `DECLINED`, `ALL`

### Branches and Tags

```bash
# List branches
orbit -p myprofile bb branch list L3SUP agents-sre

# Filter branches by name
orbit -p myprofile bb branch list L3SUP agents-sre --filter feature

# Show default branch
orbit -p myprofile bb branch default L3SUP agents-sre

# Create a branch from a ref
orbit -p myprofile bb branch create L3SUP agents-sre feature/new-thing main

# Delete a branch
orbit -p myprofile bb branch delete L3SUP agents-sre feature/old-thing

# List tags
orbit -p myprofile bb tag list L3SUP agents-sre

# Create a tag
orbit -p myprofile bb tag create L3SUP agents-sre v1.0.0 main -m "Release v1.0.0"
```

### Commits

```bash
# List recent commits (default branch)
orbit -p myprofile bb commit list L3SUP agents-sre

# List commits on a specific branch
orbit -p myprofile bb commit list L3SUP agents-sre --branch feature/new

# View commit details
orbit -p myprofile bb commit view L3SUP agents-sre abc1234def5678
```

### Users

```bash
# List users
orbit -p myprofile bb user list

# Filter users by name
orbit -p myprofile bb user list --filter john
```

### Default Reviewer Conditions (Admin)

Manage project-level default reviewer conditions that auto-assign reviewers and enforce required approvals on PRs. Alias: `rc`.

```bash
# List all default reviewer conditions for a project
orbit -p myprofile bb reviewer-condition list EPCAP
orbit -p myprofile bb rc list EPCAP -o json

# Update required approvals (e.g., temporarily set to 0 to bypass)
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 0

# Restore required approvals
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 2

# Delete a condition
orbit -p myprofile bb rc delete EPCAP 1063
```

**Bypass merge block from required reviewers:**

When a PR merge is blocked by "Not all required reviewers have approved yet", this is enforced by a project-level default reviewer condition (not a repo merge hook). Use `--bypass-review` on `pr merge` which automatically handles both merge hooks AND default reviewer conditions:
```bash
orbit -p myprofile bb pr merge EPCAP my-repo 42 --bypass-review
```

Or manually: list conditions to find the blocking one, set its required approvals to 0, merge, then restore.

## Common Patterns

**Get JSON for scripting:**
```bash
orbit -p myprofile bb pr list L3SUP agents-sre -o json | jq '.[].title'
```

**Review a PR end-to-end:**
```bash
# View PR details and reviewers (can run in parallel with activity)
orbit -p myprofile bb pr view L3SUP agents-sre 42
orbit -p myprofile bb pr activity L3SUP agents-sre 42

# Get the full diff for code review (run SEPARATELY, not in parallel)
orbit -p myprofile bb pr diff L3SUP agents-sre 42

# Outcome: approve (all good)
orbit -p myprofile bb pr approve L3SUP agents-sre 42
orbit -p myprofile bb pr comment L3SUP agents-sre 42 --body "LGTM!"

# Outcome: request changes (three issues need fixing before merge)
orbit -p myprofile bb pr comment L3SUP agents-sre 42 --body "Please address: 1) ..., 2) ..., 3) ..."
orbit -p myprofile bb pr needs-work L3SUP agents-sre 42
```

**Picking the right review outcome:**

| Outcome | Command | When to use |
|---------|---------|-------------|
| Approve | `pr approve` | PR is ready to merge; no blockers. |
| Request changes | `pr needs-work` | PR needs fixes but is still the right direction. Author addresses feedback and pushes; you flip to `approve`. PR stays open. |
| Reject | `pr decline` | PR is fundamentally wrong and shouldn't be merged at all. **Closes the PR**; reversible only via reopen. Rarely the right choice after a normal review — prefer `needs-work`. |

**Extract project key and repo slug from a URL:**
Given `https://git.cnvrmedia.net/projects/L3SUP/repos/agents-sre/pull-requests`:
- Project key: `L3SUP`
- Repo slug: `agents-sre`
Then run: `orbit -p myprofile bb pr list L3SUP agents-sre`

**Check what repos exist in a project:**
```bash
orbit -p myprofile bb repo list L3SUP
```

## Important Notes

- **Do NOT run `pr diff` in parallel with other commands.** The diff endpoint returns raw text (not JSON) and can interfere with parallel JSON-based commands. Always run `pr diff` sequentially — never in the same parallel Bash block as other orbit commands.
- **Profile required** — Always pass `-p <profile>` to select the Bitbucket connection. The profile must have a service of type `bitbucket` configured.
- **Service flag** — If a profile has multiple Bitbucket services, use `--service <name>` to disambiguate.
- **Server vs Cloud** — The service variant (`server` or `cloud`) in config determines the API prefix. Server uses `/rest/api/latest/`, Cloud uses `/2.0/`.
- **Secret references** - Credentials in config can use 1Password (`op://vault/item/field`) or Infisical (`infisical://<env>/<path>/<KEY>`) references, resolved at runtime. Run `orbit auth` once to resolve and cache all secrets (a single biometric prompt for 1Password). Use `orbit auth clear` to wipe the cache. See [Secrets](../../docs/secrets.md).
- **PR states are uppercase** — Use `OPEN`, `MERGED`, `DECLINED`, or `ALL` (case-insensitive input is accepted).
- **Pagination** — Most list commands default to 25-50 results. Use `--limit N` to adjust.
- **URL parsing** — When a user provides a Bitbucket Server URL like `https://host/projects/KEY/repos/SLUG/...`, extract the project key and repo slug from the URL path to use with orbit commands.

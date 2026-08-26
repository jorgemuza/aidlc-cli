---
name: forgejo
description: "Manage Forgejo and Gitea repositories, pull requests, issues, Actions workflow runs and releases using the orbit CLI. Use this skill whenever the user asks about a self-hosted git instance running Forgejo or Gitea - repos, PRs (pull requests), code review, issues, CI runs, job logs, or releases. Trigger on phrases like 'list PRs', 'approve the PR', 'request changes', 'merge the PR', 'why did CI fail', 'show the failed job log', 'what's running in Forgejo actions', 'list issues', 'close the issue', 'open a PR', 'check the run', or any task against a Forgejo/Gitea host - including casual references like 'is the build green', 'what broke', 'show me the PRs on the homelab git'. Also trigger when the user gives a URL on a self-hosted git host (e.g. https://git.example.com/owner/repo/pulls/17 or /actions/runs/290), or wants to call a Forgejo REST endpoint directly ('forgejo api', 'call the Forgejo API', anything the typed commands do not cover). The orbit CLI alias is `fj`."
---

# Forgejo with orbit CLI

Manage a Forgejo (or Gitea) instance through the `orbit` CLI: repositories, pull requests and reviews, issues, Forgejo Actions runs and job logs, releases, and a passthrough to every other endpoint.

## Prerequisites

1. `orbit` CLI installed - if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
2. A profile with a `forgejo` service in `~/.config/orbit/config.yaml`. `base_url` is required - Forgejo is always self-hosted, so there is no default host.
3. A personal access token from **Settings > Applications** on the instance, configured as `method: token`. It can be an `op://` or `infisical://` reference.

Verify with `orbit -p <profile> service ping <service-name>`.

## Quick Reference

All commands follow `orbit -p <profile> forgejo <command> [flags]`, alias `fj`.

Repositories are always `owner/repo`. All list and view commands support `-o json`.

For every flag, see `references/commands.md`.

## Two things that bite

**1. A workflow run has two numbers.** The one in the URL (`/actions/runs/290`) is `index_in_repo`; the API path takes the instance-wide `id` (4023 for that same run). Asking the API for the URL number returns *a different run* with a 200 and no warning. Every `orbit fj run` command takes the number you can see and resolves it safely - so pass the URL number, and never hand-roll `api /api/v1/.../actions/runs/<url-number>`.

**2. Forgejo's issue endpoint returns pull requests too.** `orbit fj issue list` filters them out. A raw `api` call to `/issues` will not.

## Core Workflows

### Reviewing a pull request

```bash
orbit -p homelab fj pr list jorgemuza/orbit-cli
orbit -p homelab fj pr view jorgemuza/orbit-cli 17
orbit -p homelab fj pr diff jorgemuza/orbit-cli 17

# Leave the actionable feedback, then set the state
orbit -p homelab fj pr comment jorgemuza/orbit-cli 17 -m "Three things to fix, see above"
orbit -p homelab fj pr request-changes jorgemuza/orbit-cli 17

# When it is ready
orbit -p homelab fj pr approve jorgemuza/orbit-cli 17 -m "LGTM"
orbit -p homelab fj pr merge jorgemuza/orbit-cli 17 --method squash --delete-branch
```

`approve` and `request-changes` verify the state Forgejo stored, and `merge` reads the pull request back to confirm it actually landed. A success line from these means it is really so.

| Intent | Command | Effect |
|--------|---------|--------|
| Ready to merge | `pr approve` | Approving review; unblocks merge |
| Needs fixes, still the right direction | `pr request-changes` | Blocks the merge, PR stays open. Alias: `needs-work` |
| Not going to be merged at all | `pr close` | Closes without merging |

### Debugging a failed CI run

```bash
# What failed recently
orbit -p homelab fj run list jorgemuza/orbit-cli --status failure --limit 5

# The run and its jobs (the number is the one from the run's URL)
orbit -p homelab fj run view jorgemuza/orbit-cli 263

# Straight to the part that matters
orbit -p homelab fj run log jorgemuza/orbit-cli 263 --failed

# A single job, by the id `run jobs` prints
orbit -p homelab fj run log jorgemuza/orbit-cli 263 --job 5592
```

`--exit-status` on `run view` exits non-zero when the run did not succeed, for scripting a check without parsing output.

### Issues

```bash
orbit -p homelab fj issue list jorgemuza/orbit-cli --state all
orbit -p homelab fj issue view jorgemuza/orbit-cli 12
orbit -p homelab fj issue create jorgemuza/orbit-cli --title "Flaky test" --body "..." --label bug
orbit -p homelab fj issue comment jorgemuza/orbit-cli 12 -m "Reproduced on v0.68.0"
orbit -p homelab fj issue close jorgemuza/orbit-cli 12
```

### Anything else

```bash
orbit -p homelab fj api /api/v1/version
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/branches --paginate
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/tags -F limit=10
```

The `/api/v1` prefix is part of the path for `api` (the typed commands add it themselves). Flags match `orbit github api`: `-X`, `-f`, `-F`, `-H`, `--input`, `--paginate`, `-i`.

## Notes

- Gitea instances speak the same API and work with this service type.
- Forgejo reports a run's outcome in `status` - there is no separate `conclusion` as on GitHub.
- Requests are pinned to the configured host: an endpoint or redirect pointing elsewhere is refused rather than followed with your token attached.

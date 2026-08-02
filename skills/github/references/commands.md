# GitHub Commands Reference

## Repository

### `github repo [owner/repo]`
View repository details.
```
orbit github repo octocat/hello-world
```

### `github repos`
List repositories for the authenticated user.
```
orbit github repos
orbit github repos --org kubernetes --limit 10
```
Flags:
- `--org` — list repos for an organization
- `--limit` — max results (default: 30)

### `github repo edit [owner/repo]`
Edit repository settings.
```
orbit gh repo edit Paybook/ai --description "New description"
orbit gh repo edit Paybook/ai --archived              # archive repo
orbit gh repo edit Paybook/ai --archived=false         # unarchive repo
orbit gh repo edit Paybook/ai --private --default-branch develop
```
Flags:
- `--description` — new description
- `--private` — set private status
- `--default-branch` — set default branch
- `--archived` — archive or unarchive the repository

### `github repo collaborator list [owner/repo]`
List repository collaborators with permissions.
```
orbit gh repo collab list Paybook/ai
```

### `github repo collaborator add [owner/repo] [username]`
Add a collaborator to a repository.
```
orbit gh repo collab add Paybook/ai jorgemuza --permission admin
```
Flags:
- `--permission` — permission level: `pull`, `triage`, `push`, `maintain`, `admin` (default: push)

### `github repo collaborator remove [owner/repo] [username]`
Remove a collaborator from a repository.
```
orbit gh repo collab remove Paybook/ai jorgemuza
```

## Pull Requests

### `github pr list [owner/repo]`
List pull requests.
```
orbit github pr list octocat/hello-world
orbit github pr list octocat/hello-world --state closed
```
Flags:
- `--state` — filter: open, closed, all
- `--limit` — max results (default: 20)

### `github pr view [owner/repo] [number]`
View a pull request.
```
orbit github pr view octocat/hello-world 42
```

### `github pr create [owner/repo]`
Create a pull request.
```
orbit github pr create octocat/hello-world --head feature/x --base main --title "Add feature"
```
Flags:
- `--head` — head branch (required)
- `--base` — base branch (required)
- `--title` — PR title (required)
- `--body` — PR body

### `github pr merge [owner/repo] [number]`
Merge a pull request.
```
orbit github pr merge octocat/hello-world 42 --method squash
```
Flags:
- `--method` — merge method: merge, squash, rebase

### `github pr comment [owner/repo] [number]`
Add a comment to a pull request.
```
orbit github pr comment octocat/hello-world 42 --body "LGTM!"
```
Flags:
- `--body` — comment body (required)

### `github pr comments [owner/repo] [number]`
List comments on a pull request.
```
orbit github pr comments octocat/hello-world 42
```
Flags:
- `--limit` — max results (default: 50)

## Issues

### `github issue list [owner/repo]`
List issues.
```
orbit github issue list octocat/hello-world --state open --labels bug
```
Flags:
- `--state` — filter: open, closed, all
- `--labels` — filter by labels (comma-separated)
- `--limit` — max results (default: 20)

### `github issue view [owner/repo] [number]`
View an issue.
```
orbit github issue view octocat/hello-world 1
```

### `github issue create [owner/repo]`
Create an issue.
```
orbit github issue create octocat/hello-world --title "Bug report" --labels bug
```
Flags:
- `--title` — issue title (required)
- `--body` — issue body
- `--labels` — labels (comma-separated)

### `github issue close [owner/repo] [number]`
Close an issue.
```
orbit github issue close octocat/hello-world 1
```

### `github issue comment [owner/repo] [number]`
Add a comment to an issue.
```
orbit github issue comment octocat/hello-world 1 --body "Working on this"
```
Flags:
- `--body` — comment body (required)

## Branches

### `github branch list [owner/repo]`
List branches.
```
orbit github branch list octocat/hello-world
```
Flags:
- `--limit` — max results (default: 50)

### `github branch view [owner/repo] [branch]`
View a branch.
```
orbit github branch view octocat/hello-world main
```

## Tags

### `github tag list [owner/repo]`
List tags.
```
orbit github tag list octocat/hello-world
```
Flags:
- `--limit` — max results (default: 50)

## Commits

### `github commit list [owner/repo]`
List commits.
```
orbit github commit list octocat/hello-world --ref main
```
Flags:
- `--ref` — branch or tag name
- `--limit` — max results (default: 20)

### `github commit view [owner/repo] [sha]`
View a commit.
```
orbit github commit view octocat/hello-world abc1234
```

## Releases

### `github release list [owner/repo]`
List releases.
```
orbit github release list octocat/hello-world
```
Flags:
- `--limit` — max results (default: 20)

### `github release view [owner/repo] [id]`
View a release by ID.
```
orbit github release view octocat/hello-world 12345
```

### `github release latest [owner/repo]`
View the latest release.
```
orbit github release latest octocat/hello-world
```

## Workflow Runs (GitHub Actions)

### `github run list [owner/repo]`
List workflow runs.
```
orbit github run list octocat/hello-world --branch main --status completed
```
Flags:
- `--branch` — filter by branch
- `--status` — filter: completed, in_progress, queued
- `--limit` — max results (default: 20)

### `github run view [owner/repo] [run-id]`
View a workflow run, its jobs and their steps.
```
orbit github run view octocat/hello-world 12345
orbit github run view octocat/hello-world 12345 --log-failed
orbit github run view octocat/hello-world 12345 --job 67890 --log
```
Flags:
- `--job` — view a single job of the run by job ID
- `--log` — print full job logs (with `--job`, that job only)
- `--log-failed` — print the log output of the failed steps only
- `--exit-status` — exit non-zero if the run did not succeed
- `--attempt` — view a specific run attempt

`--log-failed` narrows each failed job's log to the failing step's output: step headers
matched against the step order and bounded by the step's start/end times, then the
`##[group]` title, then the log's `##[error]` markers - each chunk says which of those
located it in `matched_by`, and `step order (unconfirmed)` means nothing corroborated the
match, so that slice may not be the failing one. A step that produced no output says so,
and a log that does not cover a step reports itself truncated. If a log cannot be narrowed
down, the full job log is printed with the failed step names called out above it.

Log fetching follows GitHub's redirect to a signed storage URL. orbit refuses to follow
it when a proxy is configured for the service, an environment proxy applies to the
target, or `tls_skip_verify` is set, because it cannot then verify where the request
would go; the error names the reason. Users behind a corporate proxy cannot fetch logs
today. Other `run` commands work normally in those setups. Full list of caveats:
[docs/github.md](../../../docs/github.md#log-fetching-caveats).

### `github run jobs [owner/repo] [run-id]`
List the jobs of a workflow run (ID, status, conclusion, elapsed, name). Alias: `job`.
```
orbit github run jobs octocat/hello-world 12345
```
Flags:
- `--attempt` — list jobs of a specific run attempt

### `github run log [owner/repo] [run-id]`
Print workflow run logs as plain text. Without flags, every job's log is printed under
its own header. Alias: `logs`.
```
orbit github run log octocat/hello-world 12345 --failed
orbit github run log octocat/hello-world 12345 --job 67890
```
Flags:
- `--job` — print the log of a single job by job ID
- `--failed` — print the log output of the failed steps only
- `--attempt` — use a specific run attempt

### `github run watch [owner/repo] [run-id]`
Watch a workflow run until it completes. Polls for status updates and displays job/step progress in real-time. If no run-id is provided, watches the most recent in-progress run.
```
orbit github run watch octocat/hello-world
orbit github run watch octocat/hello-world 12345
orbit github run watch octocat/hello-world --interval 10
```
Flags:
- `--interval` — polling interval in seconds (default: 5)

Exit behavior: exits 0 on success, non-zero if the run fails, is cancelled, or times out.

### `github run cancel [owner/repo] [run-id]`
Cancel a workflow run.
```
orbit github run cancel octocat/hello-world 12345
```

### `github run rerun [owner/repo] [run-id]`
Re-run a workflow run.
```
orbit github run rerun octocat/hello-world 12345
```

## Workflows (GitHub Actions)

### `github workflow list [owner/repo]`
List workflows in a repository.
```
orbit github workflow list octocat/hello-world
```

### `github workflow view [owner/repo] [workflow-id]`
View a workflow.
```
orbit github workflow view octocat/hello-world 245836153
```

### `github workflow run [owner/repo] [workflow-id]`
Trigger a workflow dispatch event. Requires `--ref` flag.
```
orbit github workflow run octocat/hello-world 245836153 --ref main
orbit github workflow run octocat/hello-world 245836153 --ref main --input key=value
```
Flags:
- `--ref` — **(required)** git ref to dispatch on (branch or tag)
- `--input` — workflow input as `key=value` (repeatable)

**Workflow:** First `workflow list` to get the numeric workflow ID, then `workflow run` with `--ref`.

### `github workflow enable [owner/repo] [workflow-id]`
Enable a disabled workflow.
```
orbit github workflow enable octocat/hello-world 245836153
```

### `github workflow disable [owner/repo] [workflow-id]`
Disable a workflow.
```
orbit github workflow disable octocat/hello-world 245836153
```

## Secrets (GitHub Actions)

### `github secret list [owner/repo]`
List repository secrets (names and timestamps only — values are never exposed).
```
orbit github secret list octocat/hello-world
```
Flags:
- `--limit` — max results (default: 30)

### `github secret set [owner/repo] [name] [value]`
Create or update a repository secret. The value is encrypted client-side using the repository's public key before being sent to the API.
```
orbit github secret set octocat/hello-world MY_SECRET "secret-value"
orbit github secret set octocat/hello-world DEPLOY_KEY "$(cat key.pem)"
```

### `github secret delete [owner/repo] [name]`
Delete a repository secret.
```
orbit github secret delete octocat/hello-world MY_SECRET
```

## Users

### `github user me`
Show current authenticated user.
```
orbit github user me
```

### `github user view [username]`
View a user profile.
```
orbit github user view octocat
```

## Raw API

### `github api [endpoint]`
Make an authenticated request to any GitHub REST API endpoint. The endpoint is a path relative to the configured base URL, with or without a leading slash, so it also works against GitHub Enterprise. The response body is pretty-printed JSON (`-o yaml` renders YAML); a non-2xx response prints the API's error body and exits non-zero.
```
orbit github api /repos/cli/cli
orbit github api repos/cli/cli/issues --paginate -F per_page=100
orbit github api /rate_limit -i
orbit github api /repos/octocat/hello-world/issues -f title="Bug" -F body=@report.md
orbit github api /repos/octocat/hello-world/issues -f title="Bug" -F 'labels[]=bug' -F 'labels[]=p1'
orbit github api /repos/octocat/hello-world/issues/1 -X PATCH -F state=closed
orbit github api /repos/octocat/hello-world/issues --input payload.json
orbit github api /repos/cli/cli/pulls/1 -H "Accept: application/vnd.github.diff"
```
Flags:
- `-X, --method` - HTTP method (default: GET; POST when fields or `--input` are supplied)
- `-f, --raw-field` - `key=value` parameter, always sent as a string
- `-F, --field` - `key=value` parameter with type inference: `42` → number, `true`/`false` → boolean, `null` → JSON null, `@file` → the file's contents (`@-` reads stdin)
- `-H, --header` - extra request header in `name:value` format (repeatable)
- `--input` - read the raw request body from a file, or `-` for stdin (mutually exclusive with the field flags)
- `--paginate` - follow `Link rel="next"` and concatenate the JSON array of every page (GET only)
- `-i, --include` - print the response status line and headers before the body

Notes:
- Both field flags accept `key[]=value`: repeat it to build a JSON array, elements are typed by the same rules as scalars, one occurrence is still an array, and a bare `key[]` sends an empty array. Mixing `key=` and `key[]=` for one name is an error. On `GET`/`HEAD` an array repeats as `key[]=a&key[]=b`.
- On `GET`/`HEAD` the fields become query parameters instead of a request body.
- `--paginate` is read-only and fetches every page (stopping at 1000 pages, or earlier if the server's `Link` headers loop, and saying so); endpoints that return an object rather than an array are reported instead of being aggregated.
- Requests stay on the host the connection is configured for: a full URL, pagination link, or redirect pointing elsewhere is refused rather than followed with your token attached.
- Any non-2xx status, including `304 Not Modified` and unfollowed redirects, exits non-zero.

## Global Flags

All commands inherit these flags:
- `-o, --output` — output format: table, json, yaml (default: table)
- `-p, --profile` — profile to use
- `--service` — github service name (if profile has multiple)

# Orbit GitHub Command Reference

Interact with GitHub repositories, pull requests, issues, actions, and more from the command line.

**Top-level command:** `orbit github` (alias: `gh`)

**Persistent flag (all subcommands):**

| Flag | Description |
|------|-------------|
| `--service` | GitHub service name, required only when a profile has multiple GitHub services configured |

**Notes:**

- The default API base URL is `https://api.github.com`.
- For GitHub Enterprise, set a custom base URL with `--base-url`.
- Repository arguments always use the `owner/repo` format (e.g., `octocat/hello-world`).

---

## Table of Contents

- [repo — View a repository](#repo)
- [repos — List repositories](#repos)
- [branch — Branch commands](#branch)
  - [branch list](#branch-list)
  - [branch view](#branch-view)
- [tag — Tag commands](#tag)
  - [tag list](#tag-list)
- [commit — Commit commands](#commit)
  - [commit list](#commit-list)
  - [commit view](#commit-view)
- [pr — Pull request commands](#pr)
  - [pr list](#pr-list)
  - [pr view](#pr-view)
  - [pr create](#pr-create)
  - [pr merge](#pr-merge)
  - [pr comment](#pr-comment)
  - [pr comments](#pr-comments)
- [issue — Issue commands](#issue)
  - [issue list](#issue-list)
  - [issue view](#issue-view)
  - [issue create](#issue-create)
  - [issue comment](#issue-comment)
  - [issue close](#issue-close)
- [release — Release commands](#release)
  - [release list](#release-list)
  - [release view](#release-view)
  - [release latest](#release-latest)
- [run — Workflow run commands](#run)
  - [run list](#run-list)
  - [run view](#run-view)
  - [run jobs](#run-jobs)
  - [run log](#run-log)
  - [Log fetching caveats](#log-fetching-caveats)
  - [run rerun](#run-rerun)
  - [run cancel](#run-cancel)
  - [run watch](#run-watch)
- [secret — Repository secret commands](#secret)
  - [secret list](#secret-list)
  - [secret set](#secret-set)
  - [secret delete](#secret-delete)
- [user — User commands](#user)
  - [user me](#user-me)
  - [user view](#user-view)
- [api - Call any GitHub API endpoint](#api)

---

## repo

View details of a repository.

```
orbit github repo [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Repository to view |

**Example:**

```bash
orbit github repo octocat/hello-world -p myprofile
```

---

## repos

List repositories accessible to the authenticated user.

```
orbit github repos [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--search` | | Filter repositories by search query |
| `--limit` | 50 | Maximum number of repositories to return |

**Examples:**

```bash
# List repositories
orbit github repos -p myprofile

# Search for repositories matching a keyword
orbit github repos --search "api" -p myprofile

# Limit results
orbit github repos --limit 10 -p myprofile
```

---

## repo edit

Edit repository settings.

```
orbit github repo edit [owner/repo] [flags] -p myprofile
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--description` | string | New repository description. |
| `--private` | bool | Set private status. |
| `--default-branch` | string | Set default branch. |
| `--archived` | bool | Archive or unarchive the repository. |

**Examples:**

```bash
# Archive a repository
orbit gh repo edit Paybook/ai --archived -p myprofile

# Unarchive a repository
orbit gh repo edit Paybook/ai --archived=false -p myprofile

# Update description
orbit gh repo edit Paybook/ai --description "Updated description" -p myprofile
```

---

## repo collaborator

Manage repository collaborators (alias: `collab`).

### repo collaborator list

List repository collaborators with permissions.

```bash
orbit gh repo collab list Paybook/ai -p myprofile
```

### repo collaborator add

Add a collaborator to a repository.

```
orbit github repo collaborator add [owner/repo] [username] [flags] -p myprofile
```

| Flag | Default | Description |
|------|---------|-------------|
| `--permission` | `push` | Permission: `pull`, `triage`, `push`, `maintain`, `admin`. |

```bash
orbit gh repo collab add Paybook/ai jorgemuza --permission admin -p myprofile
```

### repo collaborator remove

Remove a collaborator from a repository.

```bash
orbit gh repo collab remove Paybook/ai jorgemuza -p myprofile
```

---

## branch

Manage repository branches.

### branch list

List branches in a repository.

**Aliases:** `ls`

```
orbit github branch list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of branches to return |

**Examples:**

```bash
orbit github branch list octocat/hello-world -p myprofile
orbit github branch ls octocat/hello-world --limit 20 -p myprofile
```

### branch view

View details of a specific branch.

```
orbit github branch view [owner/repo] [branch] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `branch` | Branch name to view |

**Example:**

```bash
orbit github branch view octocat/hello-world main -p myprofile
```

---

## tag

Manage repository tags.

### tag list

List tags in a repository.

**Aliases:** `ls`

```
orbit github tag list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of tags to return |

**Example:**

```bash
orbit github tag list octocat/hello-world --limit 10 -p myprofile
```

---

## commit

Manage repository commits.

### commit list

List commits in a repository.

**Aliases:** `ls`

```
orbit github commit list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of commits to return |

**Example:**

```bash
orbit github commit list octocat/hello-world --limit 20 -p myprofile
```

### commit view

View details of a specific commit.

```
orbit github commit view [owner/repo] [sha] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `sha` | Full or abbreviated commit SHA |

**Example:**

```bash
orbit github commit view octocat/hello-world abc1234 -p myprofile
```

---

## pr

Manage pull requests.

**Aliases:** `pull-request`

### pr list

List pull requests in a repository.

**Aliases:** `ls`

```
orbit github pr list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | | Filter by state: `open`, `closed`, or `all` |
| `--limit` | | Maximum number of pull requests to return |

**Examples:**

```bash
orbit github pr list octocat/hello-world -p myprofile
orbit github pr ls octocat/hello-world --state closed --limit 5 -p myprofile
```

### pr view

View details of a specific pull request.

```
orbit github pr view [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

**Example:**

```bash
orbit github pr view octocat/hello-world 42 -p myprofile
```

### pr create

Create a new pull request.

```
orbit github pr create [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Required | Description |
|------|----------|-------------|
| `--head` | Yes | Source branch name |
| `--base` | Yes | Target branch name |
| `--title` | Yes | Pull request title |
| `--body` | No | Pull request description |

**Examples:**

```bash
# Minimal create
orbit github pr create octocat/hello-world \
  --head feature/login \
  --base main \
  --title "Add login page" \
  -p myprofile

# With body
orbit github pr create octocat/hello-world \
  --head feature/login \
  --base main \
  --title "Add login page" \
  --body "Implements the new login flow with OAuth support." \
  -p myprofile
```

### pr merge

Merge a pull request.

```
orbit github pr merge [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Default | Description |
|------|---------|-------------|
| `--squash` | | Squash commits before merging |

**Examples:**

```bash
orbit github pr merge octocat/hello-world 42 -p myprofile
orbit github pr merge octocat/hello-world 42 --squash -p myprofile
```

### pr comment

Add a comment to a pull request.

```
orbit github pr comment [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Required | Description |
|------|----------|-------------|
| `--body` | Yes | Comment text |

**Example:**

```bash
orbit github pr comment octocat/hello-world 42 \
  --body "LGTM, approving." \
  -p myprofile
```

### pr comments

List comments on a pull request.

```
orbit github pr comments [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Pull request number |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of comments to return |

**Example:**

```bash
orbit github pr comments octocat/hello-world 42 --limit 10 -p myprofile
```

---

## issue

Manage repository issues.

### issue list

List issues in a repository.

**Aliases:** `ls`

```
orbit github issue list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | | Filter by state: `open`, `closed`, or `all` |
| `--labels` | | Filter by comma-separated label names |
| `--limit` | | Maximum number of issues to return |

**Examples:**

```bash
orbit github issue list octocat/hello-world -p myprofile
orbit github issue ls octocat/hello-world --state open --labels "bug,critical" --limit 25 -p myprofile
```

### issue view

View details of a specific issue.

```
orbit github issue view [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

**Example:**

```bash
orbit github issue view octocat/hello-world 7 -p myprofile
```

### issue create

Create a new issue.

```
orbit github issue create [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Issue title |
| `--body` | No | Issue description |
| `--labels` | No | Comma-separated label names |

**Examples:**

```bash
orbit github issue create octocat/hello-world \
  --title "Fix broken login" \
  -p myprofile

orbit github issue create octocat/hello-world \
  --title "Fix broken login" \
  --body "The login form returns a 500 error when submitting valid credentials." \
  --labels "bug,high-priority" \
  -p myprofile
```

### issue comment

Add a comment to an issue.

```
orbit github issue comment [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

| Flag | Required | Description |
|------|----------|-------------|
| `--body` | Yes | Comment text |

**Example:**

```bash
orbit github issue comment octocat/hello-world 7 \
  --body "Reproduced on latest main. Investigating." \
  -p myprofile
```

### issue close

Close an issue.

```
orbit github issue close [owner/repo] [number] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `number` | Issue number |

**Example:**

```bash
orbit github issue close octocat/hello-world 7 -p myprofile
```

---

## release

Manage repository releases.

### release list

List releases in a repository.

**Aliases:** `ls`

```
orbit github release list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of releases to return |

**Example:**

```bash
orbit github release list octocat/hello-world --limit 5 -p myprofile
```

### release view

View details of a specific release by tag.

```
orbit github release view [owner/repo] [tag] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `tag` | Release tag name |

**Example:**

```bash
orbit github release view octocat/hello-world v1.2.0 -p myprofile
```

### release latest

View the latest release of a repository.

```
orbit github release latest [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

**Example:**

```bash
orbit github release latest octocat/hello-world -p myprofile
```

---

## run

Manage GitHub Actions workflow runs.

**Aliases:** `actions`

### run list

List workflow runs in a repository.

**Aliases:** `ls`

```
orbit github run list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | | Maximum number of runs to return |

**Example:**

```bash
orbit github run list octocat/hello-world --limit 10 -p myprofile
```

### run view

View details of a specific workflow run, its jobs and their steps. A failing run
shows which step failed, and `--log-failed` prints that step's log output in one
command.

```
orbit github run view [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

| Flag | Default | Description |
|------|---------|-------------|
| `--job` | | View a single job of the run by job ID |
| `--log` | `false` | Print full job logs (with `--job`, that job only) |
| `--log-failed` | `false` | Print the log output of the failed steps only |
| `--exit-status` | `false` | Exit with a non-zero status if the run did not succeed |
| `--attempt` | | View a specific run attempt |

**Examples:**

```bash
# Run metadata plus every job and step
orbit github run view octocat/hello-world 123456789 -p myprofile

# Which step failed, and what it printed
orbit github run view octocat/hello-world 123456789 --log-failed -p myprofile

# A single job, with its full log
orbit github run view octocat/hello-world 123456789 --job 987654321 --log -p myprofile

# Use in scripts: non-zero exit when the run did not succeed
orbit github run view octocat/hello-world 123456789 --exit-status -p myprofile
```

`--log-failed` narrows each failed job's log to the output of the steps that
failed. Steps are placed by matching the runner's `##[group]Run ...` headers
against the step sequence, bounded by each step's start and end times, which
handles a step that emits several groups, two steps sharing a display name, and
two steps running inside the same second. A step that failed without
producing any output is reported as such rather than being given a neighbour's.
If a step cannot be placed at all it is matched by its `##[group]` title, and
failing that by the `##[error]` markers in the log; when a log cannot be narrowed
down at all, the full job log is printed with the failed step names called out
above it. Every chunk says which of those located it, so the output never claims
more precision than it has.

See [Log fetching caveats](#log-fetching-caveats) for when log fetching is refused,
and for how far to trust `--log-failed`.

### run jobs

List the jobs of a workflow run, with their status, conclusion and elapsed time.

**Aliases:** `job`

```
orbit github run jobs [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

| Flag | Default | Description |
|------|---------|-------------|
| `--attempt` | | List jobs of a specific run attempt |

**Example:**

```bash
orbit github run jobs octocat/hello-world 123456789 -p myprofile
```

### run log

Print workflow run logs as plain text. Without flags every job's log is printed
under its own header.

**Aliases:** `logs`

```
orbit github run log [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

| Flag | Default | Description |
|------|---------|-------------|
| `--job` | | Print the log of a single job by job ID |
| `--failed` | `false` | Print the log output of the failed steps only |
| `--attempt` | | Use a specific run attempt |

**Examples:**

```bash
orbit github run log octocat/hello-world 123456789 --failed -p myprofile
orbit github run log octocat/hello-world 123456789 --job 987654321 -p myprofile
```

### Log fetching caveats

These apply to the commands that fetch log text: `run view --log`, `run view --log-failed`
and `run log`. The rest of `run` is unaffected.

**Log fetching is refused behind a proxy, and with `tls_skip_verify` - by design.**
GitHub serves job logs by redirecting to a signed URL on a storage host. orbit follows
that redirect only when it can verify where the request actually goes. A proxy resolves
and routes the target out of orbit's sight, and with certificate verification disabled
the host that answers never has to prove it is the host the redirect named; in both
cases the request can land on an internal service whose contents orbit would then print.
Refusing is the safe answer, so orbit refuses and says why:

```
refusing to follow the log redirect to productionresultssa12.blob.core.windows.net:
a proxy is configured for this connection (socks5h://127.0.0.1:1082), so orbit cannot
verify where the request would actually be routed [...] Remove the proxy for this
service, or fetch the log from the web UI
```

This applies when the service (or its profile) sets `proxy`, when `HTTP_PROXY` or
`HTTPS_PROXY` applies to the log host (`NO_PROXY` is honoured, using Go's own rules), or
when `tls_skip_verify` is set. **Users behind a corporate proxy cannot fetch logs with
orbit today.** What still works: `run list`, `run view`, `run jobs` and `run watch` all
run normally through a proxy, so the run, its jobs, and which step failed are still
visible - and `run view` prints the run's `URL` for reading the log in the browser.

**`--log-failed` attribution is best-effort, and reports its own confidence.** Each chunk
says how it was located, as `matched_by` in JSON and as a note in the text output:

| `matched_by` | What it means for trust |
|--------------|-------------------------|
| `step order` | The step's own section, matched by the runner's `##[group]Run ...` headers in step order and bounded by the step's start and end times. Reliable. |
| `error marker` | Located by the `##[error]` line the section ends with, because step order could not place it. The section really is the failure; which step it is attached to is inferred. |
| `step order (unconfirmed)` | The slice is the failed step's by order, but it carries no `##[error]` marker, so nothing corroborates it. **The output shown may not be the failing output.** |
| `step name` | No usable timestamps, so the section was matched by its `##[group]` title. |

Attribution degrades when the API's step data is thin, or disagrees with the log:
duplicate or missing step numbers, several steps inside one rounded second (step times
come back rounded to the second), an API step order that does not match the order of the
log, and steps that produce no output of their own. When nothing can be narrowed down at
all, the whole job log is printed with the failed step names called out above it. A step
that genuinely produced nothing is reported as such rather than being given a
neighbour's output.

**Truncated and expired logs are reported, not guessed at.** GitHub cuts and eventually
expires job logs. When the log does not cover a failed step, that chunk says
`(the job log does not cover this step - it was truncated, or has expired)` and carries
`truncated: true` in JSON, instead of appearing empty or borrowing a neighbouring step's
output.

**Not verified.** The proxy refusal policy has been exercised against a local `socks5h`
proxy and stubbed resolvers only - not against a real GitHub Enterprise install behind a
corporate proxy with split-horizon DNS and certificate verification enabled. Job
pagination past the API's 100-per-page limit is proven with a test server rather than
against a real matrix run of more than 100 jobs.

### run rerun

Re-run a workflow run.

```
orbit github run rerun [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run rerun octocat/hello-world 123456789 -p myprofile
```

### run cancel

Cancel an in-progress workflow run.

```
orbit github run cancel [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run cancel octocat/hello-world 123456789 -p myprofile
```

### run watch

Watch a workflow run until it completes.

```
orbit github run watch [owner/repo] [run-id] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `run-id` | Workflow run ID |

**Example:**

```bash
orbit github run watch octocat/hello-world 123456789 -p myprofile
```

---

## secret

Manage repository secrets.

### secret list

List secrets configured on a repository.

**Aliases:** `ls`

```
orbit github secret list [owner/repo] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |

**Example:**

```bash
orbit github secret list octocat/hello-world -p myprofile
```

### secret set

Create or update a repository secret.

```
orbit github secret set [owner/repo] [secret-name] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `secret-name` | Name of the secret |

| Flag | Required | Description |
|------|----------|-------------|
| `--value` | Yes | Secret value |

**Example:**

```bash
orbit github secret set octocat/hello-world API_KEY \
  --value "sk-abc123" \
  -p myprofile
```

### secret delete

Delete a repository secret.

```
orbit github secret delete [owner/repo] [secret-name] [flags]
```

| Argument | Description |
|----------|-------------|
| `owner/repo` | Target repository |
| `secret-name` | Name of the secret to delete |

**Example:**

```bash
orbit github secret delete octocat/hello-world API_KEY -p myprofile
```

---

## user

View GitHub user information.

### user me

Display the currently authenticated user.

```
orbit github user me [flags]
```

**Example:**

```bash
orbit github user me -p myprofile
```

### user view

View a GitHub user's profile.

```
orbit github user view [username] [flags]
```

| Argument | Description |
|----------|-------------|
| `username` | GitHub username to look up |

**Example:**

```bash
orbit github user view octocat -p myprofile
```

---

## api

Make an authenticated request to any GitHub REST API endpoint. Use it for
anything the typed commands do not cover.

```
orbit github api [endpoint] [flags]
```

| Argument | Description |
|----------|-------------|
| `endpoint` | API path relative to the configured base URL, with or without a leading slash (`/repos/cli/cli` or `repos/cli/cli`). Works unchanged against GitHub Enterprise, whose base URL usually ends in `/api/v3`. |

| Flag | Default | Description |
|------|---------|-------------|
| `-X`, `--method` | `GET` | HTTP method. Defaults to `POST` when any field flag or `--input` is supplied |
| `-f`, `--raw-field` | | Add a parameter in `key=value` format; the value is always a string |
| `-F`, `--field` | | Add a parameter in `key=value` format with type inference (see below) |
| | | Both field flags also accept the `key[]=value` array form (see below) |
| `-H`, `--header` | | Add a request header in `name:value` format; repeatable |
| `--input` | | Read the raw request body from a file, or from stdin with `-`. Mutually exclusive with the field flags |
| `--paginate` | `false` | Follow `Link rel="next"` and concatenate the JSON array of every page (GET only) |
| `-i`, `--include` | `false` | Print the response status line and headers before the body |

**Field types (`-F`):**

| Value | Sent as |
|-------|---------|
| `42` | number |
| `true` / `false` | boolean |
| `null` | JSON null |
| `@path/to/file` | the file's contents, as a string (`@-` reads stdin) |
| anything else | string |

`-f` skips the inference entirely, so `-f count=42` sends the string `"42"`.

**Array fields (`key[]=value`):**

Repeat a field with the `key[]` suffix to build a JSON array. Elements keep the
same typing rules as scalars, so `-F` infers each element and `-f` keeps them
all strings:

| Flags | Sent as |
|-------|---------|
| `-F 'labels[]=bug' -F 'labels[]=p1'` | `{"labels":["bug","p1"]}` |
| `-F 'ids[]=1' -F 'ids[]=2'` | `{"ids":[1,2]}` |
| `-f 'ids[]=1' -f 'ids[]=2'` | `{"ids":["1","2"]}` |
| `-F 'labels[]=only'` | `{"labels":["only"]}` (one element is still an array) |
| `-F 'labels[]'` | `{"labels":[]}` (empty array) |

Using both forms for one name (`-F labels=x -F 'labels[]=y'`) is an error. On
`GET`/`HEAD` an array field repeats in the query string as
`labels[]=bug&labels[]=p1`.

**Notes:**

- For `GET` and `HEAD` requests the fields become query parameters instead of a
  request body; any query already on the endpoint is preserved.
- Requests never leave the host the connection is configured for. An endpoint
  given as a full URL, a pagination link, or a redirect that points at another
  host is refused with a clear error rather than followed, because the request
  would carry your token. Same-host redirects are followed normally.
- The response body is pretty-printed as JSON, which keeps the API's key order.
  `-o yaml` renders it as YAML; a non-JSON body (a `.diff`, for example) is
  printed unchanged.
- Any non-2xx response prints the body and exits non-zero, including a `304 Not
  Modified` from a conditional request and any redirect that was not followed.
- `--paginate` is read-only: fields stay query parameters on a `GET`, and
  combining it with an explicit write method is rejected. It fetches every page,
  so pass a large `per_page` on big collections.
- `--paginate` stops at 1000 pages and says so instead of truncating silently,
  and stops if the server's `Link` headers ever point back at a page already
  fetched.
- `--paginate` aggregates JSON arrays. Endpoints that wrap their results in an
  object (`/search/*`, `/repos/{o}/{r}/actions/runs`) are reported as such
  instead of being silently mangled.

**Examples:**

```bash
# Any GET endpoint
orbit github api /repos/cli/cli -p myprofile

# Every page of a collection, five at a time
orbit github api /repos/cli/cli/issues --paginate -F per_page=5 -p myprofile

# Inspect status and headers (rate limits, ETag, deprecations)
orbit github api /repos/cli/cli -i -p myprofile

# Create an issue: method defaults to POST because fields were supplied
orbit github api /repos/octocat/hello-world/issues \
  -f title="Bug report" -F body=@report.md -p myprofile

# Array fields: repeat key[] to build a JSON array
orbit github api /repos/octocat/hello-world/issues \
  -f title="Bug report" -F 'labels[]=bug' -F 'labels[]=priority-1' -p myprofile

# Close an issue with an explicit method
orbit github api /repos/octocat/hello-world/issues/1 -X PATCH -F state=closed -p myprofile

# Send a hand-written body
orbit github api /repos/octocat/hello-world/issues --input payload.json -p myprofile

# Ask for a media type the typed commands do not expose
orbit github api /repos/cli/cli/pulls/1 -H "Accept: application/vnd.github.diff" -p myprofile
```

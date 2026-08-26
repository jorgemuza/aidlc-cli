# Orbit Forgejo Command Reference

Interact with Forgejo (and Gitea) repositories, pull requests, issues, Actions runs and releases from the command line.

**Top-level command:** `orbit forgejo` (alias: `fj`)

**Persistent flag (all subcommands):**

| Flag | Description |
|------|-------------|
| `--service` | Forgejo service name, required only when a profile has multiple Forgejo services configured |

**Notes:**

- Every Forgejo is somebody's server, so `base_url` is **required** in the service config - there is no hosted default to fall back to.
- The API lives under `/api/v1`. The typed commands add that prefix themselves; `forgejo api` does not, so pass it there.
- Repository arguments always use the `owner/repo` format (e.g. `jorgemuza/orbit-cli`).
- All list and view commands support `-o json` and `-o yaml`.
- Gitea instances speak the same API and work with this service type.

---

## Configuration

```yaml
profiles:
  - name: homelab
    services:
      - name: forgejo-homelab
        type: forgejo
        base_url: https://git.example.com
        auth:
          method: token
          token: infisical://prod/fleet/orbit/homelab/FORGEJO_TOKEN
```

Create the token on the instance under **Settings > Applications > Generate New Token**. Forgejo documents the header as `token <PAT>`; it also accepts `Bearer`, which is what orbit sends, so `method: token` is correct and needs no special casing.

Check it with:

```bash
orbit -p homelab service ping forgejo-homelab
# forgejo-homelab      OK    Forgejo authenticated as jorgemuza (Jorge Muza)
```

---

## Table of Contents

- [repo](#repo)
  - [repo view](#repo-view)
  - [repo list](#repo-list)
- [pr](#pr)
  - [pr list](#pr-list)
  - [pr view](#pr-view)
  - [pr create](#pr-create)
  - [pr merge](#pr-merge)
  - [pr diff](#pr-diff)
  - [pr comment](#pr-comment)
  - [pr comments](#pr-comments)
  - [pr approve](#pr-approve)
  - [pr request-changes](#pr-request-changes)
  - [pr reviews](#pr-reviews)
  - [pr close](#pr-close)
- [issue](#issue)
- [run](#run)
  - [run list](#run-list)
  - [run view](#run-view)
  - [run jobs](#run-jobs)
  - [run log](#run-log)
- [release](#release)
- [api](#api)

---

## repo

### `repo view`

View a repository.

```bash
orbit -p homelab fj repo view jorgemuza/orbit-cli
```

### `repo list`

List the repositories the token can see.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 30 | Max results |

```bash
orbit -p homelab fj repo list --limit 50
```

---

## pr

Pull requests are addressed by the number shown in the web UI, which is also the number the API uses for them.

### `pr list`

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | `open` | `open`, `closed` or `all` |
| `--limit <n>` | 20 | Max results |

```bash
orbit -p homelab fj pr list jorgemuza/orbit-cli
orbit -p homelab fj pr list jorgemuza/orbit-cli --state all --limit 50
```

### `pr view`

```bash
orbit -p homelab fj pr view jorgemuza/orbit-cli 17
```

### `pr create`

| Flag | Required | Description |
|------|----------|-------------|
| `--from <branch>` | Yes | Source branch |
| `--to <branch>` | Yes | Target branch |
| `--title <text>` | Yes | Title |
| `--body <text>` | No | Description |

```bash
orbit -p homelab fj pr create jorgemuza/orbit-cli --from feature/x --to main --title "Add x"
```

### `pr merge`

| Flag | Default | Description |
|------|---------|-------------|
| `--method <method>` | `merge` | `merge`, `squash`, `rebase`, `rebase-merge`, `fast-forward-only` |
| `--delete-branch` | `false` | Delete the source branch after merging |

```bash
orbit -p homelab fj pr merge jorgemuza/orbit-cli 17 --method squash --delete-branch
```

> Forgejo answers a merge with an empty `200`, which reports that the request was accepted and **not** that the branch was merged. This command reads the pull request back and fails if it is not merged, so a merge that a branch protection rule or a failing required check refused is an error rather than a success line.

### `pr diff`

Prints the unified diff, straight from `/pulls/{n}.diff`.

```bash
orbit -p homelab fj pr diff jorgemuza/orbit-cli 17
```

### `pr comment`

| Flag | Required | Description |
|------|----------|-------------|
| `-m`, `--message <text>` | Yes | Comment body |

```bash
orbit -p homelab fj pr comment jorgemuza/orbit-cli 17 -m "Please address the items above"
```

### `pr comments`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `pr approve`

Submit an `APPROVED` review.

```bash
orbit -p homelab fj pr approve jorgemuza/orbit-cli 17 -m "LGTM"
```

### `pr request-changes`

Submit a `REQUEST_CHANGES` review, which blocks the merge without closing the pull request. **Alias:** `needs-work`.

```bash
orbit -p homelab fj pr comment jorgemuza/orbit-cli 17 -m "Three things to fix"
orbit -p homelab fj pr request-changes jorgemuza/orbit-cli 17
```

> Both review commands check the state Forgejo actually stored before reporting success. A review's state is what blocks or unblocks a merge, so a `2xx` alone is not taken as proof it was recorded.

### `pr reviews`

List the reviews on a pull request, with their state and whether they have gone stale.

### `pr close`

Close a pull request without merging it.

---

## issue

`issue list`, `view`, `create`, `comment`, `comments`, `close` and `reopen` mirror the pull request commands.

| Command | Flags |
|---------|-------|
| `issue list <owner/repo>` | `--state` (default `open`), `--limit` (default 20) |
| `issue view <owner/repo> <n>` | - |
| `issue create <owner/repo>` | `--title` (required), `--body`, `--label` (repeatable), `--assignee` (repeatable) |
| `issue comment <owner/repo> <n>` | `-m`, `--message` (required) |
| `issue comments <owner/repo> <n>` | `--limit` (default 50) |
| `issue close` / `issue reopen` | - |

```bash
orbit -p homelab fj issue list jorgemuza/orbit-cli --state all
orbit -p homelab fj issue create jorgemuza/orbit-cli --title "Flaky test" --label bug
orbit -p homelab fj issue comment jorgemuza/orbit-cli 12 -m "Reproduced on v0.68.0"
```

> `issue list` sends `type=issues`. Forgejo's issue endpoint returns pull requests as well unless it is told which it wants, and a pull request listed as an issue is the kind of wrong that is easy not to notice.

---

## run

Forgejo Actions workflow runs.

### Run numbers are not run ids

**A run has two integers and they are not interchangeable.** `index_in_repo` is the number the web UI shows - the `290` in `/actions/runs/290` - while `id` is the instance-wide database id the API path takes.

Passing the number you can see to the API does not fail. It returns **a different run**, with a `200` and nothing to suggest anything is wrong:

```
GET /actions/runs/4023  ->  id 4023, index_in_repo 290, "security-audit"
GET /actions/runs/290   ->  id 290,  index_in_repo 159, a different workflow entirely
```

Every `run` command here takes the number you can see, resolves it through the `run_number` filter, and checks that the run it got back is the one that was asked for before using its id for anything.

### `run list`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 20 | Max results |
| `--status <status>` | - | `success`, `failure`, `cancelled`, `skipped`, `running`, `waiting` |
| `--event <event>` | - | `push`, `pull_request`, `schedule`, `workflow_dispatch` |
| `--branch <ref>` | - | Filter by git ref |

```bash
orbit -p homelab fj run list jorgemuza/orbit-cli
orbit -p homelab fj run list jorgemuza/orbit-cli --status failure --limit 10
```

Forgejo reports the outcome in `status`; there is no separate `conclusion` field as on GitHub.

### `run view`

Shows the run and its jobs.

| Flag | Default | Description |
|------|---------|-------------|
| `--exit-status` | `false` | Exit non-zero if the run did not succeed |

```bash
orbit -p homelab fj run view jorgemuza/orbit-cli 290
orbit -p homelab fj run view jorgemuza/orbit-cli 290 --exit-status
```

### `run jobs`

List the jobs of a run with their ids, for use with `run log --job`.

### `run log`

| Flag | Default | Description |
|------|---------|-------------|
| `--failed` | `false` | Only the jobs that failed or were cancelled |
| `--job <id>` | - | A single job, by the id `run jobs` prints |

```bash
orbit -p homelab fj run log jorgemuza/orbit-cli 263 --failed
orbit -p homelab fj run log jorgemuza/orbit-cli 290 --job 7555
```

`--failed` is usually the only part worth reading, and is a great deal shorter than the whole run. On a run with nothing failed it says so rather than printing nothing.

---

## release

| Command | Flags |
|---------|-------|
| `release list <owner/repo>` | `--limit` (default 20) |
| `release view <owner/repo> <tag>` | - |

```bash
orbit -p homelab fj release list jorgemuza/orbit-cli
orbit -p homelab fj release view jorgemuza/orbit-cli v0.68.0
```

---

## api

Make an authenticated request to any Forgejo endpoint - around 300 of them on a current instance, well beyond what the typed commands cover.

This is the same passthrough as `orbit github api` (they share one implementation), so the flags are identical: `-X/--method`, `-f/--raw-field`, `-F/--field`, `-H/--header`, `--input`, `--paginate` and `-i/--include`.

**The `/api/v1` prefix is part of the path here.**

```bash
orbit -p homelab fj api /api/v1/version
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/issues --paginate -F limit=50
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/issues -f title="Bug" -f body="..."
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/issues/12 -X PATCH -f state=closed
```

`--paginate` follows the `Link rel="next"` header, which Forgejo sends on its paginated list endpoints. The endpoints that answer with an object rather than an array - workflow runs among them - are read a page at a time with `limit` and `page` instead.

Requests never leave the host the connection is configured for: an endpoint given as a full URL, a pagination link, or a redirect that points elsewhere is refused rather than followed, because the request would carry your token.

```
$ orbit -p homelab fj api https://api.github.com/user
Error: refusing to send forgejo credentials to https://api.github.com:
service "forgejo-homelab" is configured for https://git.example.com
```

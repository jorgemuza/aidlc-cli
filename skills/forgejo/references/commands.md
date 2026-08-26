# Forgejo CLI Command Reference

Complete reference for all `orbit forgejo` (alias: `fj`) commands and flags.

Every command takes `-p <profile>` and supports `-o json` / `-o yaml`. Repositories are always `owner/repo`.

## Table of Contents

- [repo](#repo)
- [pr (pull request)](#pr-pull-request)
- [issue](#issue)
- [run (actions)](#run-actions)
- [release](#release)
- [api](#api)

---

## repo

### `forgejo repo view <owner/repo>`

View a repository: default branch, visibility, open issue and PR counts, clone URLs.

### `forgejo repo list`

List the repositories the token can see. **Alias:** `ls`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 30 | Max results |

---

## pr (pull request)

### `forgejo pr list <owner/repo>`

**Alias:** `ls`

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | `open` | `open`, `closed`, `all` |
| `--limit <n>` | 20 | Max results |

The STATE column folds Forgejo's separate flags into one word: `merged`, `draft`, `open` or `closed`.

### `forgejo pr view <owner/repo> <number>`

### `forgejo pr create <owner/repo>`

| Flag | Required | Description |
|------|----------|-------------|
| `--from <branch>` | Yes | Source branch |
| `--to <branch>` | Yes | Target branch |
| `--title <text>` | Yes | Title |
| `--body <text>` | No | Description |

### `forgejo pr merge <owner/repo> <number>`

| Flag | Default | Description |
|------|---------|-------------|
| `--method <method>` | `merge` | `merge`, `squash`, `rebase`, `rebase-merge`, `fast-forward-only` |
| `--delete-branch` | `false` | Delete the source branch after merging |

Reads the pull request back and fails if it did not actually merge - Forgejo's empty `200` says the request was accepted, not that the branch landed.

### `forgejo pr diff <owner/repo> <number>`

Prints the unified diff.

### `forgejo pr comment <owner/repo> <number>`

| Flag | Required | Description |
|------|----------|-------------|
| `-m`, `--message <text>` | Yes | Comment body |

### `forgejo pr comments <owner/repo> <number>`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `forgejo pr approve <owner/repo> <number>`

Submit an `APPROVED` review. `-m` adds a body.

### `forgejo pr request-changes <owner/repo> <number>`

Submit a `REQUEST_CHANGES` review - blocks the merge, leaves the PR open. **Alias:** `needs-work`. `-m` adds a body.

Both review commands verify the state Forgejo stored before reporting success.

### `forgejo pr reviews <owner/repo> <number>`

List reviews with reviewer, state, staleness and submission time.

### `forgejo pr close <owner/repo> <number>`

Close without merging.

---

## issue

`issue list` sends `type=issues`, so pull requests are excluded. A raw `api` call to `/issues` returns both.

### `forgejo issue list <owner/repo>`

**Alias:** `ls`

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | `open` | `open`, `closed`, `all` |
| `--limit <n>` | 20 | Max results |

### `forgejo issue view <owner/repo> <number>`

### `forgejo issue create <owner/repo>`

| Flag | Required | Description |
|------|----------|-------------|
| `--title <text>` | Yes | Issue title |
| `--body <text>` | No | Description |
| `--label <name>` | No | Label to apply (repeatable) |
| `--assignee <user>` | No | Username to assign (repeatable) |

### `forgejo issue comment <owner/repo> <number>`

| Flag | Required | Description |
|------|----------|-------------|
| `-m`, `--message <text>` | Yes | Comment body |

### `forgejo issue comments <owner/repo> <number>`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

### `forgejo issue close <owner/repo> <number>` / `forgejo issue reopen <owner/repo> <number>`

Both verify the state Forgejo stored.

---

## run (actions)

**Runs are addressed by the number in the web UI**, the `290` in `/actions/runs/290` - not by the internal API id. The two are different integers for the same run, and asking the API for the wrong one answers `200` with a different run. These commands resolve the visible number and verify they got the right run.

### `forgejo run list <owner/repo>`

**Alias:** `ls`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 20 | Max results |
| `--status <status>` | - | `success`, `failure`, `cancelled`, `skipped`, `running`, `waiting` |
| `--event <event>` | - | `push`, `pull_request`, `schedule`, `workflow_dispatch` |
| `--branch <ref>` | - | Filter by git ref |

Forgejo reports the outcome in `status`; there is no separate `conclusion`.

### `forgejo run view <owner/repo> <run-number>`

Shows the run and its jobs.

| Flag | Default | Description |
|------|---------|-------------|
| `--exit-status` | `false` | Exit non-zero if the run did not succeed |

### `forgejo run jobs <owner/repo> <run-number>`

List the jobs with their ids, for `run log --job`.

### `forgejo run log <owner/repo> <run-number>`

| Flag | Default | Description |
|------|---------|-------------|
| `--failed` | `false` | Only jobs that failed or were cancelled |
| `--job <id>` | - | A single job by id |

---

## release

### `forgejo release list <owner/repo>`

**Alias:** `ls`

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 20 | Max results |

### `forgejo release view <owner/repo> <tag>`

Shows the release and its assets.

---

## api

### `forgejo api <endpoint>`

Authenticated request to any Forgejo endpoint. **The `/api/v1` prefix is part of the path** (typed commands add it themselves).

| Flag | Default | Description |
|------|---------|-------------|
| `-X`, `--method <verb>` | `GET` | Defaults to POST when fields or `--input` are supplied |
| `-f`, `--raw-field <k=v>` | - | String parameter (`k[]=v` builds an array) |
| `-F`, `--field <k=v>` | - | Typed parameter: numbers, `true`/`false`/`null`, `@file` |
| `-H`, `--header <n:v>` | - | Request header |
| `--input <file>` | - | Raw request body (`-` reads stdin) |
| `--paginate` | `false` | Follow `Link rel="next"` and aggregate JSON arrays (GET only) |
| `-i`, `--include` | `false` | Print status line and headers |

```bash
orbit -p homelab fj api /api/v1/version
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/branches --paginate
orbit -p homelab fj api /api/v1/repos/jorgemuza/orbit-cli/issues/12 -X PATCH -f state=closed
```

Requests are pinned to the configured host - an endpoint, pagination link or redirect pointing elsewhere is refused rather than followed with your token.

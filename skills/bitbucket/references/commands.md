# Bitbucket CLI Command Reference

Complete reference for all `orbit bitbucket` (alias: `bb`) commands and flags.

## Table of Contents

- [project](#project)
- [repo](#repo)
- [branch](#branch)
- [tag](#tag)
- [commit](#commit)
- [pr (pull request)](#pr-pull-request)
- [user](#user)
- [reviewer-condition (admin)](#reviewer-condition-admin)

## Global Flags

These flags apply to all bitbucket commands:

| Flag | Description |
|------|-------------|
| `-p, --profile <name>` | Profile to use (required) |
| `-o, --output <format>` | Output format: `table` (default), `json`, `yaml` |
| `--service <name>` | Bitbucket service name (if profile has multiple) |

---

## project

### `bitbucket project view <project-key>`

View a single project.

```
orbit -p myprofile bb project view L3SUP
```

**Output fields:** Key, Name, Description, Public, Type

### `bitbucket project list`

List all projects.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

---

## repo

### `bitbucket repo view <project-key> <repo-slug>`

View repository details including clone URLs.

```
orbit -p myprofile bb repo view L3SUP agents-sre
```

**Output fields:** ID, Slug, Name, Description, State, SCM, Forkable, Project, Clone URLs

### `bitbucket repo list <project-key>`

List repositories in a project.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

---

## branch

### `bitbucket branch list <project-key> <repo-slug>`

List branches.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter branches by name |
| `--limit <n>` | 50 | Max results |

### `bitbucket branch default <project-key> <repo-slug>`

Show the default branch.

### `bitbucket branch create <project-key> <repo-slug> <name> <start-point>`

Create a branch from a ref (branch name, tag, or commit SHA).

```
orbit -p myprofile bb branch create L3SUP agents-sre feature/new-thing main
```

### `bitbucket branch delete <project-key> <repo-slug> <branch-name>`

Delete a branch.

```
orbit -p myprofile bb branch delete L3SUP agents-sre feature/old-thing
```

---

## tag

### `bitbucket tag list <project-key> <repo-slug>`

List tags.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter tags by name |
| `--limit <n>` | 50 | Max results |

### `bitbucket tag create <project-key> <repo-slug> <tag-name> <start-point>`

Create a tag from a ref.

| Flag | Description |
|------|-------------|
| `-m, --message <text>` | Tag message (annotated tag) |

```
orbit -p myprofile bb tag create L3SUP agents-sre v1.0.0 main -m "Release v1.0.0"
```

---

## commit

### `bitbucket commit list <project-key> <repo-slug>`

List recent commits.

| Flag | Default | Description |
|------|---------|-------------|
| `--branch <name>` | default branch | Branch to list commits from |
| `--limit <n>` | 20 | Max results |

### `bitbucket commit view <project-key> <repo-slug> <commit-id>`

View a commit's details.

```
orbit -p myprofile bb commit view L3SUP agents-sre abc1234def5678
```

**Output fields:** ID, Author, Date, Message

---

## pr (pull request)

Alias: `pull-request`

### `bitbucket pr list <project-key> <repo-slug>`

List pull requests.

| Flag | Default | Description |
|------|---------|-------------|
| `--state <state>` | | Filter: `OPEN`, `MERGED`, `DECLINED`, `ALL` |
| `--limit <n>` | 25 | Max results |

**Output columns:** ID, State, Author, Title

### `bitbucket pr view <project-key> <repo-slug> <pr-id>`

View pull request details.

**Output fields:** ID, Title, State, From Branch, To Branch, Author, Reviewers (with approval status), Created, Updated, URL, Description

### `bitbucket pr create <project-key> <repo-slug>`

Create a pull request.

| Flag | Required | Description |
|------|----------|-------------|
| `--title <text>` | Yes | PR title |
| `--from <branch>` | Yes | Source branch |
| `--to <branch>` | Yes | Target branch |
| `--description <text>` | No | PR description |
| `--reviewers <slugs>` | No | Comma-separated reviewer usernames |

```
orbit -p myprofile bb pr create L3SUP agents-sre \
  --from feature/new --to main --title "Add new feature" \
  --reviewers john.doe,jane.smith
```

### `bitbucket pr diff <project-key> <repo-slug> <pr-id>`

Show the unified diff for a pull request. Useful for code review.

| Flag | Default | Description |
|------|---------|-------------|
| `--context <n>` | 3 | Number of context lines around changes |

```
orbit -p myprofile bb pr diff L3SUP agents-sre 42
orbit -p myprofile bb pr diff L3SUP agents-sre 42 --context 10
```

### `bitbucket pr merge <project-key> <repo-slug> <pr-id>`

Merge a pull request. Automatically handles version for optimistic locking.

| Flag | Default | Description |
|------|---------|-------------|
| `--bypass-review` | `false` | Temporarily disable PRE_PULL_REQUEST_MERGE hooks AND set project-level default reviewer conditions to 0 required approvals before merging, then restore all settings after. Requires project/repo admin permissions. |

```
orbit -p myprofile bb pr merge L3SUP agents-sre 42
orbit -p myprofile bb pr merge L3SUP agents-sre 42 --bypass-review
```

### `bitbucket pr approve <project-key> <repo-slug> <pr-id>`

Approve a pull request as the current user.

### `bitbucket pr unapprove <project-key> <repo-slug> <pr-id>`

Remove your approval from a pull request.

### `bitbucket pr needs-work <project-key> <repo-slug> <pr-id>`

**Alias:** `request-changes`

Flag a pull request as `NEEDS_WORK` - Bitbucket Server's "request changes" signal. Blocks the merge without closing the PR (unlike `decline`). **Server only.** Pair it with `pr comment` so the author has the actionable feedback.

All three status commands verify the status Bitbucket stored, so a success line means the flag is really on the PR.

### `bitbucket pr decline <project-key> <repo-slug> <pr-id>`

Decline a pull request. Automatically handles version for optimistic locking.

### `bitbucket pr comment <project-key> <repo-slug> <pr-id>`

Add a comment to a pull request.

| Flag | Required | Description |
|------|----------|-------------|
| `--body <text>` | Yes | Comment text |

### `bitbucket pr activity <project-key> <repo-slug> <pr-id>`

List pull request activity — comments, approvals, merges, status changes.

| Flag | Default | Description |
|------|---------|-------------|
| `--limit <n>` | 50 | Max results |

**Activity types shown:** COMMENTED, APPROVED, REVIEWED (needs work), MERGED, DECLINED, OPENED

Nested comment replies are shown with indentation — all reply threads are traversed recursively so no responses are hidden.

---

## user

### `bitbucket user list`

List users.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter by username or display name |
| `--limit <n>` | 25 | Max results |

---

## reviewer-condition (admin)

Alias: `rc`

Manage project-level default reviewer conditions. These conditions auto-assign reviewers to PRs and can enforce a minimum number of required approvals before merging.

### `bitbucket reviewer-condition list <project-key>`

List all default reviewer conditions for a project. Shows condition ID, scope, source/target branch matchers, required approvals count, and reviewer names.

```
orbit -p myprofile bb rc list EPCAP
orbit -p myprofile bb rc list EPCAP -o json
```

**Output fields:** ID, Scope (PROJECT/REPOSITORY), Source, Target, Required Approvals, Reviewers

### `bitbucket reviewer-condition update <project-key> <condition-id>`

Update the required approvals count for a default reviewer condition. Preserves all reviewers and branch matchers — only changes the approval requirement.

| Flag | Required | Description |
|------|----------|-------------|
| `--required-approvals <n>` | Yes | Number of required approvals (0 = no requirement) |

```
# Temporarily bypass required approvals
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 0

# Restore to 2 required approvals
orbit -p myprofile bb rc update EPCAP 1063 --required-approvals 2
```

**Common use case — bypass merge block:**
When `pr merge` fails with "Not all required reviewers have approved yet", this is caused by a project-level condition. Set required approvals to 0, merge, then restore. The `--bypass-review` flag on `pr merge` does this automatically.

### `bitbucket reviewer-condition delete <project-key> <condition-id>`

Delete a default reviewer condition permanently.

```
orbit -p myprofile bb rc delete EPCAP 1063
```

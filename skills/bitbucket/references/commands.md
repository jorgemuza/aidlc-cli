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

### `bitbucket pr merge <project-key> <repo-slug> <pr-id>`

Merge a pull request. Automatically handles version for optimistic locking.

### `bitbucket pr decline <project-key> <repo-slug> <pr-id>`

Decline a pull request.

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

---

## user

### `bitbucket user list`

List users.

| Flag | Default | Description |
|------|---------|-------------|
| `--filter <text>` | | Filter by username or display name |
| `--limit <n>` | 25 | Max results |

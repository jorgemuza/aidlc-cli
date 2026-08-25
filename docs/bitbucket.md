# Orbit Bitbucket Command Reference

Manage Bitbucket resources from the command line.

**Command:** `orbit bitbucket`
**Alias:** `orbit bb`

## Table of Contents

- [Global Flags](#global-flags)
- [API Variants](#api-variants)
- [Projects](#project)
- [Repositories](#repo)
- [Branches](#branch)
- [Tags](#tag)
- [Commits](#commit)
- [Pull Requests](#pr)
- [Users](#user)
- [Reviewer Conditions](#reviewer-condition)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--service` | Bitbucket service name. Required only when a profile has multiple Bitbucket services configured. |
| `--debug` | Print HTTP request/response details to stderr for troubleshooting. |

All examples below use `-p myprofile` to specify the Orbit profile.

---

## API Variants

Orbit supports two Bitbucket backends:

| Variant | API | Base URL |
|---------|-----|----------|
| **Cloud** | Bitbucket Cloud REST API v2.0 | `https://api.bitbucket.org/2.0` (default) |
| **Server** | Bitbucket Data Center REST API | Requires `--base-url` |

- **Cloud** uses `workspace/repo` format for addressing repositories.
- **Server** uses `project-key/repo-slug` format.

---

## project

Manage Bitbucket projects.

### `project list`

List projects.

**Alias:** `project ls`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of projects to return. |

```bash
# List all projects
orbit bb project list -p myprofile

# List with a limit
orbit bb project list -p myprofile --limit 10
```

### `project view`

View details of a single project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key (e.g. `TEAM`). |

```bash
orbit bb project view TEAM -p myprofile
```

---

## repo

Manage repositories within a project.

### `repo list`

List repositories in a project.

**Alias:** `repo ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of repositories to return. |

```bash
# List all repos in a project
orbit bb repo list TEAM -p myprofile

# Limit results
orbit bb repo list TEAM -p myprofile --limit 5
```

### `repo view`

View details of a single repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

```bash
orbit bb repo view TEAM my-service -p myprofile
```

### `repo create`

Create a new repository in a project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | *(required)* | Repository name. |
| `--description` | string | | Repository description. |
| `--default-branch` | string | | Default branch name (e.g. `main`). |
| `--public` | bool | `false` | Make the repository publicly visible. |
| `--forkable` | bool | `true` | Allow forking. |

```bash
orbit bb repo create L3SUP --name my-service --description "Payment gateway" -p myprofile
orbit bb repo create L3SUP --name my-api --public -p myprofile
```

### `repo edit`

Update a repository's settings. **Alias:** `repo update`.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | New repository name. |
| `--description` | string | New description. |
| `--default-branch` | string | New default branch. |
| `--public` | bool | Set public visibility. |
| `--forkable` | bool | Allow/disallow forking. |

```bash
orbit bb repo edit L3SUP my-service --description "Updated description" -p myprofile
orbit bb repo edit L3SUP my-service --public=false --forkable=false -p myprofile
```

### `repo delete`

Delete a repository. **Irreversible.**

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

```bash
orbit bb repo delete L3SUP my-service -p myprofile
```

### `repo fork`

Fork a repository, optionally into a different project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | Source project key. |
| `repo-slug` | 2 | Source repository slug. |

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | Fork name (default: same as source). |
| `--target-project` | string | Target project key (default: same project). |

```bash
orbit bb repo fork L3SUP agents-sre --name agents-sre-fork -p myprofile
orbit bb repo fork L3SUP agents-sre --target-project MYPROJ -p myprofile
```

### `repo forks`

List forks of a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 25 | Maximum forks to return. |

```bash
orbit bb repo forks L3SUP agents-sre -p myprofile
```

### `repo permissions`

List user-level permissions on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Maximum results. |

```bash
orbit bb repo permissions L3SUP agents-sre -p myprofile
```

### `repo grant`

Grant a user permission on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `username` | 3 | Bitbucket username (slug). |

| Flag | Type | Description |
|------|------|-------------|
| `--permission` | string | *(required)* One of `REPO_READ`, `REPO_WRITE`, `REPO_ADMIN`. |

```bash
orbit bb repo grant L3SUP agents-sre john.doe --permission REPO_WRITE -p myprofile
```

### `repo revoke`

Revoke a user's permission on a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `username` | 3 | Bitbucket username (slug). |

```bash
orbit bb repo revoke L3SUP agents-sre john.doe -p myprofile
```

---

## branch

Manage branches in a repository.

### `branch list`

List branches.

**Alias:** `branch ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of branches to return. |

```bash
orbit bb branch list TEAM my-service -p myprofile
orbit bb branch ls TEAM my-service -p myprofile --limit 20
```

### `branch create`

Create a new branch.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `branch` | 3 | Name of the new branch. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--ref` | string | Yes | The commit SHA or branch name to branch from. |

```bash
# Create a feature branch from main
orbit bb branch create TEAM my-service feature/login --ref main -p myprofile

# Create a branch from a specific commit
orbit bb branch create TEAM my-service hotfix/urgent --ref a1b2c3d -p myprofile
```

### `branch delete`

Delete a branch.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `branch` | 3 | Name of the branch to delete. |

```bash
orbit bb branch delete TEAM my-service feature/login -p myprofile
```

### `branch default`

View or set the default branch of a repository.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--set` | string | No | Set the default branch to this value. Omit to view the current default. |

```bash
# View the current default branch
orbit bb branch default TEAM my-service -p myprofile

# Change the default branch to develop
orbit bb branch default TEAM my-service -p myprofile --set develop
```

---

## tag

Manage tags in a repository.

### `tag list`

List tags.

**Alias:** `tag ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of tags to return. |

```bash
orbit bb tag list TEAM my-service -p myprofile
orbit bb tag ls TEAM my-service -p myprofile --limit 10
```

### `tag create`

Create a new tag.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `tag` | 3 | Name of the tag. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--ref` | string | Yes | The commit SHA or branch name to tag. |

```bash
# Tag the current HEAD of main
orbit bb tag create TEAM my-service v1.0.0 --ref main -p myprofile

# Tag a specific commit
orbit bb tag create TEAM my-service v1.0.1 --ref a1b2c3d -p myprofile
```

---

## commit

Browse commit history.

### `commit list`

List commits in a repository.

**Alias:** `commit ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of commits to return. |

```bash
orbit bb commit list TEAM my-service -p myprofile
orbit bb commit ls TEAM my-service -p myprofile --limit 50
```

### `commit view`

View details of a single commit.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `sha` | 3 | The commit SHA. |

```bash
orbit bb commit view TEAM my-service a1b2c3d4e5f6 -p myprofile
```

---

## pr

Manage pull requests.

### `pr list`

List pull requests in a repository.

**Alias:** `pr ls`

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | — | Filter by state: `open`, `declined`, `merged`, or `all`. |
| `--limit` | int | — | Maximum number of pull requests to return. |

```bash
# List open PRs (default)
orbit bb pr list TEAM my-service -p myprofile

# List merged PRs
orbit bb pr list TEAM my-service -p myprofile --state merged

# List all PRs with a limit
orbit bb pr ls TEAM my-service -p myprofile --state all --limit 25
```

### `pr view`

View details of a single pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr view TEAM my-service 42 -p myprofile
```

### `pr create`

Create a new pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--source` | string | Yes | Source branch name. |
| `--target` | string | Yes | Target branch name. |
| `--title` | string | Yes | Pull request title. |
| `--description` | string | No | Pull request description body. |

```bash
# Minimal PR
orbit bb pr create TEAM my-service \
  --source feature/login \
  --target main \
  --title "Add login flow" \
  -p myprofile

# PR with description
orbit bb pr create TEAM my-service \
  --source feature/login \
  --target main \
  --title "Add login flow" \
  --description "Implements OAuth2 login with refresh token support." \
  -p myprofile
```

### `pr merge`

Merge a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr merge TEAM my-service 42 -p myprofile
```

### `pr decline`

Decline a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr decline TEAM my-service 42 -p myprofile
```

### `pr comment`

Add a comment to a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--body` | string | Yes | The comment text. |

```bash
orbit bb pr comment TEAM my-service 42 --body "LGTM, approved." -p myprofile
```

### `pr diff`

View the unified diff for a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--context` | int | 3 | Number of context lines around changes. |

```bash
orbit bb pr diff TEAM my-service 42 -p myprofile
orbit bb pr diff TEAM my-service 42 --context 10 -p myprofile
```

### `pr approve`

Approve a pull request as the current user.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr approve TEAM my-service 42 -p myprofile
```

### `pr unapprove`

Remove your approval from a pull request.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr unapprove TEAM my-service 42 -p myprofile
```

### `pr needs-work`

Flag a pull request as `NEEDS_WORK` — Bitbucket Server's "request changes" signal. Blocks merge the same way the manual "Needs Work" button in the UI does, without closing the PR (unlike `decline`). Aliases: `request-changes`. **Server only** — Bitbucket Cloud has a different approval model.

Works even when the caller isn't already listed as a reviewer on the PR, and fails loudly if Bitbucket doesn't actually store the flag - see the note under this section.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
# Leave the review comment first, then flip to needs-work
orbit bb pr comment TEAM my-service 42 -m "Please address the three items above" -p myprofile
orbit bb pr needs-work TEAM my-service 42 -p myprofile

# When the author pushes fixes, approve to clear the status
orbit bb pr approve TEAM my-service 42 -p myprofile
```

> `approve`, `unapprove`, and `needs-work` all share a single `SetReviewStatus` helper that `PUT`s to `/participants/{userSlug}` with the target status, then checks the participant Bitbucket echoes back to confirm the status was stored. The check is not paranoia: `POST /participants` accepts a `status` field, answers `200`, and stores only the role - so the earlier POST-based implementation reported success while leaving the caller `UNAPPROVED` and the PR unflagged.
>
> If the `PUT` is refused (some instances reject it for a user who isn't on the PR yet, though Bitbucket documents it as adding them implicitly), the helper `POST`s to `/participants` with `role=REVIEWER` to enrol the caller and retries the `PUT` once. If that fails too, the error you see is Bitbucket's own refusal of the status change.
>
> The caller's username is resolved from the profile's `auth.username` (basic auth) or from the `X-AUSERNAME` header on an authed GET (PAT/bearer), then memoized for the life of the process.

### `pr activity`

View the activity feed of a pull request (comments, approvals, status changes).

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `repo-slug` | 2 | The repository slug. |
| `pr-id` | 3 | The pull request ID. |

```bash
orbit bb pr activity TEAM my-service 42 -p myprofile
```

Nested comment replies are shown with indentation — all reply threads are traversed recursively.

---

## user

Manage Bitbucket users.

### `user list`

List users.

**Alias:** `user ls`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | — | Maximum number of users to return. |

```bash
orbit bb user list -p myprofile
orbit bb user ls -p myprofile --limit 100
```

---

## reviewer-condition

Manage project-level default reviewer conditions. These auto-assign reviewers to PRs and enforce required approvals.

**Alias:** `rc`

### `reviewer-condition list`

List all default reviewer conditions for a project.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |

```bash
orbit bb reviewer-condition list TEAM -p myprofile
orbit bb rc list TEAM -p myprofile -o json
```

### `reviewer-condition update`

Update the required approvals count for a condition.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `condition-id` | 2 | The condition ID. |

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--required-approvals` | int | Yes | Number of required approvals (0 = no requirement). |

```bash
# Temporarily bypass required approvals
orbit bb rc update TEAM 1063 --required-approvals 0 -p myprofile

# Restore required approvals
orbit bb rc update TEAM 1063 --required-approvals 2 -p myprofile
```

### `reviewer-condition delete`

Delete a default reviewer condition permanently.

| Argument | Position | Description |
|----------|----------|-------------|
| `project-key` | 1 | The project key. |
| `condition-id` | 2 | The condition ID. |

```bash
orbit bb rc delete TEAM 1063 -p myprofile
```

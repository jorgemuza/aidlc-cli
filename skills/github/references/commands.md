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
View a workflow run.
```
orbit github run view octocat/hello-world 12345
```

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

## Global Flags

All commands inherit these flags:
- `-o, --output` — output format: table, json, yaml (default: table)
- `-p, --profile` — profile to use
- `--service` — github service name (if profile has multiple)

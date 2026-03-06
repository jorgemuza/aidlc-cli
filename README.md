# aidlc

A unified CLI for managing connections to development lifecycle services.

Supports **Jira**, **Confluence**, **GitLab**, and **Bitbucket** (cloud and self-hosted). Organize connections into profiles to switch between projects seamlessly.

Secrets can be stored as [1Password](https://1password.com/) references (`op://vault/item/field`) and are resolved at runtime using the 1Password CLI.

## Install

### MacOS - Homebrew

```bash
brew install jorgemuza/tap/aidlc
```

### Windows - Scoop

```bash
  scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket
  scoop install aidlc
```

### From source

```bash
go install github.com/jorgemuza/aidlc-cli@latest
```

### Binary releases

Download from [GitHub Releases](https://github.com/jorgemuza/aidlc-cli/releases).

## Quick Start

```bash
# Create a profile
aidlc profile create --name my-project --default

# Add services
aidlc profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method token --token "op://Dev/jira-token/credential"

aidlc profile add-service \
  --name gitlab-onprem --type gitlab --variant server \
  --base-url https://gitlab.internal.com \
  --auth-method basic --username admin --password "op://Dev/gitlab/password"

# Test connectivity
aidlc service ping
```

## Commands

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `~/.config/aidlc/config.yaml` |
| `-p, --profile` | Profile to use (overrides default) | |
| `-o, --output` | Output format: `table`, `json`, `yaml` | `table` |

### `profile create`

Create a new profile.

```bash
aidlc profile create --name project-a --description "Project A services" --default
```

| Flag | Description |
|------|-------------|
| `-n, --name` | Profile name (required) |
| `-d, --description` | Profile description |
| `--default` | Set as default profile |

### `profile list`

List all profiles with a summary of their services.

```bash
aidlc profile list
aidlc profile list -o json
```

### `profile show`

Show details of a specific profile.

```bash
aidlc profile show my-project
```

### `profile use`

Set a profile as the default.

```bash
aidlc profile use my-project
```

### `profile delete`

Delete a profile.

```bash
aidlc profile delete old-project
```

### `profile add-service`

Add a service connection to a profile.

```bash
# Jira Cloud with API token
aidlc profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method token --token "my-api-token"

# GitLab self-hosted with basic auth and 1Password secret
aidlc profile add-service \
  --name gitlab-onprem --type gitlab --variant server \
  --base-url https://gitlab.internal.com \
  --auth-method basic --username admin \
  --password "op://Dev/gitlab/password"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-n, --name` | Service connection name (required) | |
| `-t, --type` | Service type: `jira`, `confluence`, `gitlab`, `bitbucket` (required) | |
| `--variant` | Variant: `cloud`, `server` | `cloud` |
| `--base-url` | Base URL of the service | |
| `--auth-method` | Auth method: `token`, `basic`, `oauth2` | `token` |
| `--token` | API token or PAT (supports `op://` references) | |
| `--username` | Username for basic auth | |
| `--password` | Password for basic auth (supports `op://` references) | |
| `--client-id` | OAuth2 client ID | |
| `--client-secret` | OAuth2 client secret (supports `op://` references) | |

### `profile remove-service`

Remove a service connection from a profile.

```bash
aidlc profile remove-service --name jira-cloud
```

### `service list`

List all services in the active profile.

```bash
aidlc service list
aidlc -p other-project service list
```

### `service ping`

Test connectivity to services. Pings all services in the active profile, or a specific one by name.

```bash
# Ping all services
aidlc service ping

# Ping a specific service
aidlc service ping jira-cloud
```

### `version`

Print version, commit hash, and build date.

```bash
aidlc version
```

## Supported Services

### Connection Setup

Each service requires a `--base-url`, `--variant` (cloud/server), and authentication. Auth credentials support [1Password references](#1password-integration).

| Parameter | Description | Default |
|-----------|-------------|---------|
| `--base-url` | Service base URL (optional for GitLab/Bitbucket cloud) | |
| `--variant` | `cloud` or `server` | `cloud` |
| `--auth-method` | `token`, `basic`, or `oauth2` | `token` |
| `--token` | API token / PAT | |
| `--username` | Username for basic auth | |
| `--password` | Password for basic auth | |

```bash
# Example: add Jira Cloud
aidlc profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method token --token "op://Dev/jira/token"
```

---

### Jira

Connects to Jira Cloud or Data Center/Server. `--base-url` is required for both variants.

#### `jira ping`

Tests connectivity via `GET /rest/api/2/serverInfo`. Returns version and server title.

```bash
aidlc service ping jira-cloud
```

#### `jira issue list`

List issues with filtering. Supports JQL queries.

```bash
# List issues assigned to me
aidlc jira issue list --assignee me

# Filter by status and type
aidlc jira issue list --status "In Progress" --type Bug

# Raw JQL query
aidlc jira issue list --jql "project = PROJ AND sprint in openSprints()"

# Filter by date
aidlc jira issue list --created-after 2024-01-01 --updated today

# Output formats
aidlc jira issue list -o json
aidlc jira issue list -o yaml
```

| Flag | Description |
|------|-------------|
| `-t, --type` | Filter by issue type (Bug, Story, Task, etc.) |
| `-s, --status` | Filter by status (repeatable) |
| `-y, --priority` | Filter by priority |
| `-a, --assignee` | Filter by assignee (email or display name, `me` for self) |
| `-r, --reporter` | Filter by reporter |
| `-l, --label` | Filter by label (repeatable) |
| `-C, --component` | Filter by component |
| `-P, --parent` | Filter by parent issue key |
| `-q, --jql` | Raw JQL query |
| `--created` | Created date filter (`today`, `week`, `month`, `year`, or date) |
| `--updated` | Updated date filter (same formats as `--created`) |
| `--created-after` | Created after date |
| `--created-before` | Created before date |
| `--updated-after` | Updated after date |
| `--updated-before` | Updated before date |
| `--order-by` | Order by field (default: `created`) |
| `--reverse` | Reverse display order |
| `--paginate` | Pagination `<from>:<limit>` (default: `0:100`) |

#### `jira issue view`

View a single issue with details and comments.

```bash
aidlc jira issue view PROJ-123
aidlc jira issue view PROJ-123 --comments 5
aidlc jira issue view PROJ-123 -o json
```

| Flag | Description |
|------|-------------|
| `--comments` | Number of recent comments to show (default: 1) |

#### `jira issue create`

Create a new issue.

```bash
aidlc jira issue create \
  --type Story \
  --summary "Implement login page" \
  --body "As a user, I want to log in to the app" \
  --priority High \
  --assignee "john@example.com" \
  --label backend --label auth \
  --component "Web App"

# Create a subtask
aidlc jira issue create \
  --type Sub-task \
  --parent PROJ-123 \
  --summary "Add password validation"

# With custom fields
aidlc jira issue create \
  --type Bug \
  --summary "Login fails" \
  --custom "customfield_10001=team-alpha"
```

| Flag | Description |
|------|-------------|
| `-t, --type` | Issue type: Bug, Story, Task, Epic, Sub-task (required) |
| `-s, --summary` | Issue summary/title |
| `-b, --body` | Issue description |
| `-y, --priority` | Priority |
| `-a, --assignee` | Assignee (email or display name) |
| `-r, --reporter` | Reporter |
| `-P, --parent` | Parent issue key (for subtasks) |
| `-l, --label` | Labels (repeatable) |
| `-C, --component` | Components (repeatable) |
| `--fix-version` | Fix versions (repeatable) |
| `--affects-version` | Affects versions (repeatable) |
| `-e, --original-estimate` | Time estimate (e.g. `2d 3h 30m`) |
| `--custom` | Custom field `KEY=VALUE` (repeatable) |

#### `jira issue edit`

Update an existing issue.

```bash
aidlc jira issue edit PROJ-123 \
  --summary "Updated title" \
  --priority Critical

# Add and remove labels (prefix with - to remove)
aidlc jira issue edit PROJ-123 \
  --label new-label --label -old-label
```

| Flag | Description |
|------|-------------|
| `-s, --summary` | New summary |
| `-b, --body` | New description |
| `-y, --priority` | New priority |
| `-l, --label` | Add labels (prefix with `-` to remove) |
| `-C, --component` | Add components (prefix with `-` to remove) |
| `--fix-version` | Add fix versions (prefix with `-` to remove) |
| `--affects-version` | Add affects versions (prefix with `-` to remove) |

#### `jira issue assign`

Assign or unassign an issue.

```bash
aidlc jira issue assign PROJ-123 "john@example.com"
aidlc jira issue assign PROJ-123 me        # assign to self
aidlc jira issue assign PROJ-123 x         # unassign
```

#### `jira issue move`

Transition an issue to a new workflow state.

```bash
aidlc jira issue move PROJ-123 "In Progress"
aidlc jira issue move PROJ-123 Done --comment "Fixed in v2.1"
aidlc jira issue move PROJ-123 Done --resolution Fixed
```

| Flag | Description |
|------|-------------|
| `--comment` | Add comment during transition |
| `-a, --assignee` | Assign during transition |
| `-R, --resolution` | Set resolution (e.g. Fixed, Won't Fix) |

#### `jira issue delete`

Delete an issue.

```bash
aidlc jira issue delete PROJ-123
aidlc jira issue delete PROJ-123 --cascade   # delete with subtasks
```

#### `jira issue comment`

Add a comment to an issue.

```bash
aidlc jira issue comment PROJ-123 "This is fixed now"
aidlc jira issue comment PROJ-123 --body "Internal note" --internal
```

| Flag | Description |
|------|-------------|
| `--internal` | Post as internal comment (Service Desk) |

#### `jira issue link`

Link two issues together.

```bash
aidlc jira issue link PROJ-100 PROJ-200 Blocks
aidlc jira issue link PROJ-100 PROJ-200 Duplicates
```

#### `jira issue unlink`

Remove a link between two issues.

```bash
aidlc jira issue unlink PROJ-100 PROJ-200
```

#### `jira issue worklog`

Log time spent on an issue.

```bash
aidlc jira issue worklog PROJ-123 "2h 30m" --comment "Code review"
```

#### `jira issue clone`

Clone an issue with optional modifications.

```bash
aidlc jira issue clone PROJ-123 --summary "Cloned: new title"
aidlc jira issue clone PROJ-123 --replace "v1:v2"
```

#### `jira epic list`

List epics or issues within an epic.

```bash
aidlc jira epic list                    # list all epics
aidlc jira epic list PROJ-50            # issues in epic PROJ-50
```

#### `jira epic create`

Create a new epic.

```bash
aidlc jira epic create \
  --name "Q1 Auth Overhaul" \
  --summary "Revamp authentication system" \
  --body "Replace legacy auth with OAuth2"
```

#### `jira epic add / remove`

Add or remove issues from an epic (up to 50 at once).

```bash
aidlc jira epic add PROJ-50 PROJ-101 PROJ-102 PROJ-103
aidlc jira epic remove PROJ-101 PROJ-102
```

#### `jira sprint list`

List sprints or issues in a sprint.

```bash
aidlc jira sprint list                  # list all sprints
aidlc jira sprint list 42               # issues in sprint 42
aidlc jira sprint list --current        # active sprint issues
aidlc jira sprint list --state active   # filter by sprint state
```

| Flag | Description |
|------|-------------|
| `--current` | Active sprint |
| `--prev` | Previous sprint |
| `--next` | Next planned sprint |
| `--state` | Filter by state: `future`, `active`, `closed` |

#### `jira sprint add`

Add issues to a sprint (up to 50 at once).

```bash
aidlc jira sprint add 42 PROJ-101 PROJ-102
```

#### `jira board list`

List all boards in the project.

```bash
aidlc jira board list
```

#### `jira project list`

List all accessible projects.

```bash
aidlc jira project list
```

#### `jira release list`

List versions/releases for the project.

```bash
aidlc jira release list
```

---

### Confluence

Connects to Confluence Cloud or Data Center/Server. `--base-url` is required for both variants.

- **Cloud** API path: `/wiki/rest/api`
- **Server** API path: `/rest/api`

#### `confluence ping`

Tests connectivity by querying spaces. Returns variant and space count.

```bash
aidlc service ping confluence-cloud
```

#### `confluence page read`

Read a page's content in different formats.

```bash
aidlc confluence page read 12345
aidlc confluence page read 12345 --format markdown
aidlc confluence page read 12345 --format html
aidlc confluence page read "https://myco.atlassian.net/wiki/spaces/DEV/pages/12345"
```

| Flag | Description |
|------|-------------|
| `--format` | Output format: `storage` (XML), `html`, `markdown`, `text` |

#### `confluence page info`

Get page metadata (title, space, version, author, dates).

```bash
aidlc confluence page info 12345
```

#### `confluence page create`

Create a new page in a space.

```bash
aidlc confluence page create \
  --title "Architecture Guide" \
  --space DEV \
  --body "<h1>Overview</h1><p>System architecture...</p>"

# Create as child page
aidlc confluence page create-child \
  --title "API Reference" \
  --parent 12345 \
  --body-file ./api-docs.html
```

| Flag | Description |
|------|-------------|
| `--title` | Page title (required) |
| `--space` | Space key (required for top-level pages) |
| `--parent` | Parent page ID (for child pages) |
| `--body` | Page content in storage format |
| `--body-file` | Read content from file |

#### `confluence page update`

Update an existing page's title or content.

```bash
aidlc confluence page update 12345 \
  --title "Updated Title" \
  --body "<p>New content</p>"

aidlc confluence page update 12345 --body-file ./updated-content.html
```

#### `confluence page delete`

Delete a page (moves to trash).

```bash
aidlc confluence page delete 12345
```

#### `confluence page move`

Move a page to a new parent within the same space.

```bash
aidlc confluence page move 12345 67890
aidlc confluence page move 12345 67890 --title "Renamed Page"
```

#### `confluence page children`

List child pages, optionally recursive.

```bash
aidlc confluence page children 12345
aidlc confluence page children 12345 --recursive --max-depth 3
aidlc confluence page children 12345 --format tree
```

| Flag | Description |
|------|-------------|
| `--recursive` | Include nested children |
| `--max-depth` | Max depth for recursive listing |
| `--format` | Output format: `list`, `tree`, `json` |

#### `confluence page copy-tree`

Duplicate an entire page tree.

```bash
aidlc confluence page copy-tree 12345 67890
aidlc confluence page copy-tree 12345 67890 \
  --max-depth 2 --exclude "Archive*" --dry-run
```

| Flag | Description |
|------|-------------|
| `--max-depth` | Maximum depth to copy |
| `--exclude` | Glob pattern to exclude pages |
| `--copy-suffix` | Suffix for copied pages (default: ` (copy)`) |
| `--dry-run` | Preview without making changes |
| `--delay` | Delay between API calls (ms) |

#### `confluence page export`

Export a page with optional attachments.

```bash
aidlc confluence page export 12345 --dest ./export/
aidlc confluence page export 12345 --format markdown --recursive
aidlc confluence page export 12345 --skip-attachments
```

| Flag | Description |
|------|-------------|
| `--format` | Export format: `storage`, `html`, `markdown` |
| `--dest` | Destination directory |
| `--recursive` | Export child pages too |
| `--max-depth` | Max depth for recursive export |
| `--skip-attachments` | Don't download attachments |
| `--dry-run` | Preview without downloading |
| `--overwrite` | Overwrite existing files |

#### `confluence search`

Search for content using text or CQL.

```bash
aidlc confluence search "deployment guide"
aidlc confluence search --cql "space = DEV AND type = page AND title ~ 'API*'"
aidlc confluence search "auth" --limit 20
```

| Flag | Description |
|------|-------------|
| `--cql` | Confluence Query Language expression |
| `--limit` | Max results (default: 25) |

#### `confluence space list`

List all accessible spaces.

```bash
aidlc confluence space list
```

#### `confluence page find`

Find a page by title within a space.

```bash
aidlc confluence page find "Getting Started" --space DEV
```

#### `confluence attachment list`

List attachments on a page.

```bash
aidlc confluence attachment list 12345
aidlc confluence attachment list 12345 --pattern "*.pdf" --download --dest ./files/
```

| Flag | Description |
|------|-------------|
| `--pattern` | Glob filter for filenames |
| `--download` | Download matching attachments |
| `--dest` | Download destination directory |

#### `confluence attachment upload`

Upload files to a page.

```bash
aidlc confluence attachment upload 12345 ./diagram.png
aidlc confluence attachment upload 12345 ./report.pdf --replace
```

| Flag | Description |
|------|-------------|
| `--replace` | Replace existing attachment with same name |

#### `confluence attachment delete`

Delete an attachment from a page.

```bash
aidlc confluence attachment delete 12345 att67890
```

#### `confluence comment list`

List comments on a page.

```bash
aidlc confluence comment list 12345
aidlc confluence comment list 12345 --location inline --depth all
```

| Flag | Description |
|------|-------------|
| `--location` | Filter: `inline`, `footer`, `resolved` |
| `--depth` | Comment depth: `top`, `all` |

#### `confluence comment add`

Add a comment to a page.

```bash
aidlc confluence comment add 12345 "Looks good, approved."
aidlc confluence comment add 12345 --body "Please review section 3" --inline
```

#### `confluence comment delete`

Delete a comment.

```bash
aidlc confluence comment delete 98765
```

#### `confluence property list / get / set / delete`

Manage page properties (key-value metadata).

```bash
aidlc confluence property list 12345
aidlc confluence property get 12345 status
aidlc confluence property set 12345 status '{"state":"approved"}'
aidlc confluence property delete 12345 status
```

---

### GitLab

Connects to GitLab.com or self-hosted GitLab. `--base-url` defaults to `https://gitlab.com` for cloud variant.

#### `gitlab ping`

Tests connectivity via `GET /api/v4/version`. Returns version and revision.

```bash
aidlc service ping gitlab-cloud
```

#### `gitlab project list`

List accessible projects.

```bash
aidlc gitlab project list
aidlc gitlab project list --owned
aidlc gitlab project list --search "api"
```

#### `gitlab issue list`

List issues in a project.

```bash
aidlc gitlab issue list --project my-group/my-project
aidlc gitlab issue list --state opened --assignee me
aidlc gitlab issue list --label bug --label urgent
```

| Flag | Description |
|------|-------------|
| `--project` | Project path (required) |
| `--state` | Filter: `opened`, `closed`, `all` |
| `--assignee` | Filter by assignee |
| `--label` | Filter by label (repeatable) |
| `--milestone` | Filter by milestone |
| `--search` | Search in title and description |

#### `gitlab issue view`

View a single issue.

```bash
aidlc gitlab issue view --project my-group/my-project 42
```

#### `gitlab issue create`

Create a new issue.

```bash
aidlc gitlab issue create --project my-group/my-project \
  --title "Fix login timeout" \
  --description "Users report 30s timeout on login" \
  --label bug --assignee john
```

#### `gitlab issue edit`

Update an issue.

```bash
aidlc gitlab issue edit --project my-group/my-project 42 \
  --title "Updated title" --state-event close
```

#### `gitlab mr list`

List merge requests.

```bash
aidlc gitlab mr list --project my-group/my-project
aidlc gitlab mr list --project my-group/my-project --state merged
```

#### `gitlab mr view`

View a merge request.

```bash
aidlc gitlab mr view --project my-group/my-project 15
```

#### `gitlab mr create`

Create a merge request.

```bash
aidlc gitlab mr create --project my-group/my-project \
  --source feature-branch --target main \
  --title "Add login feature" \
  --description "Implements OAuth2 login flow"
```

#### `gitlab pipeline list`

List CI/CD pipelines.

```bash
aidlc gitlab pipeline list --project my-group/my-project
aidlc gitlab pipeline list --project my-group/my-project --status failed
```

#### `gitlab pipeline view`

View pipeline details and jobs.

```bash
aidlc gitlab pipeline view --project my-group/my-project 12345
```

---

### Bitbucket

Connects to Bitbucket Cloud or Data Center/Server. `--base-url` defaults to `https://api.bitbucket.org/2.0` for cloud.

- **Cloud** uses Bitbucket Cloud REST API v2.0
- **Server** uses Bitbucket Data Center REST API

#### `bitbucket ping`

Tests connectivity. Cloud pings `/user`, Server pings `/rest/api/latest/application-properties`.

```bash
aidlc service ping bitbucket-cloud
```

#### `bitbucket repo list`

List repositories.

```bash
aidlc bitbucket repo list --workspace myteam
```

#### `bitbucket pr list`

List pull requests.

```bash
aidlc bitbucket pr list --repo myteam/myrepo
aidlc bitbucket pr list --repo myteam/myrepo --state OPEN
```

#### `bitbucket pr view`

View a pull request.

```bash
aidlc bitbucket pr view --repo myteam/myrepo 42
```

#### `bitbucket pr create`

Create a pull request.

```bash
aidlc bitbucket pr create --repo myteam/myrepo \
  --source feature-branch --target main \
  --title "Add search feature" \
  --description "Implements full-text search"
```

#### `bitbucket pipeline list`

List pipelines.

```bash
aidlc bitbucket pipeline list --repo myteam/myrepo
```

## 1Password Integration

Instead of storing secrets in plain text, use 1Password references:

```bash
aidlc profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method token \
  --token "op://DevVault/jira-token/credential"
```

Secrets are resolved at runtime via `op read`, so the [1Password CLI](https://developer.1password.com/docs/cli/) must be installed and authenticated.

## Configuration

Config is stored in YAML at `~/.config/aidlc/config.yaml`:

```yaml
profiles:
  - name: my-project
    description: "My project services"
    default: true
    services:
      - name: jira-cloud
        type: jira
        variant: cloud
        base_url: https://myco.atlassian.net
        auth:
          method: token
          token: "op://DevVault/jira-token/credential"
      - name: gitlab-onprem
        type: gitlab
        variant: server
        base_url: https://gitlab.internal.com
        auth:
          method: basic
          username: admin
          password: "op://DevVault/gitlab/password"
```

## License

MIT

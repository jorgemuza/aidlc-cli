# orbit Jira Commands Reference

Complete reference for all `orbit jira` commands with flags and examples.

> **The command tree is flat.** Every command below is a leaf - none takes a nested sub-action, so `<command> list` is never valid. Reading has its own command name at the same level as writing: `issue view --comments` (not `issue comment list`), `issue attachments` (not `issue attach list`), `field context-list`, `screen tab-list`.
>
> This is worth care because write commands take the issue key as their **first positional argument**. `orbit jira issue comment list PYMT-123` does not fail with "unknown command" - it is parsed as issue key `list`, body `PYMT-123`, and sent to Jira as a real write against an issue that does not exist. If a command in this reference is not listed below, it does not exist; check here rather than guessing a subcommand.

## Table of Contents

- [Global Flags](#global-flags)
- [Issue Commands](#issue-commands)
  - [issue list](#issue-list)
  - [issue view](#issue-view)
  - [issue create](#issue-create)
  - [issue edit](#issue-edit)
  - [issue assign](#issue-assign)
  - [issue move](#issue-move)
  - [issue delete](#issue-delete)
  - [issue comment](#issue-comment)
  - [issue attach](#issue-attach)
  - [issue attachments](#issue-attachments)
  - [issue link](#issue-link)
  - [issue unlink](#issue-unlink)
  - [issue worklog](#issue-worklog)
  - [issue clone](#issue-clone)
- [Export Command](#export-command)
  - [export](#export)
- [Epic Commands](#epic-commands)
- [Sprint Commands](#sprint-commands)
- [Board Commands](#board-commands)
- [Project Commands](#project-commands)
- [Release Commands](#release-commands)
- [Field Commands](#field-commands)
  - [field list](#field-list)
  - [field create](#field-create)
  - [field context-list](#field-context-list)
  - [field option-list](#field-option-list)
  - [field option-add](#field-option-add)
- [Screen Commands](#screen-commands)
  - [screen list](#screen-list)
  - [screen tab-list](#screen-tab-list)
  - [screen tab-create](#screen-tab-create)
  - [screen field-list](#screen-field-list)
  - [screen field-add](#screen-field-add)
  - [screen field-remove](#screen-field-remove)
  - [screen field-move](#screen-field-move)
- [Filter Commands](#filter-commands)
  - [filter list](#filter-list)
  - [filter view](#filter-view)
  - [filter search](#filter-search)
  - [filter create](#filter-create)
  - [filter update](#filter-update)
  - [filter delete](#filter-delete)
- [Dashboard Commands](#dashboard-commands)
  - [dashboard list](#dashboard-list)
  - [dashboard view](#dashboard-view)
  - [dashboard create](#dashboard-create)
  - [dashboard delete](#dashboard-delete)
  - [dashboard gadget list](#dashboard-gadget-list)
  - [dashboard gadget add](#dashboard-gadget-add)
  - [dashboard gadget remove](#dashboard-gadget-remove)
  - [dashboard gadget property list](#dashboard-gadget-property-list)
  - [dashboard gadget property get](#dashboard-gadget-property-get)
  - [dashboard gadget property set](#dashboard-gadget-property-set)
- [Status Commands](#status-commands)
  - [status list](#status-list)
- [Issue Type Commands](#issue-type-commands)
  - [issuetype-list](#issuetype-list)

---

## Global Flags

These flags are available on all jira subcommands:

| Flag | Description |
|------|-------------|
| `-p, --profile` | Profile to use (overrides default) |
| `-o, --output` | Output format: table, json, yaml (default: table) |
| `--service` | Jira service name (if profile has multiple jira services) |
| `--config` | Config file path (default: ~/.config/orbit/config.yaml) |

---

## Issue Commands

### issue list

List and search issues with filtering.

**Aliases:** `ls`, `search`

```bash
orbit -p profile jira issue list [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `-q, --jql` | string | Raw JQL query (overrides other filters) |
| `--project` | string | Filter by project key |
| `-t, --type` | string | Filter by issue type |
| `-s, --status` | []string | Filter by status (repeatable) |
| `-y, --priority` | string | Filter by priority |
| `-a, --assignee` | string | Filter by assignee |
| `-r, --reporter` | string | Filter by reporter |
| `-l, --label` | []string | Filter by label (repeatable) |
| `-C, --component` | string | Filter by component |
| `-P, --parent` | string | Filter by parent issue key |
| `--created` | string | Created date filter |
| `--updated` | string | Updated date filter |
| `--created-after` | string | Created after date |
| `--created-before` | string | Created before date |
| `--updated-after` | string | Updated after date |
| `--updated-before` | string | Updated before date |
| `--order-by` | string | Order by field (default: created) |
| `--reverse` | bool | Reverse display order |
| `--start-at` | int | Pagination start index (default: 0) |
| `--max-results` | int | Max results (default: 50) |
| `-F, --fresh` | bool | Re-verify each row against `/issue/{key}` (bypasses JQL search index lag; ~N extra API calls) |

> **When to use `--fresh`:** right after transitioning/creating/updating issues in the same session, or whenever the user says "current", "fresh", "live", "still open", or "right now". Without it, `issue list` reads the eventually-consistent search index and may lag 30–60s. With it, each row is re-fetched in parallel (concurrency 4, tuned for Jira Cloud's ~10 req/s limit) from the authoritative single-issue endpoint, and transient failures are retried automatically by orbit's shared HTTP layer. Per-key failures keep the stale row and are logged on stderr — the list still renders.

**Examples:**

```bash
orbit -p paybook jira issue list --project PYMT
orbit -p paybook jira issue list --assignee me --status "In Progress"
orbit -p paybook jira issue list --jql "project = PYMT AND issuetype = Epic"
orbit -p paybook jira issue list --project PYMT --type Epic -o json

# Fresh list after just transitioning tickets
orbit -p paybook jira issue list --project PYMT -s Open --fresh
```

---

### issue view

View detailed issue information including summary, status, assignee, description, subtasks, links, and **comments with timestamps**. This is the command for reading comments; `issue comment` only writes them.

**Aliases:** `show`

```bash
orbit -p profile jira issue view <issue-key> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--comments` | int | Number of recent comments to show (default: 10, use 0 to hide) |

The view shows: summary, type, status, priority, assignee, reporter, labels, components, description, subtasks, links, **attachments** (with filename, size, author, date, and download URL), and **comments** (with author, date, and edit indicator). Dates are formatted as `YYYY-MM-DD HH:MM` for readability. Attachments appear before comments for context when screenshots are referenced in comment text.

**Examples:**

```bash
# View issue with last 10 comments (default)
orbit -p paybook jira issue view PYMT-123

# View with all comments
orbit -p paybook jira issue view PYMT-123 --comments 50

# Hide comments
orbit -p paybook jira issue view PYMT-123 --comments 0

# JSON output (includes full comment history)
orbit -p paybook jira issue view PYMT-123 -o json
```

---

### issue create

Create a new issue.

```bash
orbit -p profile jira issue create [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--project` | string | Project key (**required**) |
| `-t, --type` | string | Issue type (**required**) |
| `-s, --summary` | string | Issue summary (**required**) |
| `-b, --body` | string | Issue description (wiki markup for Server, plain text for Cloud) |
| `-y, --priority` | string | Priority |
| `-a, --assignee` | string | Assignee username (Server) or accountId (Cloud) |
| `-r, --reporter` | string | Reporter username |
| `-P, --parent` | string | Parent issue key (for sub-tasks only) |
| `-l, --label` | []string | Labels (repeatable) |
| `-C, --component` | []string | Components (repeatable) |
| `--fix-version` | []string | Fix versions (repeatable) |
| `--affects-version` | []string | Affects versions (repeatable) |
| `-e, --original-estimate` | string | Time estimate (e.g., 2d 3h) |
| `--epic-name` | string | Epic name (auto-set from summary for Epic type) |
| `-F, --field` | []string | Custom fields as key=value (repeatable) |

**Examples:**

```bash
# Simple story
orbit -p paybook jira issue create --project PYMT --type Story \
  --summary "Add login page"

# Epic with Parent Link
orbit -p paybook jira issue create --type Epic --project PYMT \
  --summary "Okta Authentication Foundation" \
  --priority Highest \
  --field "customfield_27521=PYMT-100"

# Bug with labels and components
orbit -p paybook jira issue create --project PYMT --type Bug \
  --summary "Fix timeout" --priority High \
  --label backend --label urgent --component api
```

**Notes:**
- `--parent` only works for Sub-task types. For linking Epics to Initiatives, use `--field "customfield_27521=KEY"` (Parent Link).
- When `--type Epic` is used, the Epic Name custom field is auto-set from `--summary`. Override with `--epic-name`.
- For Cloud, plain-text description is automatically wrapped in ADF format.

---

### issue edit

Edit an existing issue.

**Aliases:** `update`

```bash
orbit -p profile jira issue edit <issue-key> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `-s, --summary` | string | New summary |
| `-b, --body` | string | New description |
| `-y, --priority` | string | New priority |
| `-l, --label` | []string | Add/remove labels (prefix with `-` to remove) |
| `-C, --component` | []string | Add/remove components (prefix with `-` to remove) |
| `--fix-version` | []string | Add/remove fix versions (prefix with `-` to remove) |
| `-F, --field` | []string | Custom fields as key=value (repeatable) |

**Examples:**

```bash
orbit -p paybook jira issue edit PYMT-123 --summary "Updated title"
orbit -p paybook jira issue edit PYMT-123 --priority Critical
orbit -p paybook jira issue edit PYMT-123 --label new-label --label -old-label
orbit -p paybook jira issue edit PYMT-123 -F customfield_10397=value:Yes
orbit -p paybook jira issue edit PYMT-123 -F customfield_10403=value:Dev

# Multi-select fields require JSON array of objects (value: shorthand returns 400)
orbit -p paybook jira issue edit PYMT-123 -F 'customfield_10398=[{"value":"Developer"},{"value":"QA"}]'
orbit -p paybook jira issue edit PYMT-123 -F 'customfield_10400=[{"value":"Confirmed"}]'
```

---

### issue assign

Assign or unassign an issue.

```bash
orbit -p profile jira issue assign <issue-key> <assignee>
```

Use `x` as assignee to unassign.

**Examples:**

```bash
orbit -p paybook jira issue assign PYMT-123 john.doe
orbit -p paybook jira issue assign PYMT-123 x    # unassign
```

---

### issue move

Transition an issue to a new workflow state.

**Aliases:** `transition`

```bash
orbit -p profile jira issue move <issue-key> <state> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--comment` | string | Add comment during transition |
| `-R, --resolution` | string | Set resolution (e.g., Fixed) |

**Examples:**

```bash
orbit -p paybook jira issue move PYMT-123 "In Progress"
orbit -p paybook jira issue move PYMT-123 "AI Review"
orbit -p paybook jira issue move PYMT-123 "In Code Review"
orbit -p paybook jira issue move PYMT-123 "In QA"
orbit -p paybook jira issue move PYMT-123 Done --comment "Fixed in v2.1" --resolution Fixed
```

---

### issue delete

Delete an issue.

**Aliases:** `rm`, `remove`

```bash
orbit -p profile jira issue delete <issue-key> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--cascade` | bool | Delete subtasks too |

---

### issue comment

Add a comment to an issue. Body can be passed as `--body` flag or as a positional argument. Supports @mentions using Jira mention syntax.

> **Write-only.** To read comments use [`issue view`](#issue-view), which prints them with author and date. There is no `comment list` subcommand. The first positional argument is always the issue key, so `orbit -p profile jira issue comment list PYMT-123` is parsed as issue key `list` with body `PYMT-123` and attempts a real write against a nonexistent issue.

```bash
orbit -p profile jira issue comment <issue-key> --body <body>
orbit -p profile jira issue comment <issue-key> "comment text"
```

**Mentioning users in comments:**

Use `[~accountId:ID]` to @mention a user in Cloud (produces a clickable mention). Use `[~username]` for Server. The account ID can be found via `jira user search` or from issue assignee details (`-o json`).

```bash
# Simple comment
orbit -p paybook jira issue comment PYMT-123 -b "This is fixed now"

# Comment with @mention (Cloud — uses accountId)
orbit -p paybook jira issue comment PYMT-123 -b "Hey [~accountId:5b10ac8d82e05b22cc7d4ef5], please review"

# Comment with @mention (Server — uses username)
orbit -p paybook jira issue comment PYMT-123 -b "Assigned back to [~jorge.padilla] for rework"

# Find a user's accountId for mentions
orbit -p paybook jira user search "Jorge"
```

---

### issue attach

Attach files to an issue. Aliases: `attachment`, `upload`.

```bash
orbit -p profile jira issue attach <issue-key> <file...>
```

```bash
# Attach a single file
orbit -p paybook jira issue attach PYMT-123 report.pdf

# Attach multiple files
orbit -p paybook jira issue attach PYMT-123 screenshot.png log.txt
```

### issue attachments

List attachments on an issue.

```bash
orbit -p paybook jira issue attachments PYMT-123
orbit -p paybook jira issue attachments PYMT-123 -o json
```

### issue download

Download attachments from an issue. Alias: `dl`.

| Flag | Default | Description |
|------|---------|-------------|
| `--output`, `-O` | `.` | Output directory |
| `--images-only` | `false` | Download only image attachments (png, jpg, gif, webp, svg) |

```bash
# Download all attachments to current directory
orbit -p paybook jira issue download PYMT-123

# Download to a specific directory
orbit -p paybook jira issue download PYMT-123 --output ./attachments

# Download only images (screenshots)
orbit -p paybook jira issue download PYMT-123 --images-only --output ./screenshots
```

### issue detach

Remove attachments by ID. Aliases: `rm-attachment`, `remove-attachment`.

```bash
orbit -p paybook jira issue detach 12345
orbit -p paybook jira issue detach 12345 12346 12347
```

Use `issue attachments` to find attachment IDs first.

---

### issue link

Link two issues together.

```bash
orbit -p profile jira issue link <inward-key> <outward-key> <link-type>
```

**Examples:**

```bash
orbit -p paybook jira issue link PYMT-100 PYMT-200 Blocks
orbit -p paybook jira issue link PYMT-100 PYMT-200 "is caused by"
```

---

### issue unlink

Remove a link between two issues.

```bash
orbit -p profile jira issue unlink <inward-key> <outward-key>
```

---

### issue worklog

Log time spent on an issue.

```bash
orbit -p profile jira issue worklog <issue-key> <time-spent> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--comment` | string | Comment for worklog entry |

**Examples:**

```bash
orbit -p paybook jira issue worklog PYMT-123 "2h 30m"
orbit -p paybook jira issue worklog PYMT-123 "1d" --comment "Code review"
```

---

### issue clone

Clone an issue with optional modifications.

```bash
orbit -p profile jira issue clone <issue-key> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `-s, --summary` | string | Override summary |
| `-H, --replace` | []string | Replace text find:replace |

**Examples:**

```bash
orbit -p paybook jira issue clone PYMT-123
orbit -p paybook jira issue clone PYMT-123 --summary "Cloned: new title"
orbit -p paybook jira issue clone PYMT-123 --replace "v1:v2"
```

---

## Export Command

### export

Export one or more Jira epics and their full hierarchy (stories, sub-tasks) into structured markdown directories.

```bash
orbit -p profile jira export <epic-key...> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `-O, --output` | string | Output directory (default: `.`) |
| `--max-results` | int | Max child issues per level (default: 100) |

**Output structure:**

```
<outdir>/<EPIC-KEY>/INDEX.md                    — epic metadata + story list
<outdir>/<EPIC-KEY>/<STORY-KEY>/INDEX.md        — story metadata + sub-task list
<outdir>/<EPIC-KEY>/<STORY-KEY>/<TASK-KEY>.md   — individual sub-task file
```

**Examples:**

```bash
orbit -p paybook jira export BS-1243 -O ./jira-export
orbit -p paybook jira export BS-1243 BS-1255 BS-1246 -O ./jira-export
orbit -p paybook jira export PROJ-100
```

---

## Epic Commands

### epic list

List epics or issues within an epic.

**Aliases:** `ls`

```bash
orbit -p profile jira epic list [epic-key] [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--project` | string | Filter by project key |
| `-s, --status` | []string | Filter by status |
| `--max-results` | int | Max results (default: 50) |

**Examples:**

```bash
orbit -p paybook jira epic list --project PYMT
orbit -p paybook jira epic list PYMT-50    # issues in epic PYMT-50
```

### epic create

Create a new epic.

```bash
orbit -p profile jira epic create [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--project` | string | Project key (**required**) |
| `-n, --name` | string | Epic name (**required**) |
| `-s, --summary` | string | Epic summary (defaults to name) |
| `-b, --body` | string | Epic description |
| `-y, --priority` | string | Priority |
| `-l, --label` | []string | Labels |
| `-C, --component` | []string | Components |

### epic add

Add issues to an epic (max 50).

```bash
orbit -p profile jira epic add <epic-key> <issue-keys...>
```

### epic remove

Remove issues from their epic (max 50).

```bash
orbit -p profile jira epic remove <issue-keys...>
```

---

## Sprint Commands

### sprint list

List sprints or issues in a sprint.

```bash
orbit -p profile jira sprint list [sprint-id] [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--board-id` | int | Board ID (required for listing sprints) |
| `--state` | string | Filter: future, active, closed |
| `--max-results` | int | Max results (default: 50) |

### sprint add

Add issues to a sprint (max 50).

```bash
orbit -p profile jira sprint add <sprint-id> <issue-keys...>
```

### sprint close

Close an active sprint.

```bash
orbit -p profile jira sprint close <sprint-id>
```

### sprint view

View sprint details.

```bash
orbit -p profile jira sprint view <sprint-id>
```

### sprint update

Update sprint name, dates, or goal.

| Flag | Description |
|------|-------------|
| `--name` | Sprint name |
| `--start-date` | Start date (YYYY-MM-DD) |
| `--end-date` | End date (YYYY-MM-DD) |
| `--goal` | Sprint goal |

```bash
orbit -p profile jira sprint update 1366 --name "Sprint 2" --goal "Ship v2.0"
orbit -p profile jira sprint update 1366 --end-date 2026-04-15
```

### sprint remove

Move issues from sprint to backlog. Alias: `backlog`.

```bash
orbit -p profile jira sprint remove PROJ-101 PROJ-102
```

### sprint delete

Delete a sprint.

```bash
orbit -p profile jira sprint delete <sprint-id>
```

---

## Board Commands

### board list

List boards.

```bash
orbit -p profile jira board list [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--project` | string | Filter by project key |

---

## Project Commands

### project list

List all accessible projects.

```bash
orbit -p profile jira project list
```

---

## Release Commands

### release list

List project versions/releases.

```bash
orbit -p profile jira release list --project <key>
```

---

## Field Commands

Manage Jira fields — list system/custom fields, create new custom fields, manage field contexts and options. Field creation and option management are **Cloud only**.

### field list

List Jira fields (system and custom).

```bash
orbit -p profile jira field list [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--filter` | string | Filter fields by name or ID (case-insensitive) |
| `--custom` | bool | Show only custom fields |

**Examples:**

```bash
orbit -p paybook jira field list --filter "AI"
orbit -p paybook jira field list --custom
orbit -p paybook jira field list --filter "customfield_10397"
```

---

### field create

Create a custom field (Cloud only).

```bash
orbit -p profile jira field create [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | Field name (**required**) |
| `--type` | string | Field type or shorthand (**required**) |
| `--description` | string | Field description |
| `--searcher` | string | Searcher key (auto-resolved if omitted) |

**Shorthand types:**

| Shorthand | Full Type |
|-----------|-----------|
| `select` | `com.atlassian.jira.plugin.system.customfieldtypes:select` |
| `multiselect` | `com.atlassian.jira.plugin.system.customfieldtypes:multiselect` |
| `number` / `float` | `com.atlassian.jira.plugin.system.customfieldtypes:float` |
| `checkbox` / `checkboxes` | `com.atlassian.jira.plugin.system.customfieldtypes:multicheckboxes` |
| `text` | `com.atlassian.jira.plugin.system.customfieldtypes:textfield` |
| `textarea` | `com.atlassian.jira.plugin.system.customfieldtypes:textarea` |

**Examples:**

```bash
orbit -p paybook jira field create --name "AI Assisted" --type select \
  --description "Was AI used on this ticket?"

orbit -p paybook jira field create --name "AI Prompt Iterations" --type number \
  --description "How many prompt cycles to get working output"

orbit -p paybook jira field create --name "Human Review Confirmed" --type checkbox \
  --description "Engineer confirms AI output was reviewed"
```

---

### field context-list

List contexts for a custom field (Cloud only).

```bash
orbit -p profile jira field context-list <field-id>
```

**Examples:**

```bash
orbit -p paybook jira field context-list customfield_10397
# Output: 10817  Default Configuration Scheme for AI Assisted [global] [any-issue-type]
```

---

### field option-list

List options for a select/multi-select field context (Cloud only).

```bash
orbit -p profile jira field option-list <field-id> <context-id>
```

**Examples:**

```bash
orbit -p paybook jira field option-list customfield_10397 10817
# Output: 10424  Yes
#         10425  No
#         10426  Partial
```

---

### field option-add

Add options to a select/multi-select field context (Cloud only).

```bash
orbit -p profile jira field option-add <field-id> <context-id> --values <values>
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--values` | []string | Option values to add (comma-separated) |

**Examples:**

```bash
orbit -p paybook jira field option-add customfield_10397 10817 --values "Yes,No,Partial"
orbit -p paybook jira field option-add customfield_10398 10818 \
  --values "Architect,Developer,QA,Reviewer,Simplifier,None"
```

---

## Screen Commands

Manage Jira screens — list screens, create tabs, add/remove/move fields on screen tabs. Screens control which fields appear on issue create, edit, and view forms.

### screen list

List all screens, optionally filtered by name.

```bash
orbit -p profile jira screen list [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--filter` | string | Filter screens by name (case-insensitive) |

**Examples:**

```bash
orbit -p paybook jira screen list
orbit -p paybook jira screen list --filter "PYMT"
orbit -p paybook jira screen list --filter "Scrum Default"
```

---

### screen tab-list

List tabs on a screen.

```bash
orbit -p profile jira screen tab-list <screen-id>
```

**Examples:**

```bash
orbit -p paybook jira screen tab-list 10089
# Output: 10189  General
#         10868  AI Workflow
```

---

### screen tab-create

Create a new tab on a screen.

```bash
orbit -p profile jira screen tab-create <screen-id> <tab-name>
```

**Examples:**

```bash
orbit -p paybook jira screen tab-create 10089 "AI Workflow"
# Output: Created tab 10868 (AI Workflow) on screen 10089
```

---

### screen field-list

List fields on a screen tab.

```bash
orbit -p profile jira screen field-list <screen-id> <tab-id>
```

**Examples:**

```bash
orbit -p paybook jira screen field-list 10089 10868
# Output: customfield_10397  AI Assisted
#         customfield_10398  Agent Role Used
#         ...
```

---

### screen field-add

Add fields to a screen tab.

```bash
orbit -p profile jira screen field-add <screen-id> <tab-id> --fields <field-ids>
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--fields` | []string | Field IDs to add (comma-separated) |

**Examples:**

```bash
orbit -p paybook jira screen field-add 10089 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399"
```

---

### screen field-remove

Remove fields from a screen tab.

```bash
orbit -p profile jira screen field-remove <screen-id> <tab-id> --fields <field-ids>
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--fields` | []string | Field IDs to remove (comma-separated) |

**Examples:**

```bash
orbit -p paybook jira screen field-remove 10089 10189 \
  --fields "customfield_10397,customfield_10398"
```

---

### screen field-move

Move fields from one tab to another on the same screen. This removes the field from the source tab and adds it to the target tab.

```bash
orbit -p profile jira screen field-move <screen-id> <source-tab-id> <target-tab-id> --fields <field-ids>
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--fields` | []string | Field IDs to move (comma-separated) |

**Examples:**

```bash
# Move AI fields from General tab to AI Workflow tab
orbit -p paybook jira screen field-move 10089 10189 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399,customfield_10400,customfield_10401,customfield_10402,customfield_10403"
```

---

## Filter Commands

### filter list

List favourite/saved filters. Alias: `ls`.

```
orbit -p profile jira filter list
```

### filter view

View details of a saved filter.

```
orbit -p profile jira filter view <filter-id>
```

**Examples:**

```bash
orbit -p paybook jira filter view 10195
```

### filter search

Search for filters by name. Returns all accessible filters when no name is specified.

```
orbit -p profile jira filter search [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | | Filter by name |
| `--max-results` | 50 | Maximum results |

**Examples:**

```bash
orbit -p paybook jira filter search
orbit -p paybook jira filter search --name "Sprint"
```

### filter create

Create a saved filter.

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Filter name |
| `--jql` | Yes | JQL query |
| `--description` | No | Description |
| `--favourite` | No | Mark as favourite |

```bash
orbit -p paybook jira filter create --name "My Bugs" --jql "project = PROJ AND type = Bug"
orbit -p paybook jira filter create --name "My Work" --jql "assignee = currentUser()" --favourite
```

### filter update

Update an existing filter.

| Flag | Description |
|------|-------------|
| `--name` | New name |
| `--jql` | New JQL query |
| `--description` | New description |

```bash
orbit -p paybook jira filter update 10195 --jql "project = PROJ AND sprint in openSprints()"
orbit -p paybook jira filter update 10195 --name "Active Sprint"
```

### filter delete

Delete a saved filter. Alias: `rm`.

```bash
orbit -p paybook jira filter delete 12345
```

---

## Dashboard Commands

### dashboard list

List dashboards visible to the current user. Alias: `ls`.

```
orbit -p profile jira dashboard list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--filter` | | Filter by name |
| `--max-results` | 50 | Maximum results |

**Examples:**

```bash
orbit -p paybook jira dashboard list
orbit -p paybook jira dashboard list --filter "Sprint"
```

### dashboard view

View details of a dashboard.

```
orbit -p profile jira dashboard view <dashboard-id>
```

**Examples:**

```bash
orbit -p paybook jira dashboard view 10117
```

### dashboard create

Create a new dashboard.

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | Yes | Dashboard name |
| `--description` | No | Dashboard description |

```bash
orbit -p paybook jira dashboard create --name "AI Metrics" --description "AI adoption tracking"
```

### dashboard delete

Delete a dashboard. Alias: `rm`.

```bash
orbit -p paybook jira dashboard delete 10408
```

### dashboard gadget list

List gadgets on a dashboard.

```bash
orbit -p paybook jira dashboard gadget list 10408
```

### dashboard gadget add

Add a gadget to a dashboard. Use `--uri` for Jira Cloud.

| Flag | Description |
|------|-------------|
| `--module-key` | Gadget module key |
| `--uri` | Gadget URI (alternative to module-key, required for Cloud) |
| `--title` | Gadget title |
| `--color` | Color: blue, red, yellow, green, cyan, purple, gray, white |
| `--column` | Column position (0-based) |
| `--row` | Row position (0-based) |

```bash
orbit -p paybook jira dashboard gadget add 10408 \
  --uri "rest/gadgets/1.0/g/com.atlassian.jira.gadgets:filter-results-gadget/gadgets/filter-results-gadget.xml" \
  --title "My Filter" --color blue
```

### dashboard gadget remove

Remove a gadget from a dashboard. Alias: `rm`.

```bash
orbit -p paybook jira dashboard gadget remove 10408 35501
```

### dashboard gadget property list

List properties of a gadget.

```bash
orbit -p paybook jira dashboard gadget property list 10408 35501
```

### dashboard gadget property get

Get a gadget property value (JSON output).

```bash
orbit -p paybook jira dashboard gadget property get 10408 35501 config
```

### dashboard gadget property set

Set a gadget property (JSON value).

```bash
orbit -p paybook jira dashboard gadget property set 10408 35501 config '{"filterId":"11233","num":"20"}'
```

---

## Status Commands

### status list

List all workflow statuses with their categories (To Do, In Progress, Done).

```bash
orbit -p profile jira status list
```

**Examples:**

```bash
orbit -p paybook jira status list
# Output: 10062  In Code Review  [In Progress]
#         10022  In QA           [In Progress]
#         ...

# Filter for specific statuses
orbit -p paybook jira status list | grep -iE "ai review|ready to deploy|blocked"
```

---

## Issue Type Commands

### issuetype-list

List all issue types available in the Jira instance.

```bash
orbit -p profile jira issuetype-list
```

**Examples:**

```bash
orbit -p paybook jira issuetype-list
# Output: 10001  Story
#         10002  Task
#         10003  Sub-task
#         10004  Bug
#         10000  Epic
```

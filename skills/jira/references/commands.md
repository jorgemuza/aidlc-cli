# orbit Jira Commands Reference

Complete reference for all `orbit jira` commands with flags and examples.

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
  - [issue link](#issue-link)
  - [issue unlink](#issue-unlink)
  - [issue worklog](#issue-worklog)
  - [issue clone](#issue-clone)
- [Epic Commands](#epic-commands)
- [Sprint Commands](#sprint-commands)
- [Board Commands](#board-commands)
- [Project Commands](#project-commands)
- [Release Commands](#release-commands)
- [Field Commands](#field-commands)

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

**Examples:**

```bash
orbit -p epsilon jira issue list --project PRT
orbit -p epsilon jira issue list --assignee me --status "In Progress"
orbit -p epsilon jira issue list --jql "project = PRT AND issuetype = Epic"
orbit -p epsilon jira issue list --project PRT --type Epic -o json
```

---

### issue view

View detailed issue information.

**Aliases:** `show`

```bash
orbit -p profile jira issue view <issue-key> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--comments` | int | Number of recent comments to show (default: 1) |

**Examples:**

```bash
orbit -p epsilon jira issue view PRT-4373
orbit -p epsilon jira issue view PRT-4373 --comments 5
orbit -p epsilon jira issue view PRT-4373 -o json
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
| `-b, --body` | string | Issue description (use wiki markup!) |
| `-y, --priority` | string | Priority |
| `-a, --assignee` | string | Assignee username |
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
orbit -p epsilon jira issue create --project PRT --type Story \
  --summary "Add login page"

# Epic with Parent Link
orbit -p epsilon jira issue create --type Epic --project PRT \
  --summary "Okta Authentication Foundation" \
  --priority Highest \
  --field "customfield_27521=PRT-4378"

# Bug with labels and components
orbit -p epsilon jira issue create --project PRT --type Bug \
  --summary "Fix timeout" --priority High \
  --label backend --label urgent --component api
```

**Notes:**
- `--parent` only works for Sub-task types. For linking Epics to Initiatives, use `--field "customfield_27521=KEY"` (Parent Link).
- When `--type Epic` is used, the Epic Name custom field (`customfield_11523`) is auto-set from `--summary`. Override with `--epic-name`.

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
| `-b, --body` | string | New description (use wiki markup!) |
| `-y, --priority` | string | New priority |
| `-l, --label` | []string | Add/remove labels (prefix with `-` to remove) |
| `-C, --component` | []string | Add/remove components (prefix with `-` to remove) |
| `--fix-version` | []string | Add/remove fix versions (prefix with `-` to remove) |

**Examples:**

```bash
orbit -p epsilon jira issue edit PRT-123 --summary "Updated title"
orbit -p epsilon jira issue edit PRT-123 --priority Critical
orbit -p epsilon jira issue edit PRT-123 --label new-label --label -old-label
orbit -p epsilon jira issue edit PRT-123 --body "h2. Updated Description

* New bullet point
* Another point"
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
orbit -p epsilon jira issue assign PRT-123 john.doe
orbit -p epsilon jira issue assign PRT-123 x    # unassign
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
orbit -p epsilon jira issue move PRT-123 "In Progress"
orbit -p epsilon jira issue move PRT-123 Done --comment "Fixed in v2.1"
orbit -p epsilon jira issue move PRT-123 Done --resolution Fixed
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

Add a comment to an issue.

```bash
orbit -p profile jira issue comment <issue-key> <body...>
```

**Examples:**

```bash
orbit -p epsilon jira issue comment PRT-123 "This is fixed now"
```

---

### issue link

Link two issues together.

```bash
orbit -p profile jira issue link <inward-key> <outward-key> <link-type>
```

**Examples:**

```bash
orbit -p epsilon jira issue link PRT-100 PRT-200 Blocks
orbit -p epsilon jira issue link PRT-100 PRT-200 Duplicates
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
orbit -p epsilon jira issue worklog PRT-123 "2h 30m"
orbit -p epsilon jira issue worklog PRT-123 "1d" --comment "Code review"
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
orbit -p epsilon jira issue clone PRT-123
orbit -p epsilon jira issue clone PRT-123 --summary "Cloned: new title"
orbit -p epsilon jira issue clone PRT-123 --replace "v1:v2"
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
orbit -p epsilon jira epic list --project PRT
orbit -p epsilon jira epic list PRT-50    # issues in epic PRT-50
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
| `-b, --body` | string | Epic description (use wiki markup!) |
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

### field-list

List Jira fields. Essential for discovering custom field IDs.

```bash
orbit -p profile jira field-list [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--filter` | string | Filter fields by name (case-insensitive) |

**Examples:**

```bash
orbit -p epsilon jira field-list --filter "parent"
orbit -p epsilon jira field-list --filter "epic"
orbit -p epsilon jira field-list --filter "client"
```

**Common custom fields (instance-specific):**

| Field | Typical ID | Usage |
|-------|-----------|-------|
| Epic Name | customfield_11523 | Required for Epic creation (auto-set by orbit) |
| Parent Link | customfield_27521 | Links epics to initiatives/capabilities |
| Epic Link | customfield_11522 | Links stories to epics |

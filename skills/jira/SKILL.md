---
name: jira
description: "Interact with Jira using the orbit CLI to create, list, view, edit, and transition issues, manage sprints and epics, export epic hierarchies to markdown, manage dashboards and gadgets, manage saved filters, manage custom fields and screen configurations, list statuses and issue types, and write properly formatted descriptions using Jira wiki markup. Use this skill whenever the user asks about Jira tasks, tickets, issues, sprints, epics, dashboards, filters, gadgets, or needs to manage project work items using orbit. Also trigger when the user says things like 'create a ticket', 'create epics', 'move this to done', 'assign the issue', 'update the description', 'format for Jira', 'create a custom field', 'add field to screen', 'list statuses', 'configure Jira', 'create a dashboard', 'add a gadget', 'list filters', 'search filters', 'export epic', 'dump epic hierarchy', or any Jira-related workflow — even casual references like 'update Jira', 'what tickets are in this sprint', 'add a comment to PROJ-123', 'set up AI tracking fields', 'show me the dashboards', or 'create a metrics dashboard'. Trigger especially when descriptions need proper formatting (headings, bullets, tables, links) since Jira Server uses wiki markup, not markdown."
---

# Jira with orbit CLI

Manage Jira issues, epics, sprints, boards, projects, releases, custom fields, screens, statuses, and issue types using the `orbit` CLI. Supports both Jira Cloud (REST API v3) and Jira Server/Data Center (REST API v2) with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. A profile with a `jira-cloud` or `jira-onprem` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials — API token for Cloud (basic auth with email + token), PAT for Server
4. Credentials can be stored as 1Password (`op://`) or Infisical (`infisical://`) references for automatic resolution

## Quick Reference

All commands follow the pattern: `orbit -p <profile> jira <resource> <action> [flags]`

The `-o` flag controls output format: `table` (default), `json`, `yaml`.

For full command details and flags, see `references/commands.md`.
For Jira wiki markup formatting (Server only), see `references/wiki-markup.md`.

## Core Workflows

### Viewing Issues

```bash
# View issue with comments (shows last 10 by default)
orbit -p myprofile jira issue view PROJ-123

# View with more comments
orbit -p myprofile jira issue view PROJ-123 --comments 50

# JSON output (includes full data + all comments)
orbit -p myprofile jira issue view PROJ-123 -o json
```

The view command shows: summary, type, status, priority, assignee, reporter, labels, components, description, subtasks, links, **attachments** (with filename, size, download URL), and **comments with author/date/edit timestamps**. Attachments appear before comments so screenshots referenced in comment text have context.

### Creating Issues

```bash
# Simple story
orbit -p myprofile jira issue create --project PROJ --type Story --summary "Add login page"

# Bug with priority and assignment
orbit -p myprofile jira issue create --project PROJ --type Bug \
  --summary "Fix timeout" --priority High --assignee john.doe

# Sub-task under a parent
orbit -p myprofile jira issue create --project PROJ --type Sub-task \
  --parent PROJ-123 --summary "Add validation"
```

### Creating Epics

Epics require the "Epic Name" custom field. The `orbit` CLI auto-sets this from the summary when `--type Epic` is used. Use `--epic-name` to override.

```bash
# Create an epic (Epic Name auto-set from summary)
orbit -p myprofile jira issue create --type Epic --project PROJ \
  --summary "Q1 Auth Revamp" --priority Highest

# Create epic with Parent Link to a Capability/Initiative
orbit -p myprofile jira issue create --type Epic --project PROJ \
  --summary "Okta Authentication Foundation" \
  --field "customfield_27521=PRT-4378" \
  --priority Highest
```

The `--field` flag sets arbitrary custom fields as `key=value`. Use `jira field list` to discover field IDs.

### Editing Issues with Formatted Descriptions

Jira Server uses **wiki markup**, not Markdown. Jira Cloud uses ADF (handled automatically by orbit for plain text). Always format `--body` values using wiki markup syntax for Server instances.

```bash
orbit -p myprofile jira issue edit PRT-123 --body "h2. Value Statement

The platform provides real authentication via Okta.

h2. User Stories

* *Story 1:* As a platform admin, I want all API requests to require a valid Okta JWT.
** Okta custom AS configured with epsilon_claims claim
** AUTH_BYPASS removed from all Terraform task definitions"
```

### Editing Custom Field Values

The `-F` flag supports multiple value formats. The correct format depends on the Jira field type:

```bash
# Single-select field — use value: shorthand
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10397=value:Yes

# Multi-select field — MUST use JSON array of objects (value: shorthand will return 400)
orbit -p myprofile jira issue edit PROJ-123 -F 'customfield_10398=[{"value":"Developer"},{"value":"QA"}]'

# Multi-select with a single value still needs array format
orbit -p myprofile jira issue edit PROJ-123 -F 'customfield_10400=[{"value":"Confirmed"}]'

# Number field
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10399=3

# Plain string field
orbit -p myprofile jira issue edit PROJ-123 -F customfield_10010=my-value
```

**Important**: Multi-select fields (checkboxes, multi-select dropdowns) require `[{"value":"X"}]` array format. Using `value:X` on a multi-select field returns a 400 error.

### Transitioning Issues

```bash
# Move to In Progress
orbit -p myprofile jira issue move PROJ-123 "In Progress"

# Close with comment and resolution
orbit -p myprofile jira issue move PROJ-123 Done --comment "Fixed in v2.1" --resolution Fixed
```

### Comments with @mentions

Use `[~accountId:ID]` to @mention users in comments (Cloud). Jira renders these as clickable user links. For Server, use `[~username]`.

```bash
# Add a comment mentioning the assignee
orbit -p myprofile jira issue comment PROJ-123 -b "Hey [~accountId:5b10ac8d82e05b22cc7d4ef5], this needs rework"

# Look up a user's accountId
orbit -p myprofile jira user search "Jorge"
```

### Attachments

```bash
# Attach a file to an issue
orbit -p myprofile jira issue attach PROJ-123 report.pdf

# Attach multiple files
orbit -p myprofile jira issue attach PROJ-123 screenshot.png log.txt

# List attachments
orbit -p myprofile jira issue attachments PROJ-123

# Download all attachments
orbit -p myprofile jira issue download PROJ-123

# Download only screenshots/images
orbit -p myprofile jira issue download PROJ-123 --images-only --output ./screenshots

# Remove an attachment by ID
orbit -p myprofile jira issue detach 12345
```

### Searching Issues

```bash
# Filter by project and type
orbit -p myprofile jira issue list --project PROJ --type Epic

# Filter by assignee and status
orbit -p myprofile jira issue list --assignee me --status "In Progress"

# Raw JQL query
orbit -p myprofile jira issue list --jql "project = PROJ AND sprint in openSprints()"

# JSON output for processing
orbit -p myprofile jira issue list --project PROJ -o json

# Fresh list — bypasses the JQL search index (useful right after transitions)
orbit -p myprofile jira issue list --jql "project = PROJ AND status = Open" --fresh
```

> **JQL index lag** — `issue list` hits Jira's JQL search index, which is eventually consistent and can trail transitions by 30–60s. When the user asks for a "current", "fresh", or "live" list, or they just transitioned a ticket in the same session, **pass `--fresh`** so each row is re-verified against the authoritative single-issue endpoint. Without `--fresh`, list results are informational; always re-verify a single key with `jira issue view KEY` before acting on it.

## Exporting Epic Hierarchies

Export full epic hierarchies (epic → stories → sub-tasks) to structured markdown directories.

```bash
# Export a single epic
orbit -p myprofile jira export BS-1243 -O ./jira-export

# Export multiple epics
orbit -p myprofile jira export BS-1243 BS-1255 BS-1246 -O ./jira-export
```

Output structure:
```
<outdir>/<EPIC-KEY>/INDEX.md                    — epic metadata + story list
<outdir>/<EPIC-KEY>/<STORY-KEY>/INDEX.md        — story metadata + sub-task list
<outdir>/<EPIC-KEY>/<STORY-KEY>/<TASK-KEY>.md   — individual sub-task file
```

Each file includes a metadata table (type, status, assignee, priority, labels, story points) and the full description converted from Jira ADF to markdown. Use `--max-results` to control how many children are fetched per level (default: 100).

## Sprint Management

```bash
# List sprints on a board
orbit -p myprofile jira sprint list --board-id 42
orbit -p myprofile jira sprint list --board-id 42 --state active

# View sprint details
orbit -p myprofile jira sprint view 1366

# List issues in a sprint
orbit -p myprofile jira sprint list 1366

# Create a sprint
orbit -p myprofile jira sprint create --board-id 42 --name "Sprint 5" --start-date 2026-04-01 --end-date 2026-04-15

# Start a sprint
orbit -p myprofile jira sprint start 1366

# Close a sprint
orbit -p myprofile jira sprint close 1366

# Update sprint name/dates/goal
orbit -p myprofile jira sprint update 1366 --name "Sprint 5 Extended" --end-date 2026-04-20

# Add issues to a sprint
orbit -p myprofile jira sprint add 1366 PROJ-101 PROJ-102

# Move issues to backlog
orbit -p myprofile jira sprint remove PROJ-101 PROJ-102

# Delete a sprint
orbit -p myprofile jira sprint delete 1366
```

## Dashboards & Gadgets

Create and manage Jira dashboards with gadgets programmatically.

```bash
# List dashboards
orbit -p myprofile jira dashboard list
orbit -p myprofile jira dashboard list --filter "Sprint"

# Create a dashboard
orbit -p myprofile jira dashboard create --name "AI Metrics" --description "AI adoption tracking"

# View dashboard details
orbit -p myprofile jira dashboard view 10408

# Delete a dashboard
orbit -p myprofile jira dashboard delete 10408

# Add gadgets to a dashboard (use URI format for Cloud)
orbit -p myprofile jira dashboard gadget add 10408 \
  --uri "rest/gadgets/1.0/g/com.atlassian.jira.gadgets:filter-results-gadget/gadgets/filter-results-gadget.xml" \
  --title "My Filter Results" --color blue

# List gadgets on a dashboard
orbit -p myprofile jira dashboard gadget list 10408

# Remove a gadget
orbit -p myprofile jira dashboard gadget remove 10408 35501

# Configure gadget properties (e.g., bind a filter)
orbit -p myprofile jira dashboard gadget property set 10408 35501 config '{"filterId":"11233"}'
orbit -p myprofile jira dashboard gadget property get 10408 35501 config
orbit -p myprofile jira dashboard gadget property list 10408 35501
```

## Filters

Create, search, and manage saved JQL filters.

```bash
# List favourite filters
orbit -p myprofile jira filter list

# Search all accessible filters
orbit -p myprofile jira filter search
orbit -p myprofile jira filter search --name "Sprint"

# View filter details
orbit -p myprofile jira filter view 10195

# Create a filter
orbit -p myprofile jira filter create --name "My Bugs" --jql "project = PROJ AND type = Bug" --favourite

# Update a filter
orbit -p myprofile jira filter update 10195 --jql "project = PROJ AND sprint in openSprints()"
orbit -p myprofile jira filter update 10195 --name "Active Sprint"

# Delete a filter
orbit -p myprofile jira filter delete 12345
```

## Field Management (Cloud Only)

Manage custom fields, their contexts, and options programmatically.

### Listing Fields

```bash
# List all fields
orbit -p myprofile jira field list

# List only custom fields
orbit -p myprofile jira field list --custom

# Filter by name or ID
orbit -p myprofile jira field list --filter "AI"
```

### Creating Custom Fields

```bash
# Create a select field (shorthand types: select, multiselect, number, checkbox, text, textarea)
orbit -p myprofile jira field create --name "AI Assisted" --type select \
  --description "Was AI used on this ticket?"

# Create a number field
orbit -p myprofile jira field create --name "AI Prompt Iterations" --type number \
  --description "How many prompt cycles to get working output"

# Create a checkbox field
orbit -p myprofile jira field create --name "Human Review Confirmed" --type checkbox \
  --description "Engineer confirms AI output was reviewed"
```

### Managing Field Contexts and Options

```bash
# List field contexts
orbit -p myprofile jira field context-list customfield_10397

# List existing options
orbit -p myprofile jira field option-list customfield_10397 10817

# Add options to a select/multiselect field
orbit -p myprofile jira field option-add customfield_10397 10817 \
  --values "Yes,No,Partial"
```

## Screen Management

Control which fields appear on issue create/edit/view forms by managing screen tabs and field assignments.

### Listing Screens and Tabs

```bash
# List all screens
orbit -p myprofile jira screen list

# Filter screens by name
orbit -p myprofile jira screen list --filter "PYMT"

# List tabs on a screen
orbit -p myprofile jira screen tab-list 10089

# List fields on a screen tab
orbit -p myprofile jira screen field-list 10089 10189
```

### Creating Tabs

```bash
# Create a new tab to group related fields
orbit -p myprofile jira screen tab-create 10089 "AI Workflow"
```

### Adding and Removing Fields

```bash
# Add fields to a screen tab
orbit -p myprofile jira screen field-add 10089 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399"

# Remove fields from a screen tab
orbit -p myprofile jira screen field-remove 10089 10189 \
  --fields "customfield_10397,customfield_10398"
```

### Moving Fields Between Tabs

```bash
# Move fields from General tab to AI Workflow tab
orbit -p myprofile jira screen field-move 10089 10189 10868 \
  --fields "customfield_10397,customfield_10398,customfield_10399"
```

## Status and Issue Type Management

### Listing Workflow Statuses

```bash
# List all statuses with their categories
orbit -p myprofile jira status list

# Filter with grep
orbit -p myprofile jira status list | grep -i "review"
```

### Listing Issue Types

```bash
orbit -p myprofile jira issuetype-list
```

## Deleting Issues

Delete issues with `issue delete`. Use `--cascade` to also delete all subtasks.

**Important:** There is no `--confirm` or `--force` flag. The command deletes immediately.

```bash
# Delete a single issue
orbit -p myprofile jira issue delete PROJ-123

# Delete an issue and all its subtasks
orbit -p myprofile jira issue delete PROJ-100 --cascade
```

**Deleting a hierarchy (parent → children → grandchildren):** Delete bottom-up to avoid orphan errors, or use `--cascade` at each level that has subtasks.

```bash
# Example: delete Story (with sub-tasks) → Epic → Capability
orbit -p myprofile jira issue delete PROJ-50 --cascade   # Story + sub-tasks
orbit -p myprofile jira issue delete PROJ-40              # Epic (now childless)
orbit -p myprofile jira issue delete PRT-200              # Capability (now childless)
```

**Tip:** `--cascade` only deletes Jira sub-tasks (child issues of type Sub-task). It does NOT delete Stories under an Epic or Epics under a Capability. For full hierarchy deletion, query children with JQL first, then delete bottom-up.

```bash
# Find all children of an issue
orbit -p myprofile jira issue list --jql '"Parent" = PROJ-40' -o json
```

## Important Notes

- **Cloud vs Server** — Use service type `jira-cloud` for Atlassian Cloud (uses API v3 with ADF for descriptions). Use `jira-onprem` for Server/Data Center (uses API v2 with wiki markup).
- **Description formatting** — For Server instances, always use Jira wiki markup for `--body` values. For Cloud, orbit auto-converts plain text to ADF. See `references/wiki-markup.md` for wiki markup syntax.
- **Field management** — `field create`, `field context-list`, `field option-list`, and `field option-add` are Cloud-only features.
- **Screen management** — Adding fields to screens makes them appear on issue create/edit forms. Use `tab-create` to group related fields into a dedicated tab.
- **Epic type cannot use `--parent` flag** — Jira rejects it because Epic is not a sub-task type. Use the `--field "customfield_27521=KEY"` (Parent Link) instead.
- **Secret references** - Credentials in config can use 1Password (`op://vault/item/field`) or Infisical (`infisical://<env>/<path>/<KEY>`) references, resolved at runtime. Run `orbit auth` once to resolve and cache all secrets (a single biometric prompt for 1Password). Use `orbit auth clear` to wipe the cache. See [Secrets](../../docs/secrets.md). Without `orbit auth`, secrets are still resolved on each command but may trigger repeated biometric prompts.

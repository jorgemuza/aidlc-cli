---
name: jira-skill
description: "Interact with Jira using the orbit CLI to create, list, view, edit, and transition issues, manage sprints and epics, and write properly formatted descriptions using Jira wiki markup. Use this skill whenever the user asks about Jira tasks, tickets, issues, sprints, epics, or needs to manage project work items using orbit. Also trigger when the user says things like 'create a ticket', 'create epics', 'move this to done', 'assign the issue', 'update the description', 'format for Jira', or any Jira-related workflow — even casual references like 'update Jira', 'what tickets are in this sprint', or 'add a comment to PROJ-123'. Trigger especially when descriptions need proper formatting (headings, bullets, tables, links) since Jira Server uses wiki markup, not markdown."
---

# Jira with orbit CLI

Manage Jira issues, epics, sprints, boards, projects, and releases using the `orbit` CLI. This tool connects to Jira Server/Data Center via REST API with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` binary built and accessible
2. A profile with a `jira-onprem` service configured in `~/.config/orbit/config.yaml`
3. A valid Jira personal access token (can be stored in 1Password with `op://` prefix)

## Quick Reference

All commands follow the pattern: `orbit -p <profile> jira <resource> <action> [flags]`

The `-o` flag controls output format: `table` (default), `json`, `yaml`.

For full command details and flags, see `references/commands.md`.
For Jira wiki markup formatting, see `references/wiki-markup.md`.

## Core Workflows

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

The `--field` flag sets arbitrary custom fields as `key=value`. Use `jira field-list` to discover field IDs.

### Editing Issues with Formatted Descriptions

Jira Server uses **wiki markup**, not Markdown. Always format `--body` values using wiki markup syntax.

```bash
orbit -p myprofile jira issue edit PRT-123 --body "h2. Value Statement

The platform provides real authentication via Okta.

h2. User Stories

* *Story 1:* As a platform admin, I want all API requests to require a valid Okta JWT.
** Okta custom AS configured with epsilon_claims claim
** AUTH_BYPASS removed from all Terraform task definitions
** Unauthenticated requests return 401

h2. Functional Requirements

||Req ID||Description||Priority||
|FR-001|Okta custom authorization server setup|Must|
|FR-002|AWS Secrets Manager wiring|Must|"
```

### Discovering Custom Fields

```bash
# List all fields, filter by name
orbit -p myprofile jira field-list --filter "parent"
orbit -p myprofile jira field-list --filter "epic"

# Common custom field IDs (vary by Jira instance):
# customfield_11523 = Epic Name (auto-set for Epic type)
# customfield_27521 = Parent Link (links epics to initiatives/capabilities)
# customfield_11522 = Epic Link (links stories to epics)
```

### Transitioning Issues

```bash
# Move to In Progress
orbit -p myprofile jira issue move PROJ-123 "In Progress"

# Close with comment and resolution
orbit -p myprofile jira issue move PROJ-123 Done --comment "Fixed in v2.1" --resolution Fixed
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
```

### Bulk Epic Creation Pattern

When creating multiple epics under a parent (Initiative/Capability), use this pattern:

```bash
# Create epics with Parent Link and formatted descriptions
for epic in "Epic 1 Name" "Epic 2 Name" "Epic 3 Name"; do
  orbit -p myprofile jira issue create --type Epic --project PROJ \
    --summary "$epic" \
    --field "customfield_27521=PROJ-100" \
    --priority Highest
done

# Then update descriptions with wiki markup formatting
orbit -p myprofile jira issue edit PROJ-201 --body "h2. Value Statement
..."
```

## Important Notes

- **Epic type cannot use `--parent` flag** — Jira rejects it because Epic is not a sub-task type. Use the `--field "customfield_27521=KEY"` (Parent Link) instead.
- **Description formatting** — Always use Jira wiki markup for `--body` values. Markdown will render as plain text. See `references/wiki-markup.md` for the full syntax.
- **Custom field values** — Most custom fields accept string values. For option/multi-select fields, you may need JSON values via the REST API directly.
- **1Password integration** — Token fields in config can use `op://vault/item/field` and are resolved at runtime.

---
name: azure-devops
description: "Manage Azure DevOps (VSTS) projects, work items, and saved queries using the orbit CLI. Use this skill whenever the user asks about Azure DevOps work items, boards, queries, WIQL, bugs, tasks, user stories, sprints, iterations, or project management on Azure DevOps / VSTS. Trigger on phrases like 'list work items', 'show bugs', 'run query', 'create a bug', 'update work item', 'ADO bugs', 'VSTS query', 'Azure DevOps tasks', 'what's assigned to me', 'list projects', 'saved queries', 'work item 12345', or any Azure DevOps-related task — even casual references like 'check the board', 'what bugs are open', 'show me the Fusion backlog', 'ADO status', 'query results', or mentions of WIQL. The orbit CLI alias is `ado`."
---

# Azure DevOps with orbit CLI

Manage Azure DevOps (VSTS) projects, work items, and saved queries from the command line. Supports both modern (`dev.azure.com`) and legacy (`*.visualstudio.com`) instances using PAT-based basic auth.

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
2. A profile with an `azure-devops` service configured in `~/.config/orbit/config.yaml`
3. A Personal Access Token (PAT) with appropriate scopes (Work Items Read/Write, Project Read)
4. Auth method: `basic` with your email as username and the PAT as password

## Configuration

```yaml
profiles:
  myprofile:
    services:
      - name: azure-devops
        type: azure-devops
        base_url: https://myorg.visualstudio.com   # or https://dev.azure.com/myorg
        proxy: none                                 # bypass profile-level proxy if needed
        auth:
          method: basic
          username: me@company.com
          password: "MY_PAT_HERE"                   # supports op:// references
        options:
          default_project: MyProject                # optional — avoids --project on every call
```

## Quick Reference

All commands follow the pattern: `orbit -p <profile> ado <command> [flags]`

All commands support `-o json` and `-o yaml` for structured output.

## Command Groups

| Group | Alias | Description |
|-------|-------|-------------|
| `project` | | List projects in the organization |
| `work-item` | `wi` | View, list, create, update, comment on work items |
| `query` | | Run and list saved queries |
| `version` | | Verify connectivity and show authenticated user |

## Core Workflows

### Verifying Connection

```bash
orbit -p myprofile ado version
```

### Listing Projects

```bash
orbit -p myprofile ado project list
orbit -p myprofile ado project list -o json
```

### Work Items

```bash
# View a single work item
orbit -p myprofile ado wi view 12345
orbit -p myprofile ado wi view 12345 -o json

# List work items with filters (builds WIQL automatically)
orbit -p myprofile ado wi list --project Fusion --type Bug --state New
orbit -p myprofile ado wi list --project Fusion --type "User Story" --state Active
orbit -p myprofile ado wi list --project Fusion --assigned-to "me@company.com"
orbit -p myprofile ado wi list --project Fusion --max-results 100

# Raw WIQL query
orbit -p myprofile ado wi list --project Fusion --wiql "SELECT [System.Id], [System.Title] FROM WorkItems WHERE [System.State] = 'Active' ORDER BY [System.CreatedDate] DESC"

# Create a work item
orbit -p myprofile ado wi create --project Fusion --type Bug \
  --field "System.Title=Login page broken" \
  --field "System.Description=Steps to reproduce..."

# Update a work item (any field — uses "replace" to avoid triggering state-reset rules)
orbit -p myprofile ado wi update 12345 --field "System.State=Active"
orbit -p myprofile ado wi update 12345 --field "System.AssignedTo=someone@company.com"
orbit -p myprofile ado wi update 12345 \
  --field "System.State=Resolved" \
  --field "Microsoft.VSTS.Common.ResolvedReason=Fixed"

# Add a discussion comment (preserves current state automatically)
orbit -p myprofile ado wi comment 12345 -m "Root cause: null check in auth handler"
orbit -p myprofile ado wi comment 12345 -m "Fixed in commit abc123. Deployed to staging."
```

> **Comments:** always use `wi comment -m "..."` instead of `wi update --field "System.History=..."`. The comment command fetches the current state and re-sends it in the same PATCH, preventing Azure DevOps work item rules from silently resetting the state.

### Saved Queries

```bash
# Run a saved query by its UUID
orbit -p myprofile ado query run a670954b-e739-47ca-a09a-acc10f623123 --project Fusion

# List saved queries and folders
orbit -p myprofile ado query list --project Fusion
orbit -p myprofile ado query list --project Fusion --depth 3
```

## Common Field Names

Azure DevOps uses fully-qualified field names. The most common:

| Field | Description |
|-------|-------------|
| `System.Title` | Work item title |
| `System.State` | State (New, Active, Resolved, Closed) |
| `System.WorkItemType` | Type (Bug, Task, User Story, Epic, Feature) |
| `System.AssignedTo` | Assigned user (email or display name) |
| `System.Description` | HTML description |
| `System.AreaPath` | Area path (e.g. `Fusion\Backend`) |
| `System.IterationPath` | Iteration/sprint path |
| `System.Tags` | Semicolon-separated tags |
| `Microsoft.VSTS.Common.Priority` | Priority (1-4) |
| `Microsoft.VSTS.Common.Severity` | Severity (1-4) |
| `Microsoft.VSTS.Common.ResolvedReason` | Why it was resolved |

## WIQL (Work Item Query Language)

WIQL is Microsoft's query language for work items, similar to SQL:

```sql
SELECT [System.Id], [System.Title], [System.State]
FROM WorkItems
WHERE [System.TeamProject] = 'Fusion'
  AND [System.WorkItemType] = 'Bug'
  AND [System.State] <> 'Closed'
  AND [System.AssignedTo] = @Me
ORDER BY [System.CreatedDate] DESC
```

Key differences from Jira JQL:
- Field names use brackets: `[System.Title]`
- String values use single quotes: `'Active'`
- Current user: `@Me`
- Not equal: `<>`
- Table: `FROM WorkItems` (always)

## Important Notes

- **PAT auth**: Azure DevOps uses the PAT as the password in basic auth. The username is your email.
- **API version**: All requests use `api-version=7.1` (latest stable).
- **JSON Patch**: Create/update use `application/json-patch+json` format internally — the CLI handles this; just pass `--field Key=Value`.
- **default_project**: Set in `options.default_project` to avoid `--project` on every command.
- **Proxy bypass**: If your profile has a proxy for VPN services but Azure DevOps is on the public internet, set `proxy: none` on the service to bypass.
- **Secret references**: PAT in config can use 1Password (`op://vault/item/field`) or Infisical (`infisical://<env>/<path>/<KEY>`) references, resolved at runtime via `orbit auth`.

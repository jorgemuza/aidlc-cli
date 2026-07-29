---
name: qmetry
description: "Manage QMetry test cases, test cycles, and test executions using the orbit CLI. Use this skill whenever the user asks about QMetry test management, creating test cases, listing test cases, test cycles, test execution, QMetry projects, or test planning. Trigger on phrases like 'create test case', 'list test cases', 'test cycle', 'QMetry', 'QTC-', 'create tests in QMetry', 'test management', 'test plan', or any QMetry-related task — even casual references like 'add tests for this ticket', 'create QMetry cases from AC', or 'what test cases exist'. The orbit CLI alias is `qm`."
---

# QMetry with orbit CLI

Manage QMetry test cases, test cycles, folders, and projects through the `orbit` CLI. Works with QMetry Cloud via REST API with multi-profile support and secret resolution (1Password and Infisical).

## Prerequisites

1. `orbit` CLI installed — if `which orbit` fails, install with:
   - **macOS/Linux (Homebrew):** `brew install jorgemuza/tap/orbit`
   - **macOS/Linux (script):** `curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh`
   - **Windows (Scoop):** `scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket && scoop install orbit`
2. A profile with a `qmetry` service configured in `~/.config/orbit/config.yaml`
3. QMetry API key - generated from QMetry Settings → API Keys; can be stored as 1Password (`op://`) or Infisical (`infisical://<env>/<path>/<KEY>`) references

## Configuration

```yaml
services:
  - name: qmetry
    type: qmetry
    base_url: https://qtmcloud.qmetry.com
    headers:
      referer: "https://qtmcloud.qmetry.com/"
    auth:
      method: apiKey
      apiKey: "op://Vault/qmetry-api-key/credential"
```

The `apiKey` auth method sends the key as an `apiKey` HTTP header on every request. The `referer` header is required by QMetry's CORS validation. API keys are generated from QMetry Settings → API Keys.

## Quick Reference

```bash
# List projects
orbit -p myprofile qm project

# List folders in a project
orbit -p myprofile qm folder --project-id 10061

# List test cases
orbit -p myprofile qm tc list --project-id 10061
orbit -p myprofile qm tc list --project-id 10061 --folder-id 2382094

# View a test case
orbit -p myprofile qm tc view QTC-101

# Create a test case
orbit -p myprofile qm tc create --project-id 10061 --summary "TC: Login validation [PYMT-123]"

# Create with BDD steps
orbit -p myprofile qm tc create --project-id 10061 --folder-id 2382094 \
  --summary "TC: Payment flow [PYMT-456]" \
  --description "Validate credit card payment processing" \
  --steps '[{"stepDetails":"Given user has items in cart","expectedResult":"Cart total displayed"},{"stepDetails":"When user enters valid card","testData":"4111111111111111","expectedResult":"Payment processed"}]'

# List test cycles
orbit -p myprofile qm cycle list --project-id 10061

# Create a test cycle
orbit -p myprofile qm cycle create --project-id 10061 --name "Sprint 5 Regression"
```

## Project & Folder Mapping

| Jira Project | QMetry projectId | Default folderId | Notes |
|---|---|---|---|
| PYMT | `10061` | `2382094` | Syncfy payment flows |

Use `orbit qm project` and `orbit qm folder --project-id ID` to discover IDs.

## Workflow: Creating Test Cases from Jira Acceptance Criteria

1. Read the Jira ticket to extract acceptance criteria
2. For each AC, create a QMetry test case with BDD steps (Given/When/Then)
3. Use the returned QTC keys as references in automation scripts

```bash
# Step 1: Read the ticket
orbit -p paybook jira issue view PYMT-123

# Step 2: Create test cases
orbit -p paybook qm tc create --project-id 10061 --folder-id 2382094 \
  --summary "TC: User can login with valid credentials [PYMT-123]" \
  --steps '[{"stepDetails":"Given user is on login page","expectedResult":"Login form visible"},{"stepDetails":"When user enters valid credentials","testData":"user@example.com / password123","expectedResult":""},{"stepDetails":"Then user is redirected to dashboard","expectedResult":"Dashboard loads with user name"}]'
```

For full command details and flags, see `references/commands.md`.

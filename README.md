# orbit

A unified CLI for managing connections to development lifecycle services.

Supports **Jira**, **Confluence**, **GitLab**, **GitHub**, **Forgejo**/**Gitea**, **Bitbucket**, **GoCD**, and **Azure DevOps** (cloud and self-hosted). Includes **build provenance attestation** verification via Sigstore and **AI token usage tracking** (tkm) with cost analysis per model. Pairs with [RTK](https://github.com/rtk-ai/rtk) for command output compression (60-90% token savings). Organize connections into profiles to switch between projects seamlessly.

Secrets can be stored as [1Password](https://1password.com/) references (`op://vault/item/field`) or [Infisical](https://infisical.com/) references (`infisical://<env>/<path>/<KEY>`) and are resolved at runtime. See [Secrets](docs/secrets.md).

## Install

### macOS / Linux — Homebrew

```bash
brew install jorgemuza/tap/orbit
```

Installing the formula directly (as above) trusts only this one formula and needs no
`brew tap`. If you instead tap the whole repo (`brew tap jorgemuza/tap`), recent
Homebrew versions print a "tap is not trusted" warning — mark it trusted once with:

```bash
brew trust jorgemuza/tap
```

### Linux / macOS — Install script

```bash
curl -sSfL https://raw.githubusercontent.com/jorgemuza/orbit/main/install.sh | sh
```

This detects your OS and architecture, downloads the correct binary from GitHub Releases, and installs it to `/usr/local/bin`.

### Windows — Scoop

```bash
scoop bucket add jorgemuza https://github.com/jorgemuza/scoop-bucket
scoop install orbit
```

### Binary releases

Download pre-built binaries for all platforms from [GitHub Releases](https://github.com/jorgemuza/orbit/releases).

## Verifying releases

Releases are built reproducibly by GitHub Actions and the checksums file is signed,
so you can confirm a binary really came from this project:

- **Signed checksums (cosign).** `checksums.txt` is signed keylessly with
  [cosign](https://docs.sigstore.dev/), as a single Sigstore bundle
  (`checksums.txt.sigstore.json`) carrying both the signature and the certificate.
  Download both files from the release, verify the signature, then verify your
  download's checksum against the file:

  ```bash
  cosign verify-blob checksums.txt \
    --bundle checksums.txt.sigstore.json \
    --certificate-identity-regexp 'https://github.com/jorgemuza/orbit-cli/.*' \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com
  shasum -a 256 -c checksums.txt --ignore-missing
  ```

  Releases up to and including **v0.66.0** were signed in the older two-file form
  instead. Verify those with `--signature checksums.txt.sig --certificate
  checksums.txt.pem` in place of `--bundle`, and note that cosign v4 removed those
  flags, so verifying them needs a cosign 3.x.

- **Build provenance (SLSA).** The release workflow also emits GitHub build-provenance
  attestations — verifiable with `gh attestation verify <file> --repo jorgemuza/orbit-cli`.
  GitHub only persists these for public repositories, so they're available once the
  build repo is public; until then, use the cosign verification above.

> Note: these are out-of-band guarantees you opt into — Homebrew does not yet
> auto-verify attestations for third-party taps. The `brew trust` step above is a
> separate, consumer-side acknowledgement of *who* you trust to run formula code.

## Quick Start

```bash
# Create a profile
orbit profile create --name my-project --default

# Add services
orbit profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method basic --username me@co.com --token "op://Dev/jira-token/credential"

orbit profile add-service \
  --name gitlab-onprem --type gitlab --variant server \
  --base-url https://gitlab.internal.com \
  --auth-method basic --username admin --password "op://Dev/gitlab/password"

orbit profile add-service \
  --name github-cloud --type github \
  --auth-method token --token "op://Dev/github-token/credential"

# Test connectivity
orbit service ping
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `~/.config/orbit/config.yaml` |
| `-p, --profile` | Profile to use (overrides default) | |
| `-o, --output` | Output format: `table`, `json`, `yaml` | `table` |
| `--debug` | Print HTTP request/response details to stderr | `false` |

## Request Reliability

Every orbit service client (jira, confluence, github, forgejo, gitlab, bitbucket, gocd, draxarp, qmetry, kroki, tkm) shares a single HTTP layer with automatic retry for transient failures:

- **Idempotent methods only** — `GET` and `HEAD` retry up to 3 times; `POST`, `PUT`, `PATCH`, `DELETE` execute exactly once so transitions/creates/updates can't be double-submitted.
- **Exponential backoff** — 150ms → 300ms → 600ms between attempts.
- **Retry triggers** — network/connection errors, `429 Too Many Requests`, and `5xx` (except `501 Not Implemented`).
- **Fail-fast** — `4xx` other than `429` returns immediately (e.g. `401`, `403`, `404`, `422`) since retrying won't help.
- Use `--debug` to see each attempt's request/response on stderr.

**One path deliberately does not use this layer: rendering a diagram.** `confluence publish` probes the Kroki renderer directly, with its own schedule - 4 attempts, **1s → 4s → 9s** backoff, inside a per-diagram budget of 30s that covers every attempt and every wait (`--timeout`, or a kroki connection's `http_timeout`, overrides that budget). It also retries **`403`**, which the shared layer fails fast on, because that is how the public `https://kroki.io` signals Cloudflare rate limiting rather than a real refusal. A kroki service connection's own client - what `orbit service ping` uses - is an ordinary service client and does share the layer above.

Bulk read commands that fan out over many keys — such as `orbit jira issue list --fresh` — share a bounded-concurrency worker (`internal/concurrency.ParallelDo`) so parallel fetches respect service rate limits (default concurrency: 4 for Jira Cloud) and collect partial failures instead of aborting the whole operation.

## Commands Overview

### Profile & Service Management

```bash
orbit profile create --name proj --default         # Create a profile
orbit profile list                                  # List profiles
orbit profile show proj                             # Show profile details
orbit profile use proj                              # Set default profile
orbit profile delete proj                           # Delete profile
orbit profile add-service --name svc --type jira ...  # Add service
orbit profile remove-service svc                    # Remove service

orbit service list                                  # List services in active profile
orbit service ping                                  # Test all connections
orbit service ping jira-cloud                       # Test specific service

orbit auth                                          # Pre-resolve and cache secrets
orbit auth clear                                    # Clear cached secrets

orbit version                                       # Print version info
```

### Jira

Manage issues, epics, sprints, boards, projects, releases, dashboards, filters, custom fields, screens, statuses, and issue types. Supports Jira Cloud (API v3) and Server/Data Center (API v2).

```bash
orbit jira issue list --project PROJ --assignee me
orbit jira issue create --project PROJ --type Story --summary "Add login"
orbit jira issue edit PROJ-123 --field customfield_10397="Yes"
orbit jira issue move PROJ-123 "In Progress"
orbit jira epic list --project PROJ
orbit jira sprint list --board-id 1 --state active
orbit jira field list --custom
orbit jira field create --name "AI Assisted" --type select
orbit jira screen list --filter "PYMT"
orbit jira screen tab-create 10089 "AI Workflow"
orbit jira screen field-add 10089 10868 --fields "customfield_10397,customfield_10398"
orbit jira status list
orbit jira issuetype-list
```

**[Full Jira reference →](docs/jira.md)**

### Confluence

Manage pages, publish markdown directories, and control page layout. Supports Cloud and Server/Data Center.

```bash
orbit confluence page 12345
orbit confluence children 12345
orbit confluence hierarchy 12345 --depth 3
orbit confluence create --space DEV --parent 12345 --title "Guide" --file docs/guide.md
orbit confluence update 12345 --file docs/guide.md
orbit confluence publish ./docs --space DEV --parent 12345 --dry-run
orbit confluence set-width 12345 --recursive
```

**[Full Confluence reference →](docs/confluence.md)**

### GitLab

Manage projects, merge requests, pipelines, jobs, branches, tags, issues, runners, variables, and more. Alias: `gl`.

```bash
orbit gl project view my-team/app
orbit gl project create --name "my-app" --namespace my-team
orbit gl project list --group my-team
orbit gl mr list my-team/app --state opened
orbit gl mr create my-team/app --source feature/x --target main --title "Add feature"
orbit gl pipeline list my-team/app
orbit gl pipeline jobs my-team/app 12345
orbit gl job log my-team/app 67890
orbit gl branch list my-team/app
orbit gl variable list my-team/app
```

**[Full GitLab reference →](docs/gitlab.md)**

### GitHub

Manage repositories, pull requests, issues, Actions workflow runs, releases, secrets, and more. Alias: `gh`.

```bash
orbit gh repos --org my-org
orbit gh pr list octocat/hello-world
orbit gh pr create octocat/hello-world --head feature/x --base main --title "Add feature"
orbit gh run list octocat/hello-world
orbit gh run watch octocat/hello-world 12345
orbit gh run view octocat/hello-world 12345 --log-failed
orbit gh issue list octocat/hello-world --labels bug
orbit gh secret set octocat/hello-world MY_SECRET --value "s3cret"
orbit gh release latest octocat/hello-world
```

**[Full GitHub reference →](docs/github.md)**

### Forgejo

Manage a self-hosted Forgejo (or Gitea) instance: repositories, pull requests and reviews, issues, Actions runs and job logs, releases. Alias: `fj`.

```bash
orbit fj repo view jorgemuza/orbit-cli
orbit fj pr list jorgemuza/orbit-cli
orbit fj pr approve jorgemuza/orbit-cli 17
orbit fj pr request-changes jorgemuza/orbit-cli 17 -m "Three things to fix"
orbit fj pr merge jorgemuza/orbit-cli 17 --method squash --delete-branch
orbit fj run list jorgemuza/orbit-cli --status failure
orbit fj run log jorgemuza/orbit-cli 263 --failed
orbit fj issue list jorgemuza/orbit-cli
orbit fj api /api/v1/repos/jorgemuza/orbit-cli/branches --paginate
```

Runs are addressed by the number in the web UI. The API path takes a different id, and asking it for the visible number returns a *different* run with a 200 - so orbit resolves the number and verifies it got the run you asked for. Reviews and merges are likewise checked against what the server stored rather than trusting the 2xx.

**[Full Forgejo reference →](docs/forgejo.md)**

### Bitbucket

Manage projects, repositories, pull requests, branches, tags, reviewer conditions, and users. Supports Cloud and Data Center. Alias: `bb`.

```bash
orbit bb repo list MY-PROJ
orbit bb pr list MY-PROJ my-repo --state open
orbit bb pr create MY-PROJ my-repo --source feature/x --target main --title "Add feature"
orbit bb pr approve MY-PROJ my-repo 42
orbit bb branch list MY-PROJ my-repo
orbit bb reviewer-condition list MY-PROJ
```

**[Full Bitbucket reference →](docs/bitbucket.md)**

### GoCD

Manage pipelines, pipeline groups, agents, environments, config repos, server administration, users, roles, authorization, plugins, backups, materials, artifact stores, elastic agent profiles, cluster profiles, stages, jobs, server configuration, templates, packages, notification filters, access tokens, secret configs, and dashboard. Alias: `cd`.

```bash
orbit cd pipeline list
orbit cd pipeline status my-pipeline
orbit cd pipeline trigger my-pipeline
orbit cd pipeline compare my-pipeline --from 1 --to 5
orbit cd agent list
orbit cd env list
orbit cd dashboard
orbit cd template list
orbit cd access-token list
orbit cd secret-config list
orbit cd version
```

**[Full GoCD reference →](docs/gocd.md)**

### Attestation

Verify build provenance attestations using Sigstore bundles with in-toto format and SLSA provenance predicates.

```bash
orbit attestation verify ./my-binary --bundle attestation.jsonl
orbit attestation verify ./artifact --bundle bundle.json --owner my-org --repo my-repo --signer-identity "https://github.com/my-org/my-repo"
orbit attestation inspect attestation.jsonl
orbit attestation download sha256:abc123... --repo owner/repo
```

**[Full Attestation reference →](docs/attestation.md)**

## Supported Services

| Service | Type | Variants | Default Base URL |
|---------|------|----------|-----------------|
| Jira | `jira` | `cloud`, `server` | *(required)* |
| Confluence | `confluence` | `cloud`, `server` | *(required)* |
| GitLab | `gitlab` | `cloud`, `server` | `https://gitlab.com` |
| GitHub | `github` | `cloud`, `server` | `https://api.github.com` |
| Forgejo / Gitea | `forgejo` | — | *(required — always self-hosted)* |
| Bitbucket | `bitbucket` | `cloud`, `server` | `https://api.bitbucket.org/2.0` |
| GoCD | `gocd` | — | *(required — always self-hosted)* |
| Azure DevOps | `azure-devops` | `cloud` | *(required)* — `dev.azure.com/{org}` or `{org}.visualstudio.com` |
| QMetry | `qmetry` | `cloud` | *(required)* |
| Kroki | `kroki` | — | *(required on the connection - but the connection is optional: it is the diagram renderer `confluence publish` uses, and with none configured it falls back to `https://kroki.io`)* |

### Authentication Methods

| Method | Flags | Use Case |
|--------|-------|----------|
| `token` | `--token` | API tokens, PATs (most common) |
| `basic` | `--username`, `--password` | Basic auth (Jira Cloud uses email + API token) |
| `oauth2` | `--client-id`, `--client-secret` | OAuth2 client credentials |

## 1Password Integration

Instead of storing secrets in plain text, use 1Password references:

```bash
orbit profile add-service \
  --name jira-cloud --type jira --variant cloud \
  --base-url https://myco.atlassian.net \
  --auth-method basic \
  --username me@company.com \
  --password "op://DevVault/jira-token/credential"
```

Secrets are resolved at runtime via `op read`, so the [1Password CLI](https://developer.1password.com/docs/cli/) must be installed and authenticated. Resolved values are cached locally to avoid repeated biometric prompts.

Pre-resolve secrets before running multiple commands:

```bash
orbit auth          # triggers biometric once, caches secrets
orbit auth clear    # remove cached secrets
```

By default, cached secrets expire after **8 hours**. You can change this (or disable expiration entirely) in `config.yaml`:

```yaml
settings:
  secrets_cache_ttl_hours: 0   # 0 = never expire (default: 8)
```

## Configuration

Config is stored in YAML at `~/.config/orbit/config.yaml`:

```yaml
settings:
  secrets_cache_ttl_hours: 8   # 0 = never expire (default: 8)

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
          method: basic
          username: me@company.com
          password: "op://DevVault/jira-token/credential"
      - name: gitlab-onprem
        type: gitlab
        variant: server
        base_url: https://gitlab.internal.com
        tls_skip_verify: true                          # skip TLS cert validation (self-signed certs)
        proxy: socks5://127.0.0.1:1080                 # optional: SOCKS5, HTTP, or HTTPS proxy
        auth:
          method: basic
          username: admin
          password: "op://DevVault/gitlab/password"
      - name: github-cloud
        type: github
        auth:
          method: token
          token: "op://DevVault/github-token/credential"
```

### Network Configuration

Proxy and TLS settings can be set at the **profile level** (applies to all services) or per **service** (overrides profile):

```yaml
profiles:
  - name: corp
    proxy: socks5h://127.0.0.1:1080       # applies to all services in this profile
    tls_skip_verify: true                   # applies to all services in this profile
    services:
      - name: gitlab
        type: gitlab
        base_url: https://gitlab.internal.com
        # inherits proxy + tls_skip_verify from profile
      - name: jira-cloud
        type: jira
        base_url: https://myco.atlassian.net
        proxy: ""                            # override: no proxy for this service
        tls_skip_verify: false               # override: verify TLS for this service
```

| Field | Level | Description |
|-------|-------|-------------|
| `proxy` | Profile or Service | Proxy URL: `socks5://`, `socks5h://`, `http://`, `https://` |
| `tls_skip_verify` | Profile or Service | Skip TLS cert validation (self-signed certs) |

Service-level settings override profile-level. Without any proxy config, orbit respects `HTTP_PROXY`/`HTTPS_PROXY`/`ALL_PROXY` environment variables.

## Claude Code Plugin

Orbit is available as a [Claude Code](https://docs.anthropic.com/en/docs/claude-code) plugin that gives Claude deep knowledge of each service's CLI commands, flags, and workflows.

### Install via marketplace

```
/plugin marketplace add jorgemuza/claude-plugins
/plugin install orbit@jorgemuza
```

### Install directly

```
/plugin install jorgemuza/orbit
```

### Install individual skills (alternative)

```bash
npx @anthropic-ai/claude-code-skills --skills jira,confluence,github,forgejo,gitlab,bitbucket,gocd,format-docs --from github:jorgemuza/orbit
```

### Available Skills

| Skill | Description |
|-------|-------------|
| `jira` | Issue CRUD, epics, sprints, boards, dashboards, filters, fields, screens, statuses, wiki markup formatting |
| `confluence` | Page view/create/update, hierarchy, markdown publishing, page width control |
| `github` | Repos, PRs, Actions runs, issues, releases, secrets (alias: `gh`) |
| `gitlab` | Projects, MRs, pipelines, branches, tags, variables (alias: `gl`) |
| `bitbucket` | Projects, repos, PRs, branches, tags, reviewer conditions, approvals (alias: `bb`) |
| `forgejo` | Repos, PRs and reviews, issues, Actions runs and job logs, releases on Forgejo/Gitea (alias: `fj`) |
| `gocd` | Pipelines, agents, environments, config repos, server admin, stages, jobs (alias: `cd`) |
| `draxarp` | Projects, tasks, specs, docs, memories, sprints, knowledge graph, context captures |
| `attestation` | Verify, download, and inspect build provenance attestations (Sigstore/SLSA) |
| `format-docs` | Prepare and restructure markdown for Confluence publishing |

## License

MIT

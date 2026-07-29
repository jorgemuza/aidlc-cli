# Secrets

orbit resolves secret references in your config at runtime, so credentials never
have to be stored as plaintext. A value is only treated as a secret if it starts
with a known scheme prefix; everything else is passed through as a literal.

Two backends are supported:

| Backend | Reference format | Resolved via |
|---------|------------------|--------------|
| 1Password | `op://vault/item/field` | the 1Password CLI (`op`), gated by Touch ID on macOS |
| Infisical | `infisical://<env>/<path>/<KEY>` | the Infisical Go SDK with Universal Auth (machine identity) |

Any config value that accepts a credential can use either scheme - `auth.token`,
`auth.password`, `auth.client_secret`, `auth.apiKey`, `auth.username`, and header
values.

## Caching

Resolved secrets are cached to `~/.config/orbit/.secret-cache` (mode `0600`) so
repeated commands do not re-resolve. The TTL is controlled by
`settings.secrets_cache_ttl_hours` (default: 8 hours; `0` disables expiration).

```bash
orbit auth          # resolve every reference in the config and cache them
orbit auth clear    # remove the cache
```

For 1Password, `orbit auth` batches all references into a single `op inject`
call so only one biometric prompt is needed.

## 1Password (`op://`)

```yaml
services:
  - name: jira-cloud
    type: jira
    auth:
      method: basic
      username: me@example.com
      token: "op://Dev/jira-token/credential"
```

References are resolved by shelling out to the `op` CLI, which must be installed
and signed in.

## Infisical (`infisical://`)

Infisical is the right backend for headless hosts (no Touch ID): a machine
identity authenticates via Universal Auth, and the reference stays fully
symbolic, so `config.yaml` is safe to commit.

### Reference format

```
infisical://<env>/<path...>/<KEY>
```

- `<env>` - the Infisical environment slug (e.g. `prod`, `dev`).
- `<path...>` - the secret folder path. Optional; defaults to `/` when omitted.
- `<KEY>` - the secret name.

Examples:

```yaml
services:
  - name: jira
    type: jira
    auth:
      method: basic
      username: "infisical://prod/orbit/JIRA_USER"
      token: "infisical://prod/orbit/JIRA_TOKEN"   # env=prod, path=/orbit, key=JIRA_TOKEN
  - name: gitlab
    type: gitlab
    auth:
      method: token
      token: "infisical://prod/GITLAB_TOKEN"       # env=prod, path=/, key=GITLAB_TOKEN
```

### Connection configuration

The non-secret connection parameters live in config settings (safe to commit):

```yaml
settings:
  infisical:
    site_url: "https://secrets.example.com"   # omit for the Infisical SaaS default
    project_id: "your-project-id"
```

Universal Auth credentials are **never** stored in config. They are read from the
environment at runtime:

| Variable | Purpose |
|----------|---------|
| `INFISICAL_UNIVERSAL_AUTH_CLIENT_ID` | Universal Auth machine-identity client ID |
| `INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET` | Universal Auth machine-identity client secret |
| `INFISICAL_API_URL` | Overrides `settings.infisical.site_url` (default: `https://app.infisical.com`) |
| `INFISICAL_PROJECT_ID` | Overrides `settings.infisical.project_id` |

Precedence for site URL and project ID is: environment variable > config setting >
default.

### Injecting the bootstrap credentials

On hosts that already run [secretspec](https://secretspec.dev/), run orbit under
it so the machine identity is injected with no new secret written to disk:

```bash
secretspec run --profile prod -- orbit service ping
```

Any mechanism that exports `INFISICAL_UNIVERSAL_AUTH_CLIENT_ID` and
`INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET` into orbit's environment works. If they
are missing, resolution fails with a clear error.

# Orbit Attestation Command Reference

Verify, download, and inspect build provenance attestations using Sigstore bundles with in-toto attestation format and SLSA provenance predicates.

**Top-level command:** `orbit attestation` (alias: `attest`)

**Notes:**

- Attestation verification uses Sigstore bundle format.
- Provenance follows the in-toto attestation specification with SLSA predicates.
- All commands support `-o json` and `-o yaml` for structured output.
- **Signatures are not cryptographically verified.** See [Limitations](#limitations).

---

## Table of Contents

- [verify -- Verify artifact provenance](#verify)
- [download -- Download attestation bundle](#download)
- [inspect -- Inspect attestation bundle](#inspect)
- [Limitations](#limitations)

---

## verify

Verify an artifact's build provenance against its Sigstore attestation bundle.

Computes the artifact digest (or accepts a pre-computed digest), loads the attestation bundle, parses the SLSA provenance predicate, checks that the artifact is one of the attestation's subjects, and verifies the signer identity and source repository.

`Verification: PASSED` requires all of the following:

- The bundle contains a DSSE envelope.
- The envelope's `predicateType` is a SLSA provenance type (`v0.1`, `v0.2`, or `v1`).
- The artifact digest appears among the statement's `subject[].digest` entries, under the same algorithm. Every subject is checked, not just the first.
- If `--signer-identity` is given, the attestation names a signer that it matches as an *anchored prefix* (see [Matching semantics](#matching-semantics)).
- If `--owner` or `--repo` is given, the provenance's config source URI parses into a host, owner, and repository, and each of them matches what was asked for. The host defaults to `github.com` and is set with `--source-host`. Nested namespaces such as `mygroup/subgroup` are supported.

See [Limitations](#limitations) for what `PASSED` does *not* mean.

```
orbit attestation verify <artifact-path-or-digest> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--bundle` | Path to local attestation bundle file | |
| `--owner` | Org/owner (or nested namespace, e.g. `mygroup/subgroup`) that must own the provenance source repository | |
| `--repo` | Repository the provenance source must name, as `repo` or `owner/repo` | |
| `--signer-identity` | Expected signer workflow identity; matched as a prefix of the builder ID | |
| `--source-host` | Forge the provenance source must be hosted on; applies whenever `--owner` or `--repo` is used | `github.com` |
| `--digest-algorithm` | Hash algorithm: `sha256`, `sha512` | `sha256` |

**Examples:**

```bash
# Verify a local binary against a bundle
orbit attestation verify ./my-binary --bundle attestation.jsonl

# Verify with owner, repo, and signer identity checks
orbit attest verify ./artifact --bundle bundle.json --owner my-org --repo my-repo \
  --signer-identity "https://github.com/my-org/my-repo"

# Verify against a self-hosted forge rather than github.com
orbit attest verify ./artifact --bundle bundle.json --source-host git.example.com --repo my-org/my-repo

# Verify a project in a nested GitLab namespace
orbit attest verify ./artifact --bundle bundle.json --source-host gitlab.com \
  --owner mygroup/subgroup --repo myproject

# Verify a pre-computed digest
orbit attestation verify abc123def456... --bundle att.json --digest-algorithm sha256

# Output as JSON
orbit attestation verify ./my-binary --bundle att.json -o json
```

### Matching semantics

`--signer-identity`, `--source-host`, `--owner`, and `--repo` are all anchored.
None of them is a substring search over attacker-controlled text.

**`--signer-identity`** is an anchored prefix of the provenance's `builder.id`.
The value must match from the *start* of the identity and stop on a genuine
boundary - a `/`, the `@` that introduces a git ref, or the end of the string.
Those two characters are the whole boundary set: `-` and `.` are legal *inside*
a name, so treating them as boundaries would let `.../my-org` match a builder
under `.../my-org-evil`. This keeps the usual prefix form working while rejecting
an identity that merely embeds the expected value:

| `--signer-identity` | `builder.id` | Result |
|---|---|---|
| `https://github.com/my-org/my-repo` | `https://github.com/my-org/my-repo/.github/workflows/release.yml@refs/heads/main` | match |
| `https://github.com/my-org/my-repo` | `https://github.com/my-org/my-repo` | match |
| `https://github.com/my-org/my-repo` | `https://evil.example.com/https://github.com/my-org/my-repo/hosted` | **no match** |
| `https://github.com/my-org/my-re` | `https://github.com/my-org/my-repo` | **no match** |
| `https://github.com/my-org` | `https://github.com/my-org-evil/my-repo/...` | **no match** |

A builder ID is two things spliced together, and `--signer-identity` treats them
differently. Its **forge coordinate** - the `scheme://host/owner/repo` head - is
matched **case-insensitively**, exactly as `--owner` and `--repo` are below: it
names the same GitHub account and repository they do, and rejecting
`https://github.com/My-Org/My-Repo/...` while accepting `--owner my-org` in the
same command would reject an account's real release workflow over nothing but
spelling. Everything **after** the coordinate is matched **case-sensitively**: a
workflow file path and a git ref name a specific file on a specific branch, and
folding them would let a filter for `.../Release.yml@refs/heads/Main` be
satisfied by a different workflow on a different branch. An identity that is not
URI-shaped has no coordinate and is compared case-sensitively end to end.

| `--signer-identity` | `builder.id` | Result |
|---|---|---|
| `https://github.com/my-org/my-repo` | `https://github.com/My-Org/My-Repo/.github/workflows/release.yml@refs/heads/main` | match (coordinate folds) |
| `.../my-repo/.github/workflows/Release.yml` | `.../my-repo/.github/workflows/release.yml@...` | **no match** (tail is exact) |
| `https://github.com/my-org` | `https://github.com/My-Org-evil/my-repo/...` | **no match** (boundary still applies) |

**`--source-host`, `--owner`, and `--repo`** are compared against the
provenance's `invocation.configSource.uri`. That URI is *parsed* into a host, an
owner, and a repository, and each is then compared for equality with what you
asked for. Nothing is matched by position in the raw string, so an extra or
empty path segment cannot shift one field into another's place.

These source URI shapes are understood:

| Shape | Example |
|---|---|
| `scheme://host/owner/repo` | `https://github.com/my-org/my-repo` |
| VCS-qualified scheme | `git+https://github.com/my-org/my-repo@refs/tags/v1.0.0` |
| ssh with userinfo and port | `git+ssh://git@github.com:2222/my-org/my-repo@refs/heads/main` |
| scp-like | `git@github.com:my-org/my-repo.git` |

A trailing `@<ref>` and at most one trailing `.git` are stripped. The path must
have at least two segments and no empty ones; a shorter path, or one containing
an empty segment, is reported as unreadable rather than having a field guessed
at. Verification also fails if the provenance records no config source URI at
all.

**Nested namespaces** are supported, because they are GitLab's normal structure.
The *last* path segment is the repository and everything before it is the owner:

| Source URI | Owner | Repository |
|---|---|---|
| `git+https://github.com/my-org/my-repo@refs/heads/main` | `my-org` | `my-repo` |
| `git+https://gitlab.com/mygroup/subgroup/myproject@refs/heads/main` | `mygroup/subgroup` | `myproject` |

`--owner` matches the **whole** namespace, never just its first segment, so
`--owner mygroup` does *not* accept a project under `mygroup/subgroup` - use
`--owner mygroup/subgroup`. Matching the first segment alone would let any
subgroup, including one an attacker controls, satisfy a filter naming the parent.
The `--repo owner/repo` shorthand stays a single-slash GitHub shape; for a nested
namespace pass `--owner` and `--repo` separately.

**The host is checked too.** Using `--owner` or `--repo` constrains the forge as
well, defaulting to `github.com`. Without that, an attestation whose source URI
merely *mimics* the path on a host you never named would satisfy the filter:

```
git+ssh://git@evil.example.com/my-org/my-repo@refs/heads/main
```

is rejected under `--owner my-org --repo my-repo`, because its host is not
`github.com`. Override the expected forge with `--source-host git.example.com`;
given on its own, `--source-host` is a filter in its own right.

**Host, owner, and repository are compared case-insensitively.** DNS names are
case-insensitive, and so are GitHub owner and repository names, so
`git+https://github.com/My-Org/My-Repo@refs/heads/main` satisfies
`--owner my-org --repo my-repo`. Rejecting it would be rejecting a genuine
attestation. The fold covers ASCII letters only - an internationalised domain
reaches this code as punycode, and Unicode folding equates some non-ASCII
characters with ASCII ones, which an attacker-chosen owner could exploit. (Case
folding applies to these three and to the identity coordinate above;
`--digest-algorithm` and the `predicateType` check remain exact.)

So `--owner acme` does *not* accept
`git+https://github.com/mallory/acme-mirror@refs/heads/main`, and `--repo` is
enforced rather than ignored.

`--repo` accepts two forms:

| Form | Constrains |
|---|---|
| `--repo my-repo` | the repository segment only |
| `--repo my-org/my-repo` | the owner segment **and** the repository segment |

The owner-qualified form carries an owner of its own, so combining it with a
`--owner` that names a *different* owner is a caller error and is rejected up
front, rather than one of the two being silently preferred:

```
$ orbit attestation verify ./artifact --bundle b.json --owner someone-else --repo my-org/my-repo
Error: conflicting filters: owner "someone-else" and repo "my-org/my-repo" name different owners
("someone-else" and "my-org"); drop the owner filter or qualify the repo filter with the same owner
```

Passing the same owner in both is fine.

**Output (table format):**

```
Verification: PASSED
Digest:       sha256:abc123def456...
Signer:       https://github.com/my-org/my-repo/.github/workflows/release.yml@refs/tags/v1.0.0
Builder:      https://github.com/actions/runner
Build Type:   https://slsa.dev/provenance/v1
Source:       git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
Commit:       abc123def456
Materials:    3

Note: signatures are not cryptographically verified; this checks the
attestation's contents only. Use `gh attestation verify` for full Sigstore
verification.
```

---

## download

Download an attestation bundle for an artifact digest.

```
orbit attestation download <artifact-digest> [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--repo` | Repository (owner/repo) — **required** | |
| `--digest-algorithm` | Hash algorithm: `sha256`, `sha512` | `sha256` |

**Examples:**

```bash
# Download attestation bundle
orbit attestation download sha256:abc123... --repo owner/repo

# With explicit algorithm
orbit attest download abc123... --repo owner/repo --digest-algorithm sha256
```

---

## inspect

Display the contents of an attestation bundle, including SLSA provenance, signer identity, and build information.

```
orbit attestation inspect <bundle-file> [flags]
```

**Examples:**

```bash
# Inspect a bundle file
orbit attestation inspect attestation.jsonl

# Output as JSON for processing
orbit attest inspect bundle.json -o json
```

**Output (table format):**

```
Media Type:  application/vnd.dev.sigstore.bundle.v0.3+json
Payload:     application/vnd.in-toto+json
Signatures:  1
Signer:      https://github.com/my-org/my-repo/.github/workflows/release.yml@refs/tags/v1.0.0

Provenance:
  Builder:    https://github.com/actions/runner
  Build Type: https://slsa.dev/provenance/v1
  Source:     git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
  Entry:      .github/workflows/release.yml
  sha1:       abc123def456
  Materials:
    - git+https://github.com/my-org/my-repo@refs/tags/v1.0.0
      sha1: abc123def456
```

---

## Limitations

**No signature is cryptographically verified.**

`orbit attestation verify` parses the Sigstore bundle and checks what the
attestation *says*. It does not check that the attestation is *authentic*:

- The DSSE signatures over the payload are never checked. `signature_count` in
  the JSON output is the number of signatures present, not the number verified
  (which is always zero).
- The signing certificate is never validated against Fulcio, and no certificate
  identity extensions (workflow, issuer, repository) are checked. The
  `--signer-identity`, `--source-host`, `--owner`, and `--repo` checks are
  anchored matches (see [Matching semantics](#matching-semantics)), but they run
  against the SLSA predicate's `builder.id` and `invocation.configSource.uri`,
  which are unsigned payload content: whoever hands you the bundle chooses them.
  In particular, the source host check constrains what the *payload claims* the
  forge was, not where the bundle actually came from.
- Rekor is never consulted, so there is no transparency-log inclusion proof.

Anyone able to hand you a bundle file can therefore produce one that reports
`Verification: PASSED` for any artifact. Treat the result as a provenance
*inspection* with content assertions, not as a trust decision.

For a real trust decision, use the GitHub CLI, which performs full Sigstore
verification:

```bash
gh attestation verify ./my-binary --repo my-org/my-repo
```

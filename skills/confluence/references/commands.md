# orbit Confluence Commands Reference

Complete reference for all `orbit confluence` commands with flags and examples.

## Table of Contents

- [Global Flags](#global-flags)
- [page — View Page](#page)
- [search — Search Pages](#search)
- [children — List Child Pages](#children)
- [hierarchy — Show Page Hierarchy](#hierarchy)
- [create — Create Page](#create)
- [update — Update Page](#update)
- [move — Re-parent Page](#move)
- [export — Export Page](#export)
- [comment — List/Add/Delete Comments](#comment)
- [publish — Publish Directory](#publish)
- [set-width — Set Page Width](#set-width)

---

## Global Flags

These flags are available on all confluence subcommands:

| Flag | Description |
|------|-------------|
| `-p, --profile` | Profile to use (overrides default) |
| `-o, --output` | Output format: table, json, yaml (default: table) |
| `--service` | Confluence service name (if profile has multiple) |
| `--config` | Config file path (default: ~/.config/orbit/config.yaml) |

---

## page

View a Confluence page's details including title, version, space, and URL.

```bash
orbit -p profile confluence page <page-id> [flags]
```

**Examples:**

```bash
orbit -p paybook confluence page 473676972036
orbit -p paybook confluence page 473676972036 -o json
```

**Output (table format):**

```
ID:      473676972036
Title:   AI Development Process
Version: 3
Space:   FO
URL:     https://paybook.atlassian.net/wiki/spaces/FO/pages/473676972036
```

**Output (json format):** Full page object including `body.storage.value` with the page content in Confluence storage format.

---

## search

Search Confluence pages using CQL or convenience filters.

```
orbit confluence search [query] [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--space` | string | | Filter by space key |
| `--title` | string | | Search by title (fuzzy match, CQL `~`) |
| `--label` | string | | Filter by label |
| `--text` | string | | Full-text search |
| `--cql` | string | | Raw CQL query (overrides other filters) |
| `--limit` | int | 25 | Maximum results |

A positional argument is treated as a text search term.

```bash
orbit -p myprofile confluence search --space FO --title "Architecture"
orbit -p myprofile confluence search --space ISMS --limit 100
orbit -p myprofile confluence search --space FO --label design
orbit -p myprofile confluence search --cql 'space=FO AND label=design AND type=page'
orbit -p myprofile confluence search "deployment pipeline"
```

---

## children

List child pages of a given page.

```bash
orbit -p profile confluence children <page-id> [flags]
```

**Examples:**

```bash
orbit -p paybook confluence children 473677299713
orbit -p paybook confluence children 473677299713 -o json
```

**Output:**

```
ID              VERSION  TITLE
--              -------  -----
473676972036    3        AI Development Process
473677103107    2        Foundations
473676873739    1        Organization
```

---

## hierarchy

Show the full page hierarchy — ancestor chain and descendant tree.

```bash
orbit -p profile confluence hierarchy <page-id> [--depth N]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--depth` | 2 | Maximum depth for children tree |

**Examples:**

```bash
# Show hierarchy (default 2 levels deep)
orbit -p paybook confluence hierarchy 473676972036

# Deeper tree
orbit -p paybook confluence hierarchy 473676972036 --depth 5

# JSON output (structured tree)
orbit -p paybook confluence hierarchy 473676972036 -o json
```

**Output:**

Shows ancestors top-to-bottom, marks the target page, then renders children with tree characters.

---

## create

Create a new Confluence page. Automatically sets wide width on the created page.

```bash
orbit -p profile confluence create [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--title` | string | Page title (**required**) |
| `--space` | string | Space key (**required**) |
| `--parent` | string | Parent page ID |
| `-b, --body` | string | Page body in storage format (XHTML) |
| `-f, --file` | string | Markdown file to convert and upload |

**Examples:**

```bash
# From markdown file
orbit -p paybook confluence create --space FO --parent 473677299713 \
  --title "Sprint Ceremonies" --file docs/sprint-ceremonies.md

# With inline content
orbit -p paybook confluence create --space FO --parent 473677299713 \
  --title "Quick Note" --body "<p>Remember to update the docs</p>"
```

**Notes:**
- When `--file` is provided, the markdown is automatically converted to Confluence storage format
- The `--body` flag expects raw Confluence storage format (XHTML), not markdown
- Created pages are automatically set to full-width layout

---

## update

Update an existing Confluence page. Automatically increments the version number.

```bash
orbit -p profile confluence update <page-id> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--title` | string | New title (keeps current if empty) |
| `-b, --body` | string | New body in storage format (XHTML) |
| `-f, --file` | string | Markdown file to convert and upload |

**Examples:**

```bash
# Update content from markdown
orbit -p paybook confluence update 473676972036 --file docs/overview.md

# Update title and content
orbit -p paybook confluence update 473676972036 \
  --title "Updated Overview" --file docs/overview.md

# Update with inline content
orbit -p paybook confluence update 473676972036 \
  --body "<p>Updated content</p>"
```

**Notes:**
- The command automatically fetches the current page version and increments it
- If `--title` is not provided, the existing title is preserved
- When using `--file`, the title from the markdown file is NOT used — provide `--title` explicitly if you want to change it

---

## move

Re-parent a Confluence page under a different parent. Preserves the page body, title, and labels — only the parent changes. Useful when the repo layout shifts and the published page tree needs to follow.

```bash
orbit -p profile confluence move <page-id> --parent <new-parent-id>
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--parent` | string | New parent page ID (**required**) |

**Examples:**

```bash
# Move page under a new parent
orbit -p paybook confluence move 895422142 --parent 878970411

# JSON output (includes new version number and URL)
orbit -p paybook confluence move 895422142 --parent 878970411 -o json
```

**Notes:**
- The new parent must be in the same space — Confluence rejects cross-space moves at the API level
- Confluence rejects cycles (moving a page under one of its own descendants)
- The command fetches the current page first to read the version, then issues a `PUT /content/{id}` with the new `ancestors` field and an incremented version

---

## export

Export a Confluence page as markdown or raw storage format.

```
orbit confluence export [page-id] [flags]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `markdown` | Output format: `markdown`, `storage` |
| `--output` | string | | Output directory (prints to stdout if omitted) |

Markdown export includes frontmatter with `confluence_page_id`, `confluence_space`, and `confluence_url`. Storage export saves the raw Confluence XHTML.

```bash
# Export page as markdown to stdout
orbit -p myprofile confluence export 12345

# Export page as markdown to a directory
orbit -p myprofile confluence export 12345 --format markdown --output docs/

# Export raw storage format
orbit -p myprofile confluence export 12345 --format storage --output backup/
```

---

## comment

List, add, and delete comments on a page. Comment bodies are written in Markdown by default (converted to Confluence storage format).

```
orbit confluence comment list <page-id>
orbit confluence comment add <page-id> [flags]
orbit confluence comment delete <comment-id> [flags]
```

**`comment add` flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--body` | string | | Comment body (Markdown by default). Mutually exclusive with `--from-file`. |
| `--from-file` | string | | Read the comment body from a Markdown file. |
| `--storage` | bool | `false` | Treat the body as raw Confluence storage XHTML instead of Markdown. |
| `--reply-to` | string | | Comment ID to reply to (threads the comment under it). |
| `--inline-text` | string | | Anchor an inline comment to the first occurrence of this page text (Confluence Cloud, best-effort — rejected if the text isn't found). |

**`comment delete` flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-y`, `--yes` | bool | `false` | Skip the interactive confirmation prompt. |

`comment list` renders each comment as `[footer|inline] <id> by <author>` with a one-line body excerpt; use `-o json` for the full objects. Footer comments work on both Confluence Server/DC and Cloud; inline anchoring (`--inline-text`) is Cloud-only and best-effort — when unsure, omit it and quote the line in the comment body.

```bash
# List comments on a page
orbit -p myprofile confluence comment list 12345

# Add a comment (Markdown body)
orbit -p myprofile confluence comment add 12345 --body "LGTM, but the sandbox URL is stale."

# Add a longer review from a file
orbit -p myprofile confluence comment add 12345 --from-file review.md

# Reply to a comment
orbit -p myprofile confluence comment add 12345 --body "Fixed in v0.60.2" --reply-to 67890

# Inline comment anchored to page text (Cloud)
orbit -p myprofile confluence comment add 12345 --body "Inaccurate — deploys are on-demand" --inline-text "deploys run nightly"

# Delete a comment without prompting
orbit -p myprofile confluence comment delete 67890 --yes
```

---

## publish

Recursively publish a directory of markdown files to Confluence, preserving folder hierarchy.

```bash
orbit -p profile confluence publish <directory> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--space` | string | Confluence space key (**required**) |
| `--parent` | string | Parent page ID (**required**) |
| `--dry-run` | bool | Preview without creating pages |
| `--no-kroki` | bool | Disable Kroki diagram rendering for the whole run; leave diagram code blocks as preformatted text |

**Examples:**

```bash
# Preview the publish plan
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713 --dry-run

# Publish for real
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713

# Publish without rendering diagrams via Kroki (leave them as code blocks)
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713 --no-kroki
```

**Directory Structure Rules:**

```
docs/
├── INDEX.md          → becomes parent page for "docs"
├── overview.md       → child page under docs
├── quick-ref.md      → child page under docs
├── foundations/
│   ├── INDEX.md      → becomes parent page for "foundations" (under docs)
│   ├── concepts.md   → child page under foundations
│   └── setup.md      → child page under foundations
└── workflow/
    ├── INDEX.md      → becomes parent page for "workflow" (under docs)
    └── sprints.md    → child page under workflow
```

**Title Resolution Order:**
1. YAML frontmatter `title:` field
2. First `# heading` in the file
3. Filename converted to title case (e.g., `quick-ref.md` → "Quick Ref")

**Frontmatter Fields:**

The `publish` command recognizes the following YAML frontmatter fields:

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title (overrides filename) |
| `confluence_page_id` | string | Existing page ID for updates |
| `confluence_url` | string | Page URL (set after publish) |
| `confluence_labels` | list | Labels/tags applied to the page after create/update |
| `confluence_properties` | map | Page Properties macro prepended to page content |
| `confluence_disable_kroki` | bool | When `true`, skip Kroki rendering for this file; leave diagram code blocks as preformatted text |

### Labels

Add `confluence_labels` to frontmatter to tag pages with Confluence labels. Labels are applied via a POST to `/content/{id}/label` after each page is created or updated.

```yaml
---
title: "My Policy"
confluence_labels:
  - ai-process
  - foundation
---
```

### Page Properties (details macro)

Add `confluence_properties` to frontmatter to generate a Page Properties macro (`ac:structured-macro ac:name="details"`) prepended to the page content. This creates a structured metadata block that can be queried by Page Properties Report macros on other pages.

```yaml
---
title: "Compounding Engineering & System Evolution"
confluence_properties:
  id: status
  fields:
    Owner: AI Tooling Guild
    Classification: Internal
    Status: "{status:Green|approved}"
    Reviewed on: 2026-03-06
    Approved on: 2026-03-06
---
```

**Sub-keys:**

| Key | Description |
|-----|-------------|
| `id` | Macro ID used by `detailssummary` reports to target specific property blocks |
| `fields` | Ordered key-value pairs rendered as a two-column table inside the macro |

**Field value formats:**

| Format | Rendered As |
|--------|-------------|
| `{status:Color\|Text}` | Status badge macro (e.g., `{status:Green\|approved}`) |
| `YYYY-MM-DD` | Confluence `<time>` date macro |
| Plain text | Plain text |

### Page Properties Report (detailssummary macro)

Use either directive format in markdown to generate a dynamic table that pulls Page Properties from child or labeled pages.

**HTML comment format (full control):**
```markdown
<!-- confluence:properties-report cql="label = 'ai-process' and space = currentSpace()" firstcolumn="Document" headings="Status, Classification, Reviewed on" -->
```

**Shorthand format:**
```markdown
{properties-report: label="ai-process", columns="Status, Classification, Reviewed on"}
```

**Parameters:**

| Parameter | Description |
|-----------|-------------|
| `cql` | CQL query to find pages (e.g., `label = "policy" and space = currentSpace()`) |
| `label` | Shorthand for CQL: generates `label = "value" and space = currentSpace()` |
| `firstcolumn` | Name of the first column (defaults to "Title") |
| `columns` / `headings` | Comma-separated list of property names to show as columns |
| `sortBy` | Optional column to sort by |

**Notes:**
- Subdirectories are processed before sibling files
- Hidden directories (starting with `.`) are skipped
- Only `.md` files are processed
- Each `INDEX.md` creates a parent page; other `.md` files create child pages under it
- Labels are applied after page creation/update via the Confluence REST API
- Page Properties macros are prepended to the converted page body before upload
- Properties Report directives are converted inline during markdown-to-storage-format conversion
- `<!-- confluence:ignore-start -->` / `<!-- confluence:ignore-end -->` blocks are stripped — content between the markers is completely skipped during conversion, useful for static markdown tables that should not appear on Confluence

---

## set-width

Set the content width of one or more pages.

```bash
orbit -p profile confluence set-width <page-id...> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--width` | string | Page width: wide or fixed (default: wide) |
| `--recursive` | bool | Apply to all child pages recursively |

**Examples:**

```bash
# Set single page to wide
orbit -p paybook confluence set-width 473676972036

# Set multiple pages
orbit -p paybook confluence set-width 473676972036 473677103107

# Set page and all descendants to wide
orbit -p paybook confluence set-width 473677299713 --recursive

# Set to fixed width
orbit -p paybook confluence set-width 473676972036 --width fixed
```

**Notes:**
- Sets both draft and published appearance properties
- "wide" maps to `full-width` appearance, "fixed" maps to `fixed`
- With `--recursive`, traverses all child pages depth-first

# orbit Confluence Commands Reference

Complete reference for all `orbit confluence` commands with flags and examples.

## Table of Contents

- [Global Flags](#global-flags)
- [page — View Page](#page)
- [children — List Child Pages](#children)
- [create — Create Page](#create)
- [update — Update Page](#update)
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

**Examples:**

```bash
# Preview the publish plan
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713 --dry-run

# Publish for real
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713
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

**Notes:**
- Subdirectories are processed before sibling files
- Hidden directories (starting with `.`) are skipped
- Only `.md` files are processed
- Each `INDEX.md` creates a parent page; other `.md` files create child pages under it

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

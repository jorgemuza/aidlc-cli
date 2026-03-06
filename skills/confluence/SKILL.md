---
name: confluence-skill
description: "Manage Confluence pages using the orbit CLI — create, update, view, publish markdown directories, and control page width. Use this skill whenever the user asks about Confluence pages, wiki content, publishing documentation, uploading markdown to Confluence, syncing docs, or managing page hierarchies using orbit. Trigger on phrases like 'create a Confluence page', 'update the wiki', 'publish these docs to Confluence', 'upload markdown', 'set page width', 'view page', 'list child pages', or any Confluence-related task — even casual references like 'push this to Confluence', 'sync the docs', or 'check what pages are under X'. Also trigger when the user needs to convert markdown to Confluence storage format or wants to track which markdown files map to which Confluence pages via frontmatter metadata (confluence_page_id, confluence_url)."
---

# Confluence with orbit CLI

Manage Confluence pages, publish markdown documentation, and control page layout using the `orbit` CLI. Supports Confluence Cloud and Server/Data Center via REST API with multi-profile support and 1Password secret resolution.

## Prerequisites

1. `orbit` binary built and accessible
2. A profile with a `confluence-cloud` or `confluence-server` service configured in `~/.config/orbit/config.yaml`
3. Valid credentials (API token for Cloud, PAT for Server) — can be stored in 1Password with `op://` prefix
4. For Cloud: auth type is `basic` with email as username and API token as password

## Quick Reference

All commands follow the pattern: `orbit -p <profile> confluence <command> [flags]`

For full command details and flags, see `references/commands.md`.
For Confluence storage format (XHTML) details, see `references/storage-format.md`.

## Core Workflows

### Viewing Pages

```bash
# View page details
orbit -p myprofile confluence page 473676972036

# View as JSON (includes body content)
orbit -p myprofile confluence page 473676972036 -o json

# List child pages
orbit -p myprofile confluence children 473676972036
```

### Creating Pages from Markdown

When creating a page from a markdown file, the CLI automatically converts it to Confluence storage format (XHTML). The converter handles headings, lists, tables, code blocks, blockquotes, inline formatting, and images.

```bash
# Create page from markdown file
orbit -p myprofile confluence create --space FO --parent 473677299713 \
  --title "My New Page" --file docs/overview.md

# Create page with inline body (storage format XHTML)
orbit -p myprofile confluence create --space FO --parent 473677299713 \
  --title "Quick Page" --body "<p>Hello world</p>"
```

Pages are automatically created with **wide width** (full-width layout).

### Updating Existing Pages

```bash
# Update page content from markdown file
orbit -p myprofile confluence update 473676972036 --file docs/overview.md

# Update title and content
orbit -p myprofile confluence update 473676972036 \
  --title "Updated Title" --file docs/overview.md

# Update with inline storage format
orbit -p myprofile confluence update 473676972036 --body "<p>New content</p>"
```

### Publishing a Directory of Markdown Files

The `publish` command recursively converts an entire directory of markdown files to Confluence pages, preserving the folder hierarchy:

- `INDEX.md` files become the parent page for their directory
- Other `.md` files become child pages under the directory's parent
- Subdirectories are processed recursively
- Page titles come from frontmatter `title:`, first `# heading`, or filename

```bash
# Publish entire docs directory
orbit -p myprofile confluence publish ./docs --space FO --parent 473677299713

# Preview what would be created (no API calls)
orbit -p myprofile confluence publish ./docs --space FO --parent 473677299713 --dry-run
```

### Setting Page Width

```bash
# Set single page to wide
orbit -p myprofile confluence set-width 473676972036

# Set page and all children to wide
orbit -p myprofile confluence set-width 473676972036 --recursive

# Set to fixed width
orbit -p myprofile confluence set-width 473676972036 --width fixed
```

## Markdown Frontmatter Tracking

When publishing or converting markdown files for Confluence, track the mapping between local files and Confluence pages using YAML frontmatter. After creating or updating a page, update the source markdown file's frontmatter with:

```yaml
---
title: "My Page Title"
confluence_page_id: "473676972036"
confluence_url: "https://mysite.atlassian.net/wiki/spaces/FO/pages/473676972036/My+Page+Title"
---
```

This enables:
- Re-running updates without needing to look up page IDs
- Tracking which files have been published
- Building scripts that sync local changes to Confluence

**When the user asks to publish or sync markdown files to Confluence:**
1. Run the orbit command to create/update the page
2. Capture the page ID and URL from the output
3. Update the markdown file's frontmatter with `confluence_page_id` and `confluence_url`
4. If the frontmatter doesn't exist, add it at the top of the file

**Example workflow for syncing a file:**

```bash
# If confluence_page_id exists in frontmatter, update:
orbit -p myprofile confluence update 473676972036 --file docs/overview.md

# If no confluence_page_id, create new:
orbit -p myprofile confluence create --space FO --parent 473677299713 \
  --title "Overview" --file docs/overview.md
# Then update the frontmatter with the returned page ID and URL
```

## Markdown-to-Confluence Conversion

The converter handles the following transformations:

| Markdown | Confluence Storage Format |
|----------|--------------------------|
| `# Heading` (first one) | Skipped — Confluence shows the page title |
| `## Heading` → `###### Heading` | `<h2>` → `<h6>` |
| `**bold**` | `<strong>` |
| `` `code` `` | `<code>` |
| `~~strike~~` | `<del>` |
| `- bullet` / `* bullet` | `<ul><li>` |
| `1. numbered` | `<ol><li>` |
| `> blockquote` | Info panel macro |
| ` ```lang ` code blocks | Code macro with language |
| `\| table \| rows \|` | `<table>` with `<thead>` |
| `![alt](url)` | `<ac:image>` |
| `[text](url)` | `<a href>` (relative .md links stripped) |
| `---` horizontal rules | Ignored (not rendered in Confluence) |
| YAML frontmatter | Stripped |
| `**Key**: Value` metadata lines | Two-column table (gray label column) |
| `## Table of Contents` section | Replaced with Confluence `toc` macro (no bullets) |

**Key behaviors to communicate to users:**
- The first `# heading` is always skipped because Confluence already displays the page title
- Relative markdown links (`.md`, `./`) are converted to plain text since they'll be separate Confluence pages
- A "Table of Contents" heading (or "Contents" / "TOC") replaces the entire section with Confluence's built-in TOC macro (bullet style: none)
- Code blocks and tables use full-width layout
- Horizontal rules (`---`) are ignored
- Document metadata lines (`**Key**: Value`) right after frontmatter are converted to a styled two-column table with gray labels. Example in markdown:
  ```
  **Document Version**: 2.0
  **Status**: Published
  **Classification**: Internal
  **Last Updated**: March 5, 2026
  ```
- Code blocks with language hints get syntax highlighting via the code macro

## Important Notes

- **Wide width by default** — All pages created via `orbit` are automatically set to full-width layout. Use `set-width --width fixed` to revert.
- **Cloud vs Server** — Use service type `confluence-cloud` for Atlassian Cloud (requires `/wiki/` prefix in API paths, handled automatically). Use `confluence-server` for Data Center.
- **Auth for Cloud** — Basic auth with your email as username and an API token (not your password) as the password field.
- **1Password integration** — Credentials in config can use `op://vault/item/field` and are resolved at runtime.
- **Dry run before publish** — Always use `--dry-run` first when publishing a directory to preview the page hierarchy before making API calls.

# Confluence Command Reference

Command reference for `orbit confluence` -- manage Confluence pages from the terminal.

## Table of Contents

- [Global Flags](#global-flags)
- [page](#page) -- View a Confluence page
- [children](#children) -- List child pages
- [hierarchy](#hierarchy) -- Show page hierarchy (ancestors + descendants)
- [create](#create) -- Create a new page
- [update](#update) -- Update an existing page
- [delete](#delete) -- Delete a page
- [export](#export) -- Export a page as markdown or storage
- [comment](#comment) -- List, add, and delete page comments
- [publish](#publish) -- Publish a directory of markdown files
- [set-width](#set-width) -- Set page width
- [Notes](#notes)
- [Markdown Authoring Guide](#markdown-authoring-guide) -- Best practices for structuring docs
  - [GitHub Alerts](#github-alerts)
  - [MkDocs Admonitions](#mkdocs-admonitions)
  - [Task Lists](#task-lists)
  - [Collapsible Sections](#collapsible-sections)
  - [Code Block Titles](#code-block-titles)
  - [Status Badges](#status-badges)
  - [Properties Report Macro](#properties-report-macro)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `-p, --profile` | Profile to use (applies to all `orbit` commands) |
| `--service` | Confluence service name, when a profile has multiple Confluence services configured |

---

## page

View a Confluence page. Outputs the page title, space, version number, URL, and body content.

```
orbit confluence page [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the Confluence page |

### Examples

```bash
# View a page by ID
orbit confluence page 123456789 -p myprofile

# View a page using a specific service
orbit confluence page 123456789 -p myprofile --service wiki-prod
```

---

## search

Search Confluence pages using CQL or convenience filters.

```
orbit confluence search [query] [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--space` | string | | Filter by space key. |
| `--title` | string | | Search by title (fuzzy match). |
| `--label` | string | | Filter by label. |
| `--text` | string | | Full-text search. |
| `--cql` | string | | Raw CQL query (overrides other filters). |
| `--limit` | int | 25 | Maximum results. |

**Examples:**

```bash
# Search by space
orbit confluence search --space ISMS --limit 100 -p myprofile

# Search by title
orbit confluence search --space FO --title "Architecture" -p myprofile

# Search by label
orbit confluence search --space FO --label design -p myprofile

# Full-text search
orbit confluence search --space FO --text "deployment pipeline" -p myprofile

# Raw CQL query
orbit confluence search --cql 'space=FO AND label=design AND type=page' -p myprofile
```

---

## children

List the child pages of a given Confluence page.

```
orbit confluence children [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the parent Confluence page |

### Examples

```bash
# List all child pages
orbit confluence children 123456789 -p myprofile
```

---

## hierarchy

Show the full page hierarchy — ancestor chain (from root to the page) and descendant tree.

```
orbit confluence hierarchy [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the Confluence page |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--depth` | 2 | Maximum depth for the children tree |

### Examples

```bash
# Show hierarchy with default depth (2 levels of children)
orbit confluence hierarchy 473676972036 -p myprofile

# Show deeper tree
orbit confluence hierarchy 473676972036 --depth 5 -p myprofile

# JSON tree output
orbit confluence hierarchy 473676972036 -o json -p myprofile
```

### Output

```
Ancestors (top → bottom):
  473267830879 — Foundation Home
    473676972036 — AI Development Process  ← (this page)

Children:
  ├── 473677103107 — Foundations
  │   ├── 473676742659 — Overview & Philosophy
  │   ├── 473677103122 — Compounding Engineering
  ├── 473676873739 — Organization
  │   ├── 473677365249 — AI Layer Setup
```

---

## create

Create a new Confluence page. Pages created via `orbit` are automatically set to wide width.

```
orbit confluence create [flags]
```

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Title for the new page |
| `--space` | Yes | Space key where the page will be created |
| `--parent` | No | Parent page ID to nest the new page under |
| `-b, --body` | No | Page body content as a string (Confluence storage format / XHTML) |
| `-f, --file` | No | Path to a markdown file to convert and upload as the page body |
| `--diagrams` | No | How to render diagram code blocks in `--file`: `kroki` (default), `markdown-macro`, or `code-block` |

### Examples

```bash
# Create a page with inline body content
orbit confluence create --title "Sprint Retrospective" --space DEV --parent 123456789 -p myprofile

# Create a page from a markdown file
orbit confluence create --title "Architecture Overview" --space ENG --parent 123456789 -f ./docs/architecture.md -p myprofile

# Create a top-level page in a space (no parent)
orbit confluence create --title "Project Home" --space PROJ -b "<p>Welcome to the project.</p>" -p myprofile
```

---

## update

Update an existing Confluence page. You can change its title, body, or both.

```
orbit confluence update [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the page to update |

### Flags

| Flag | Description |
|------|-------------|
| `--title` | New title for the page |
| `-b, --body` | New body content as a string |
| `-f, --file` | Path to a markdown file to convert and use as the new body |
| `--diagrams` | How to render diagram code blocks in `--file`: `kroki` (default), `markdown-macro`, or `code-block` |

### Examples

```bash
# Update a page title
orbit confluence update 123456789 --title "Updated Title" -p myprofile

# Update a page body from a markdown file
orbit confluence update 123456789 -f ./docs/revised-design.md -p myprofile

# Update both title and body
orbit confluence update 123456789 --title "Q2 Plan (Final)" -b "<p>Finalized plan.</p>" -p myprofile
```

---

## delete

Delete a Confluence page. The page is moved to the Confluence trash.

```
orbit confluence delete [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the page to delete |

### Examples

```bash
# Delete a page (moves to trash)
orbit confluence delete 123456789 -p myprofile
```

---

## export

Export a Confluence page as markdown or raw storage format.

```
orbit confluence export [page-id] [flags] -p myprofile
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `markdown` | Output format: `markdown` or `storage`. |
| `--output` | string | | Output directory (prints to stdout if omitted). |

Markdown export includes YAML frontmatter with page metadata (`confluence_page_id`, `confluence_space`, `confluence_url`).

**Examples:**

```bash
# Export as markdown to stdout
orbit confluence export 12345 -p myprofile

# Export as markdown to a directory
orbit confluence export 12345 --format markdown --output docs/ -p myprofile

# Export raw storage format (Confluence XHTML)
orbit confluence export 12345 --format storage --output backup/ -p myprofile
```

---

## comment

List, add, and delete comments on a page. Comment bodies are written in Markdown by default (converted to Confluence storage format). Footer (page-level) comments work on both Confluence Server/DC and Cloud; inline comments are Cloud-only and best-effort.

```
orbit confluence comment list <page-id> [flags]
orbit confluence comment add <page-id> [flags]
orbit confluence comment delete <comment-id> [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the page (for `list`/`add`) |
| `comment-id` | The numeric ID of the comment (for `delete`) |

### Flags

**`comment add`:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--body` | string | | Comment body (Markdown by default). Mutually exclusive with `--from-file`. |
| `--from-file` | string | | Read the comment body from a Markdown file. |
| `--storage` | bool | `false` | Treat the body as raw Confluence storage XHTML instead of Markdown. |
| `--reply-to` | string | | Comment ID to reply to (threads the comment). |
| `--inline-text` | string | | Anchor an inline comment to the first occurrence of this page text (Cloud, best-effort — rejected if not found). |

**`comment delete`:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-y`, `--yes` | bool | `false` | Skip the interactive confirmation prompt. |

### Examples

```bash
# List comments (footer + inline, with author and a one-line excerpt)
orbit confluence comment list 12345 -p myprofile
orbit confluence comment list 12345 -o json -p myprofile

# Add a Markdown comment
orbit confluence comment add 12345 --body "LGTM, but the sandbox URL is stale." -p myprofile

# Add a longer review from a file
orbit confluence comment add 12345 --from-file review.md -p myprofile

# Reply to a comment (thread)
orbit confluence comment add 12345 --body "Fixed in v0.60.2" --reply-to 67890 -p myprofile

# Inline comment anchored to page text (Confluence Cloud)
orbit confluence comment add 12345 --body "Inaccurate — deploys are on-demand" --inline-text "deploys run nightly" -p myprofile

# Delete a comment without prompting
orbit confluence comment delete 67890 --yes -p myprofile
```

---

## publish

Publish an entire directory of markdown files to Confluence. This command creates or updates a tree of pages that mirrors the directory structure.

```
orbit confluence publish [directory] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `directory` | Path to the directory containing markdown files |

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--space` | Yes | Confluence space key |
| `--parent` | Yes | Parent page ID under which the pages will be published |
| `--dry-run` | No | Preview what would be published without making any changes |
| `--diagrams` | No | How to render diagram code blocks for the whole run: `kroki` (default), `markdown-macro`, or `code-block` |
| `--no-kroki` | No | Shorthand for `--diagrams code-block` |

### Directory Structure Rules

- **`INDEX.md`** files become the parent page for their directory.
- Other **`.md`** files become child pages under that parent.
- **Subdirectories** are processed recursively, maintaining the hierarchy.
- Files with **`confluence_ignore: true`** in frontmatter are skipped. If the file was previously published (has a `confluence_page_id`), the Confluence page is deleted. When an `INDEX.md` is ignored, the entire subdirectory is skipped.

### Page Title Resolution

Page titles are determined in the following order of precedence:

1. YAML frontmatter `title:` field
2. First `# heading` in the file
3. Filename (without `.md` extension)

### Frontmatter Tracking

After a successful publish, each markdown file has `confluence_page_id` and `confluence_url` added to its YAML frontmatter. These fields are used on subsequent runs to update existing pages rather than creating duplicates.

### Upsert Behavior

When a page has no `confluence_page_id` in frontmatter, orbit searches by title (CQL) first. If create fails because a page with the same title already exists (common with special characters like `&` in titles), orbit falls back to listing the parent's children and matching by title. All pages are set to full-width layout on every create and update.

### Examples

```bash
# Publish a docs directory
orbit confluence publish ./docs --space ENG --parent 123456789 -p myprofile

# Preview what would be published (no changes made)
orbit confluence publish ./docs --space ENG --parent 123456789 --dry-run -p myprofile

# Publish without rendering diagrams at all (leave them as code blocks)
orbit confluence publish ./docs --space ENG --parent 123456789 --diagrams code-block -p myprofile

# Render mermaid with the Confluence markdown macro instead of Kroki
orbit confluence publish ./docs --space ENG --parent 123456789 --diagrams markdown-macro -p myprofile
```

### Diagram Rendering

By default, fenced diagram blocks (` ```mermaid `, ` ```plantuml `, etc.) are rendered as images via Kroki. `--diagrams` picks a different strategy. It is one choice with one value, not a set of switches — a diagram is rendered exactly one way, so there is no combination that can contradict itself.

| Strategy | What it does | What it needs |
|----------|--------------|---------------|
| `kroki` | Renders every diagram type to a PNG served by kroki.io and links the image. The default. | kroki.io reachable at publish time |
| `markdown-macro` | Wraps mermaid diagrams in the Confluence `markdown` macro. An app turns it into a `language-mermaid` code block plus a `mermaid.init` call, so the reader's browser draws the diagram: it stays live, vector and editable in Confluence. **Never contacts kroki.io** — other diagram languages are left as code blocks. | An app providing the `markdown` macro, installed on the Confluence |
| `code-block` | Leaves every diagram as preformatted (ASCII) text. Makes no network call at all. | Nothing |

`markdown-macro` is never the default. The macro comes from an app that is not installed on every Confluence, and a page carrying a macro the instance does not have renders an error where the diagram should be.

#### Choosing a strategy

The strategy can be set in three places. The most specific one wins:

1. **The invocation** — `--diagrams <strategy>` on `publish`, `create` or `update`.
2. **The file** — `confluence_diagrams: <strategy>` in a file's YAML frontmatter.
3. **The instance** — `options.diagrams: <strategy>` on the Confluence service in your profile.

The instance sits last on purpose: whether the markdown macro's app is installed is a property of the Confluence rather than of any one publish, so configuring it there means you say it once and every page published to that instance gets it — while a specific run or a specific file can still override it.

One layer overriding another is not a conflict; it is what the layers are for. What *is* refused is two names for the option **in the same place** disagreeing — `--no-kroki` against `--diagrams`, or `confluence_disable_kroki` against `confluence_diagrams`. Those stop the run with an error naming the file, rather than warning and carrying on with a value you did not ask for. An unrecognised strategy name is refused the same way, wherever it appears: a typo in a setting that decides whether diagram source leaves your network is not something to guess at.

**On the service** (so every publish to this Confluence uses it):

```yaml
profiles:
  - name: myprofile
    services:
      - name: confluence
        type: confluence
        variant: server
        base_url: https://wiki.example.com
        options:
          diagrams: markdown-macro
        auth:
          method: token
          token: op://Private/Confluence/token
```

**On a run:**

```bash
orbit confluence publish ./docs --space ENG --parent 123456789 --diagrams markdown-macro -p myprofile
```

**On a file:**

````markdown
---
title: Architecture Overview
confluence_diagrams: code-block
---

# Architecture

```mermaid
graph TD; A-->B;
```
````

That `mermaid` block publishes as a code block rather than a Kroki image, even on a run that did not ask for it.

`--no-kroki` and `confluence_disable_kroki: true` still work and name the `code-block` strategy. They are shorthand for the same one option, not a second switch beside it: asking for `--no-kroki` and `--diagrams markdown-macro` in the same invocation is an error rather than a silent winner.

#### What the markdown macro renders

Only mermaid, and that was checked against a live instance rather than assumed. A page carrying ` ```mermaid `, ` ```plantuml ` and ` ```graphviz ` fences in three macros came back with all three rendered as `<code class="language-...">` blocks, and the app attached its `mermaid.init` call to `.language-mermaid` alone — plantuml and graphviz render as syntax-highlighted source, not diagrams.

That is the shape you would expect, since the macro renders Markdown and mermaid is a feature of Markdown renderers (as it is on GitHub) rather than a diagram service the macro talks to.

So under `markdown-macro`, mermaid goes to the macro and **every other diagram language is left as a code block** — not wrapped in a macro that would show its source anyway, and not sent to Kroki either.

Leaving them as code rather than falling through to Kroki is deliberate. Rendering them via Kroki would produce more pictures, but it would do it by sending diagram source to a third party on behalf of someone who chose the one strategy that exists to keep diagrams on the instance — and silently, for exactly the languages they are least likely to check. The guarantee is worth more than the extra pictures, and it is one sentence long: **`markdown-macro` never contacts kroki.io.**

If you want mermaid via the macro *and* PlantUML via Kroki, ask for it: set `confluence_diagrams: kroki` on the files that contain PlantUML, or run those separately with `--diagrams kroki`.

The diagram source is passed through untouched. The Mermaid compatibility fixes orbit applies before sending a diagram to Kroki (see the diagram compatibility notes) are not applied here: they work around the Mermaid build behind kroki.io, and the app renders with its own, configured with `htmlLabels: true` — so `<br/>` in a label is a line break rather than something to strip.

This also means a page published this way survives being pulled with `orbit confluence export` and published again unchanged, which was verified end to end against a live Confluence: publish, fetch the stored body, export to markdown, re-convert, and the macro comes back byte for byte.

A diagram source containing `]]>` is safe. It would otherwise end the macro's CDATA section early and corrupt everything after it on the page, so orbit splits it across two sections; Confluence rejoins them on read, and `orbit confluence export` puts the original `]]>` back.

#### When Kroki is slow or down

Diagrams are rendered four at a time, and each one gets 30 seconds covering every retry and the backoff between them. A page of eight diagrams against an unresponsive Kroki therefore costs about two 30s windows rather than eight. Rendering order does not affect the published page: the markup is emitted in document order whatever order the renders finish in.

If kroki.io does not answer within a diagram's budget, the publish stops and names the file and diagram that timed out rather than hanging:

```
Error: converting docs/architecture.md: kroki did not answer within 30s for the
mermaid diagram starting "graph TD;A-->B;": rerun with --diagrams code-block to
publish diagrams as code blocks (or set confluence_diagrams: code-block in this
file's frontmatter), or raise the budget with --timeout
```

The global `--timeout` flag overrides the 30s budget, so `--timeout 2m` allows a slow link and `--timeout 5s` fails faster. `Ctrl-C` stops a publish immediately, including mid-diagram and during a retry backoff; press it again to force quit.

Given this directory:

```
docs/
  INDEX.md
  getting-started.md
  api/
    INDEX.md
    endpoints.md
    authentication.md
```

The resulting Confluence page tree would be:

```
Parent Page (123456789)
  +-- docs (from INDEX.md)
        +-- Getting Started (from getting-started.md)
        +-- api (from api/INDEX.md)
              +-- Endpoints (from endpoints.md)
              +-- Authentication (from authentication.md)
```

---

## set-width

Set the display width of one or more Confluence pages. Accepts multiple page IDs and can optionally apply the change recursively to all child pages.

```
orbit confluence set-width [page-id...] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id...` | One or more page IDs to update |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--width` | `wide` | Width setting: `wide` or `fixed` |
| `--recursive` | `false` | Apply the width setting to all descendant pages |

### Examples

```bash
# Set a single page to wide width
orbit confluence set-width 123456789 -p myprofile

# Set multiple pages to fixed width
orbit confluence set-width 111111111 222222222 --width fixed -p myprofile

# Set a page and all its children to wide width
orbit confluence set-width 123456789 --recursive -p myprofile
```

---

## Notes

### Confluence API Variants

Orbit automatically selects the correct API path based on the configured instance type:

| Variant | API Prefix |
|---------|------------|
| Cloud | `/wiki/rest/api` |
| Server / Data Center | `/rest/api` |

### Markdown Conversion

When a markdown file is provided via `--file` / `-f` (or through `publish`), Orbit converts it to Confluence storage format (XHTML). The converter handles:

- Headings (`#`, `##`, etc.)
- Ordered, unordered, and nested lists
- Task lists (`- [ ]` / `- [x]`)
- Tables with column alignment
- Fenced code blocks with language hints and optional titles
- GitHub Alerts (`> [!NOTE]`, `> [!WARNING]`, `> [!TIP]`, `> [!CAUTION]`, `> [!IMPORTANT]`)
- MkDocs admonitions (`!!! type "Title"`)
- Blockquotes (converted to info panels)
- Collapsible sections (`<details>` / `<summary>`)
- Inline formatting (bold, italic, strikethrough, code, links)
- Status badges (`{status:Color|Text}`)
- Images
- Horizontal rules
- Table of Contents directive (`<!-- confluence:toc-start -->` / `<!-- confluence:toc-end -->`)
- Page Properties macro (from `confluence_properties` frontmatter)
- Properties Report macro (from HTML comments or shorthand syntax)
- Metadata lines

Additional conversion rules:

- The first `# heading` is **skipped** because Confluence already displays the page title.
- YAML frontmatter is **stripped** during conversion.
- Relative `.md` links between published files are resolved to Confluence `<ac:link>` macros during `publish`.

### Confluence HTML Comments

Use HTML comments to control Confluence-specific behavior inside markdown files:

| Comment | Purpose |
|---------|---------|
| `<!-- confluence:ignore-start -->` / `<!-- confluence:ignore-end -->` | Exclude a section from Confluence publishing |
| `<!-- confluence:toc-start -->` / `<!-- confluence:toc-end -->` | Replace the wrapped TOC section with a Confluence TOC macro |
| `<!-- confluence:properties-report cql="..." -->` | Generate a Page Properties Report table |

The TOC directive supports parameters to customize the macro:

```markdown
<!-- confluence:toc-start maxLevel=2 minLevel=2 style="none" -->

## Table of Contents

1. [Section One](#section-one)
2. [Section Two](#section-two)

<!-- confluence:toc-end -->
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `maxLevel` / `maxHeadingLevel` | `3` | Maximum heading depth to include |
| `minLevel` / `minHeadingLevel` | `2` | Minimum heading depth to include |
| `style` | `none` | List style (`none`, `circle`, `disc`, `square`, `decimal`) |
| `outline` | `false` | Show section numbers (set `true` to enable numbered outline) |
| `printable` | `true` | Whether the TOC is included in print view |

These comments are invisible in standard markdown renderers and only affect the Confluence conversion.

---

## Markdown Authoring Guide

Best practices for structuring markdown documents that publish cleanly to Confluence via `orbit confluence publish`.

### Frontmatter

Every markdown file should include YAML frontmatter. Orbit uses it for page title resolution and Confluence sync tracking.

```yaml
---
title: Agent Roles and Responsibilities
subtitle: "Workflow — AI Development Process v2.0"
date: "March 5, 2026"
author: "AI Tooling Guild"
confluence_ignore: false
confluence_page_id: "329383962"
confluence_url: "https://wiki.example.com/display/ENG/Agent+Roles"
confluence_labels:
  - ai-process
  - workflow
confluence_properties:
  - id: status
    fields:
      Owner: "AI Tooling Guild"
      Classification: "Internal"
      Status: "{status:Green|Published}"
      Reviewed on: "2026-03-05"
---
```

| Field | Required | Purpose |
|-------|----------|---------|
| `title` | Yes | Page title in Confluence (takes precedence over `# heading`) |
| `subtitle` | No | Section context and version — useful for readers, stripped during conversion |
| `date` | No | Last updated date — helps readers gauge freshness |
| `author` | No | Accountability — who owns this document |
| `confluence_ignore` | No | When `true`, skip publishing; delete page if previously published (default: `false`) |
| `confluence_page_id` | Auto | Set by `orbit confluence publish` after first publish — used for updates |
| `confluence_url` | Auto | Set by `orbit confluence publish` — direct link to the page |
| `confluence_labels` | No | Labels applied to the Confluence page for search and filtering |
| `confluence_properties` | No | Confluence page properties macro — renders as a metadata table at the top |
| `confluence_diagrams` | No | How to render diagram code blocks in this file: `kroki` (default), `markdown-macro`, or `code-block`. An unrecognised value is an error, not a fallback |
| `confluence_disable_kroki` | No | When `true`, the same as `confluence_diagrams: code-block` (default: `false`). Setting both to different strategies is an error |

> **Tip:** Do not manually set `confluence_page_id` or `confluence_url`. Let `orbit confluence publish` manage them automatically.

### Directory Structure

Organize documents into topic directories with `INDEX.md` files as the parent page for each group.

```
docs/
  INDEX.md                    ← Root parent page
  foundations/
    INDEX.md                  ← "Foundations" parent page
    overview.md
    context-engineering.md
  workflow/
    INDEX.md                  ← "Workflow" parent page
    agent-roles.md
    piv-loops.md
  process/
    INDEX.md
    jira-configuration.md
    sprint-ceremonies.md
  reference/
    INDEX.md
    quick-reference.md
```

Guidelines:

- **One level of nesting** — keep directories flat within each topic group. Avoid deeply nested subdirectories; they produce hard-to-navigate Confluence trees.
- **`INDEX.md` in every directory** — acts as the landing page and links to its child pages.
- **Hyphenated lowercase filenames** — `ai-layer-setup.md`, not `AI Layer Setup.md`. Avoids URL encoding issues.
- **No version numbers in filenames** — track versions in frontmatter and revision history instead.

### INDEX.md Files

Each `INDEX.md` serves as a table of contents for its directory. Structure it as a brief introduction followed by a table linking to child pages.

```markdown
---
title: Workflow
subtitle: "How engineers work with AI agents"
---

# Workflow

How the team uses AI agents in day-to-day development.

| Document | Description |
|----------|-------------|
| [Agent Roles](./agent-roles.md) | The 5 agent roles and their responsibilities |
| [PIV Loops](./piv-loops.md) | Plan → Implement → Validate execution cycle |
| [Parallelization](./parallelization.md) | Running multiple agents concurrently |
```

### Document Structure

Follow this consistent structure within each document:

```markdown
---
(frontmatter)
---

# Document Title

> **One-line summary** of what this document covers and who it's for.

## Table of Contents

1. [Section One](#section-one)
2. [Section Two](#section-two)
3. [Cross-References](#cross-references)
4. [Revision History](#revision-history)

---

## Section One

Content...

---

## Section Two

Content...

---

## Cross-References

| Document | Relevance |
|----------|-----------|
| [Related Doc](../category/related-doc.md) | Why it's relevant |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 2.0 | 2026-03-05 | Restructured into modular format |
| 1.0 | 2026-01-15 | Initial document |
```

Key patterns:

- **Opening blockquote** — one-line summary at the top, sets context immediately.
- **Numbered TOC** — helps readers navigate long documents. Use anchor links.
- **Horizontal rules** (`---`) — separate major sections visually.
- **Cross-References section** — every document links to related docs with a short relevance note. This creates a navigable web across the Confluence space.
- **Revision History** — append-only table at the bottom tracks document evolution.

### Markdown Features That Convert Well

These elements convert cleanly to Confluence storage format:

| Element | Markdown | Notes |
|---------|----------|-------|
| Headings | `##`, `###`, `####` | H1 is skipped (used as page title) |
| Tables | Pipe tables | Supports column alignment (`:---`, `:---:`, `---:`) |
| Code blocks | ` ```bash ` | Language hint enables syntax highlighting |
| Code block titles | ` ```python title="example.py" ` | Adds a title header to the code block |
| Inline code | `` `field_name` `` | Use for field names, commands, paths |
| Bold / Italic | `**bold**` / `*italic*` | Standard inline formatting |
| Strikethrough | `~~text~~` | Renders as struck-through text |
| Links | `[text](url)` | Both relative and absolute URLs work |
| Ordered lists | `1. 2. 3.` | For procedures and sequential steps |
| Unordered lists | `- item` | For guidelines, options, checklists |
| Nested lists | Indented `- item` | Supports multiple levels of nesting |
| Task lists | `- [ ]` / `- [x]` | Renders as status badges (incomplete / complete) |
| Blockquotes | `> text` | Converts to Confluence info panel |
| GitHub Alerts | `> [!NOTE]`, `> [!WARNING]`, etc. | Converts to styled Confluence panels (see below) |
| MkDocs admonitions | `!!! note "Title"` | Indented content block — converts to panel macro |
| Collapsible sections | `<details><summary>` | Converts to Confluence expand macro |
| Status badges | `{status:Green\|Text}` | Renders as colored Confluence status lozenge |
| Horizontal rules | `---`, `***`, `___` | Visual section separators |
| Images | `![alt](path)` | Relative paths resolved during publish |

### Tips for Clean Conversion

1. **Use `##` as your top-level section heading** — since `#` is skipped (it becomes the page title), start your content sections at `##`.

2. **Keep tables simple** — avoid merged cells or multi-line cell content. If a table gets complex, split it into multiple tables or use a list instead.

3. **Prefer fenced code blocks with language hints** — ` ```bash `, ` ```json `, ` ```markdown ` enable syntax highlighting in Confluence.

4. **Use ASCII diagrams for simple flows** — they render in Confluence as preformatted text:
   ```
   Plan → Implement → Validate
     |                    |
     +--- codify <-------+
   ```

   Diagram code blocks (` ```mermaid `, ` ```plantuml `, etc.) are otherwise rendered as images via Kroki. Pass `--diagrams code-block` to `orbit confluence publish` (whole run) or set `confluence_diagrams: code-block` in a file's frontmatter (that file only) to keep them as preformatted code, or `--diagrams markdown-macro` / `confluence_diagrams: markdown-macro` to render mermaid with the Confluence markdown macro.

5. **Cross-link with relative paths** — `[PIV Loops](../workflow/piv-loops.md)` works both locally and in Confluence (resolved during publish).

6. **Avoid HTML in markdown** — Confluence has its own storage format. Raw HTML may not convert correctly. Exception: `<!-- confluence: -->` comments for publish control.

7. **One idea per file** — if a document exceeds ~500 lines, consider splitting it into a subdirectory with its own `INDEX.md`.

### GitHub Alerts

GitHub-style alert blockquotes are converted to styled Confluence panels with matching colors and icons.

```markdown
> [!NOTE]
> Useful information that users should know.

> [!TIP]
> Helpful advice for doing things better.

> [!IMPORTANT]
> Key information users need to know.

> [!WARNING]
> Urgent info that needs immediate attention.

> [!CAUTION]
> Advises about risks or negative outcomes.
```

| Alert Type | Confluence Panel | Color |
|------------|-----------------|-------|
| `[!NOTE]` | Info panel | Blue |
| `[!TIP]` | Tip panel | Green |
| `[!IMPORTANT]` | Note panel | Yellow |
| `[!WARNING]` | Warning panel | Yellow |
| `[!CAUTION]` | Warning panel | Red |

### MkDocs Admonitions

MkDocs-style admonitions (`!!! type "Title"`) are also supported. Content must be indented with 4 spaces.

```markdown
!!! note "Important Note"
    This is the admonition content.
    It can span multiple lines.

!!! warning "Deprecation Warning"
    This API will be removed in v3.0.
```

### Task Lists

Task lists render as status badges in Confluence:

```markdown
- [x] Completed task
- [ ] Pending task
- [x] Another completed task
```

Checked items render as a green `DONE` status badge, unchecked items render as a grey `TO DO` badge.

### Collapsible Sections

HTML `<details>` / `<summary>` blocks convert to Confluence expand macros:

```markdown
<details>
<summary>Click to expand</summary>

Hidden content goes here. Supports full markdown inside.

</details>
```

### Code Block Titles

Add a `title` attribute to fenced code blocks for a descriptive header:

````markdown
```python title="example.py"
def hello():
    print("Hello, world!")
```
````

### Status Badges

Inline status lozenges using the `{status:Color|Text}` syntax:

```markdown
Current status: {status:Green|Approved}
Risk level: {status:Red|High}
Phase: {status:Blue|In Progress}
```

This is particularly useful inside `confluence_properties` fields and table cells.

### Properties Report Macro

Aggregate page properties across child pages using either syntax:

**HTML comment syntax:**

```markdown
<!-- confluence:properties-report cql="label = 'team-status'" firstcolumn="Owner" headings="Status,Classification" -->
```

**Shorthand syntax:**

```markdown
{properties-report: label="team-status", columns="Owner,Status,Classification", sortBy="Owner"}
```

Both generate a Confluence `detailssummary` macro that pulls data from Page Properties macros on child pages matching the query.

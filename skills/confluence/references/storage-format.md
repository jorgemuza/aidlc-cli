# Confluence Storage Format Reference

Confluence stores page content in an XHTML-based format called "storage format". When using `--body` directly (without `--file`), you need to provide content in this format. When using `--file` with a `.md` file, the orbit CLI converts markdown to storage format automatically.

## Basic Elements

### Paragraphs

```xml
<p>This is a paragraph.</p>
<p>This is another paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
```

### Headings

```xml
<h1>Heading 1</h1>
<h2>Heading 2</h2>
<h3>Heading 3</h3>
<h4>Heading 4</h4>
<h5>Heading 5</h5>
<h6>Heading 6</h6>
```

### Lists

```xml
<ul>
  <li>Unordered item 1</li>
  <li>Unordered item 2</li>
</ul>

<ol>
  <li>Ordered item 1</li>
  <li>Ordered item 2</li>
</ol>
```

### Tables

```xml
<table>
  <thead>
    <tr>
      <th>Header 1</th>
      <th>Header 2</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Cell 1</td>
      <td>Cell 2</td>
    </tr>
  </tbody>
</table>
```

### Links

```xml
<a href="https://example.com">Link text</a>
```

### Images

```xml
<ac:image>
  <ri:url ri:value="https://example.com/image.png" />
</ac:image>
```

### Inline Formatting

```xml
<strong>Bold</strong>
<em>Italic</em>
<code>Inline code</code>
<del>Strikethrough</del>
```

## Confluence Macros

Macros use the `ac:structured-macro` element with Atlassian-specific namespaces.

### Code Block

```xml
<ac:structured-macro ac:name="code">
  <ac:parameter ac:name="language">python</ac:parameter>
  <ac:plain-text-body><![CDATA[def hello():
    print("Hello World")]]></ac:plain-text-body>
</ac:structured-macro>
```

### Markdown (Mermaid Diagrams)

Emitted only under the `markdown-macro` diagram strategy (`--diagrams markdown-macro`,
`confluence_diagrams: markdown-macro`, or `options.diagrams` on the service). Requires an app
providing the `markdown` macro to be installed on the Confluence.

```xml
<ac:structured-macro ac:name="markdown">
  <ac:plain-text-body><![CDATA[```mermaid
graph TD;A-->B;
```]]></ac:plain-text-body>
</ac:structured-macro>
```

The fenced block is *content inside the CDATA*, not a child element — the same plain-text-body
shape as the code macro above. Notes:

- `ac:macro-id` is **not** emitted. Confluence assigns one per macro on save; sending our own
  would duplicate the server's, or repeat a single id across a whole run.
- Only mermaid is wrapped this way. Other diagram languages render inside the macro as
  syntax-highlighted source rather than diagrams, so under this strategy they are left as code
  blocks — never sent to Kroki.
- A source containing `]]>` is split across two CDATA sections; a source containing a fence gets
  a longer fence. Both are reversed on read.

### Info Panel (Blockquote)

```xml
<ac:structured-macro ac:name="info">
  <ac:rich-text-body>
    <p>This is an informational callout.</p>
  </ac:rich-text-body>
</ac:structured-macro>
```

Other panel types: `note`, `warning`, `tip`

### Table of Contents

```xml
<ac:structured-macro ac:name="toc">
  <ac:parameter ac:name="printable">true</ac:parameter>
  <ac:parameter ac:name="style">disc</ac:parameter>
  <ac:parameter ac:name="maxLevel">3</ac:parameter>
  <ac:parameter ac:name="minLevel">2</ac:parameter>
  <ac:parameter ac:name="outline">false</ac:parameter>
</ac:structured-macro>
```

### Expand (Collapse)

```xml
<ac:structured-macro ac:name="expand">
  <ac:parameter ac:name="title">Click to expand</ac:parameter>
  <ac:rich-text-body>
    <p>Hidden content here.</p>
  </ac:rich-text-body>
</ac:structured-macro>
```

### Status Badge

```xml
<ac:structured-macro ac:name="status">
  <ac:parameter ac:name="title">IN PROGRESS</ac:parameter>
  <ac:parameter ac:name="colour">Yellow</ac:parameter>
</ac:structured-macro>
```

Colors: Grey, Red, Yellow, Green, Blue, Purple

**Emoji shortcode syntax** — GitHub-style emoji circles are automatically converted to status macros:

| Markdown | Rendered as |
|----------|------------|
| `:green_circle: Green` | `{status:Green\|Green}` → Green status badge |
| `:yellow_circle: Yellow` | `{status:Yellow\|Yellow}` → Yellow status badge |
| `:red_circle: Red` | `{status:Red\|Red}` → Red status badge |
| `:blue_circle: Blue` | `{status:Blue\|Blue}` → Blue status badge |

This allows status tables in markdown to render natively in Confluence without rewriting.

### Date Macro

```xml
<time datetime="2026-03-06" />
```

Renders as a formatted date in the Confluence UI. Generated automatically from `YYYY-MM-DD` values in `confluence_properties` fields.

### Page Properties (details macro)

Wraps a two-column table of metadata. Used by the Page Properties Report macro to aggregate data across pages.

```xml
<ac:structured-macro ac:name="details" ac:macro-id="status">
  <ac:rich-text-body>
    <table>
      <tbody>
        <tr>
          <td>Owner</td>
          <td>AI Tooling Guild</td>
        </tr>
        <tr>
          <td>Status</td>
          <td>
            <ac:structured-macro ac:name="status">
              <ac:parameter ac:name="title">approved</ac:parameter>
              <ac:parameter ac:name="colour">Green</ac:parameter>
            </ac:structured-macro>
          </td>
        </tr>
        <tr>
          <td>Reviewed on</td>
          <td><time datetime="2026-03-06" /></td>
        </tr>
      </tbody>
    </table>
  </ac:rich-text-body>
</ac:structured-macro>
```

- `ac:macro-id` is set from the `id` sub-key of `confluence_properties` frontmatter
- Field values matching `{status:Color|Text}` become nested status badge macros
- Field values matching `YYYY-MM-DD` become `<time>` date macros
- Generated by `BuildPagePropertiesMacro()` and prepended to page content during `publish`

### Page Properties Report (detailssummary macro)

Generates a dynamic table that aggregates Page Properties from pages matching a CQL query.

```xml
<ac:structured-macro ac:name="detailssummary">
  <ac:parameter ac:name="cql">label = "ai-process" and space = currentSpace()</ac:parameter>
  <ac:parameter ac:name="firstcolumn">Document</ac:parameter>
  <ac:parameter ac:name="headings">Status, Classification, Reviewed on</ac:parameter>
  <ac:parameter ac:name="sortBy">Title</ac:parameter>
</ac:structured-macro>
```

Parameters:
- `cql` -- CQL query to select pages
- `firstcolumn` -- Name of the first column (defaults to "Title")
- `headings` -- Comma-separated list of property field names to display as columns
- `sortBy` -- Optional column name to sort by

Generated from markdown directives:
- HTML comment: `<!-- confluence:properties-report cql="..." firstcolumn="..." headings="..." -->`
- Shorthand: `{properties-report: label="...", columns="..."}`

When `label` is used instead of `cql`, it expands to `label = "value" and space = currentSpace()`.

## Horizontal Rule

```xml
<hr />
```

## What the Markdown Converter Does

When you use `--file` with a markdown file, the orbit CLI:

1. **Strips YAML frontmatter** — The `---` delimited block at the top is removed
2. **Skips first h1** — The first `# Heading` is dropped since Confluence displays the page title
3. **Replaces TOC sections** — `## Table of Contents` (or "Contents"/"TOC") headings replace the entire section with the Confluence `toc` macro
4. **Converts headings** — `##` through `######` become `<h2>` through `<h6>`
5. **Converts nested lists** — `- ` and `* ` become `<ul>`, `1. ` becomes `<ol>`, with proper nesting via indentation
6. **Task lists** — `- [ ]` and `- [x]` items become status macros (TODO/DONE)
7. **Converts tables** — Markdown tables with `|` pipes become `<table>` with `<thead>`, `<tbody>`, and column alignment support (`:---`, `:---:`, `---:`)
8. **Converts code blocks** — Fenced code blocks with language hints become code macros with syntax highlighting. Supports `title="..."` parameter. **Diagram blocks** (`mermaid`, `plantuml`, `graphviz`, `d2`, `ditaa`, etc.) are rendered as PNG images via [kroki.io](https://kroki.io) — max 800px height, clickable to full image
9. **GitHub Alerts** — `> [!NOTE]`, `> [!WARNING]`, `> [!TIP]`, `> [!CAUTION]`, `> [!IMPORTANT]` map to Confluence note/warning/tip/info macros
10. **Plain blockquotes** — `> ` lines become info panel macros
11. **MkDocs admonitions** — `!!! note "title"` with indented content maps to Confluence panel macros
12. **Collapsible sections** — `<details><summary>Title</summary>` blocks become expand macros
13. **Status badges** — Inline `{status:Color|Text}` syntax becomes status macros. Emoji shortcodes (`:green_circle: Green`, `:yellow_circle: Yellow`, `:red_circle: Red`) are also converted automatically
14. **Horizontal rules** — `---`, `***`, `___` become `<hr />`
15. **Image alt text** — Preserved as `ac:alt` attribute on `ac:image` elements
16. **Converts inline** — `**bold**`, `` `code` ``, `~~strike~~`, links, images
17. **Strips relative links** — Links to `.md` files or `./` paths are converted to plain text (they'll be separate Confluence pages)
18. **Page Properties macro** — `confluence_properties` in frontmatter generates a `details` macro (two-column table) prepended to the page body. Field values support `{status:Color|Text}` (status badge) and `YYYY-MM-DD` (date macro)
19. **Page Properties Report** — `<!-- confluence:properties-report ... -->` comments and `{properties-report: ...}` shorthand directives become `detailssummary` macros
20. **Labels** — `confluence_labels` in frontmatter are applied to the page via the REST API after create/update (not part of the storage format body, but applied during `publish`)
21. **Ignore blocks** — `<!-- confluence:ignore-start -->` and `<!-- confluence:ignore-end -->` markers cause all content between them to be skipped during conversion. Useful for static markdown tables on INDEX pages that should only render in markdown viewers while Confluence uses dynamic Properties Report macros. Ignore markers inside fenced code blocks are treated as literal text

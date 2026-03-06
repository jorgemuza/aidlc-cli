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

Colors: Grey, Red, Yellow, Green, Blue

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
5. **Converts lists** — `- ` and `* ` become `<ul>`, `1. ` becomes `<ol>`
6. **Converts tables** — Markdown tables with `|` pipes become `<table>` with proper `<thead>` and `<tbody>`
7. **Converts code blocks** — Fenced code blocks with language hints become code macros with syntax highlighting
8. **Converts blockquotes** — `> ` lines become info panel macros
9. **Converts inline** — `**bold**`, `` `code` ``, `~~strike~~`, links, images
10. **Strips relative links** — Links to `.md` files or `./` paths are converted to plain text (they'll be separate Confluence pages)

# Diagram Compatibility Reference

Orbit renders fenced diagram code blocks as PNG images via [kroki.io](https://kroki.io). Mermaid diagrams are auto-sanitized before rendering, but some issues require manual fixes.

Everything on this page describes the Kroki path — the `kroki` diagram strategy, which is the default. None of it applies under the other two:

- `--diagrams code-block` (whole run) or `confluence_diagrams: code-block` (one file) skips rendering entirely and keeps diagram blocks as preformatted (ASCII) text. `--no-kroki` and `confluence_disable_kroki: true` are older spellings of the same thing.
- `--diagrams markdown-macro` renders mermaid with the Confluence `markdown` macro instead, which needs an app on the Confluence and passes the diagram source through untouched — none of the auto-fixes below run, because they work around the Mermaid build behind kroki.io and the app renders with its own, configured with `htmlLabels: true` (so `<br/>` in a label is a line break, not something to strip). This strategy never contacts kroki.io at all: other diagram languages are left as code blocks, so nothing on this page applies to it.

## Auto-Fixed Issues (Mermaid)

These are automatically corrected by orbit during publishing:

| Issue | Example | Auto-fix |
|-------|---------|----------|
| `<br/>` in participant names | `participant API as Django API<br/>:8000` | Stripped, replaced with `. ` |
| Parenthesized suffixes | `participant S3 as S3 (LocalStack)` | `S3 - LocalStack` |
| Trailing `()` | `participant D as dispatch_event()` | `dispatch_event` |
| Port numbers after colon | `participant API as Django API:8000` | `Django API` |
| `<br/>` in notes | `Note right of X: line1<br/>line2` | `line1. line2` |
| `<br/>` in messages | `A->>B: first<br/>second` | `first. second` |
| Reverse arrows `<<--` | `Client<<--API: event` | `API-->>Client: event` (sender/receiver flipped) |

## Manual Fixes Required

These cannot be auto-fixed and need to be corrected in the source markdown:

### ASCII art is not a diagram

```
BAD:  ```
      ┌──────────┐
      │  Model   │
      └──────────┘
      ```
```

Use a proper diagram language instead:

```
GOOD: ```mermaid
      erDiagram
          Model {
              UUID id PK
              string name
          }
      ```
```

### ASCII sequence diagrams need conversion

```
BAD:  ```
      Developer    API    ActionRequest
          |-- DELETE ->|        |
          |<-- 202 ---|        |
      ```

GOOD: ```mermaid
      sequenceDiagram
          Developer->>API: DELETE /specs/id
          API-->>Developer: 202
      ```
```

### Large diagrams may exceed Kroki limits

Very complex diagrams with 15+ participants or 50+ messages may fail. Split into smaller diagrams.

## Diagram Type Selection

| Content | Best format | Why |
|---------|------------|-----|
| Request/response flows | `mermaid` sequenceDiagram | Native sequence support |
| State machines | `mermaid` stateDiagram-v2 | Clean state transitions |
| Simple flowcharts | `mermaid` flowchart TD/LR | Widely supported |
| Component trees | `mermaid` graph TD | Parent-child layout |
| Entity/data models | `mermaid` erDiagram | Field definitions + relationships |
| Nested architecture | `d2` | Superior nesting support |
| Layered systems | `d2` | Named containers with connections |

## Supported Diagram Languages

`mermaid`, `plantuml`, `graphviz`, `dot`, `d2`, `ditaa`, `erd`, `nomnoml`, `svgbob`, `vega`, `vegalite`, `wavedrom`, `pikchr`, `structurizr`, `excalidraw`, `c4plantuml`

## Rendering Details

- Output format: PNG (SVG has xlink namespace issues in Confluence)
- Smart sizing: max 600px wide, max 800px tall — auto-scaled preserving aspect ratio
- Clickable: links to full-resolution PNG in new tab
- Fallback: invalid diagrams render as syntax-highlighted code blocks
- Encoding: zlib + base64url (per Kroki spec)

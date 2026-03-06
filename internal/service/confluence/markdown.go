package confluence

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// MarkdownToStorage converts Markdown to Confluence storage format (XHTML).
// The first h1 heading is skipped because Confluence already displays the page title.
// Any "Table of Contents" section is replaced with the Confluence toc macro.
func MarkdownToStorage(md string) string {
	return MarkdownToStorageWithLinks(md, nil)
}

// MarkdownToStorageWithLinks converts Markdown to Confluence storage format.
// linkMap maps relative .md paths (e.g. "./piv-loops.md", "../org/governance.md")
// to exact Confluence page titles for resolving internal links.
func MarkdownToStorageWithLinks(md string, linkMap map[string]string) string {
	// Strip YAML frontmatter
	md = stripFrontmatter(md)

	// Convert document metadata lines (**Key**: Value) into a two-column table
	md = convertMetadataBlock(md)

	lines := strings.Split(md, "\n")
	var out []string
	inCodeBlock := false
	codeLang := ""
	var codeLines []string
	inList := false
	listTag := ""
	inTable := false
	firstH1Skipped := false
	inTocSection := false

	inline := func(text string) string {
		return convertInline(text, linkMap)
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// Close code block
				code := escapeXML(strings.Join(codeLines, "\n"))
				if codeLang != "" {
					out = append(out, fmt.Sprintf(`<div class="code-block" data-layout="full-width"><ac:structured-macro ac:name="code"><ac:parameter ac:name="language">%s</ac:parameter><ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body></ac:structured-macro></div>`, codeLang, code))
				} else {
					out = append(out, fmt.Sprintf(`<div class="code-block" data-layout="full-width"><ac:structured-macro ac:name="code"><ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body></ac:structured-macro></div>`, code))
				}
				inCodeBlock = false
				codeLang = ""
				codeLines = nil
			} else {
				if inList {
					out = append(out, closeList(listTag))
					inList = false
				}
				inCodeBlock = true
				codeLang = strings.TrimPrefix(line, "```")
				codeLang = strings.TrimSpace(codeLang)
				codeLines = nil
			}
			continue
		}
		if inCodeBlock {
			codeLines = append(codeLines, line)
			continue
		}

		trimmed := strings.TrimSpace(line)

		// Blank line — close list/table
		if trimmed == "" {
			if inList {
				out = append(out, closeList(listTag))
				inList = false
			}
			if inTable {
				out = append(out, "</tbody></table>")
				inTable = false
			}
			continue
		}

		// Horizontal rule — skip entirely, not useful in Confluence
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			continue
		}

		// Headings
		if heading, level := parseHeading(trimmed); heading != "" {
			if inList {
				out = append(out, closeList(listTag))
				inList = false
			}
			if inTable {
				out = append(out, "</tbody></table>")
				inTable = false
			}

			// Skip the first h1 — Confluence already shows the page title
			if level == 1 && !firstH1Skipped {
				firstH1Skipped = true
				continue
			}

			// Replace "Table of Contents" heading with Confluence toc macro
			// and skip all lines until the next heading or hr
			headingLower := strings.ToLower(heading)
			if headingLower == "table of contents" || headingLower == "contents" || headingLower == "toc" {
				out = append(out, `<ac:structured-macro ac:name="toc"><ac:parameter ac:name="printable">true</ac:parameter><ac:parameter ac:name="style">none</ac:parameter><ac:parameter ac:name="maxLevel">3</ac:parameter><ac:parameter ac:name="minLevel">2</ac:parameter></ac:structured-macro>`)
				inTocSection = true
				continue
			}

			// If we hit a new heading, we're out of the TOC section
			inTocSection = false

			out = append(out, fmt.Sprintf("<h%d>%s</h%d>", level, inline(heading), level))
			continue
		}

		// Skip lines inside a TOC section (the manual list of links)
		if inTocSection {
			// Stop skipping at horizontal rules or blank-then-heading
			if trimmed == "---" || trimmed == "***" || trimmed == "___" {
				inTocSection = false
				// Don't emit the hr either, the toc macro is enough
			}
			continue
		}

		// Tables
		if isTableRow(trimmed) {
			// Skip separator rows
			if isTableSeparator(trimmed) {
				continue
			}
			if inList {
				out = append(out, closeList(listTag))
				inList = false
			}
			cells := parseTableRow(trimmed)
			if !inTable {
				// First row is header
				out = append(out, `<table data-layout="full-width"><thead><tr>`)
				for _, cell := range cells {
					out = append(out, fmt.Sprintf("<th><strong>%s</strong></th>", inline(cell)))
				}
				out = append(out, "</tr></thead><tbody>")
				inTable = true
			} else {
				out = append(out, "<tr>")
				for _, cell := range cells {
					out = append(out, fmt.Sprintf("<td>%s</td>", inline(cell)))
				}
				out = append(out, "</tr>")
			}
			continue
		}

		if inTable {
			out = append(out, "</tbody></table>")
			inTable = false
		}

		// Blockquotes
		if strings.HasPrefix(trimmed, "> ") {
			if inList {
				out = append(out, closeList(listTag))
				inList = false
			}
			// Collect consecutive blockquote lines
			var quoteLines []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "> ") {
				quoteLines = append(quoteLines, strings.TrimPrefix(strings.TrimSpace(lines[i]), "> "))
				i++
			}
			i-- // back up since for loop will increment
			quoteContent := inline(strings.Join(quoteLines, "<br />"))
			out = append(out, fmt.Sprintf(`<ac:structured-macro ac:name="info"><ac:rich-text-body><p>%s</p></ac:rich-text-body></ac:structured-macro>`, quoteContent))
			continue
		}

		// Unordered list items
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := trimmed[2:]
			if !inList || listTag != "ul" {
				if inList {
					out = append(out, closeList(listTag))
				}
				out = append(out, "<ul>")
				inList = true
				listTag = "ul"
			}
			out = append(out, fmt.Sprintf("<li>%s</li>", inline(content)))
			continue
		}

		// Ordered list items
		if matched, _ := regexp.MatchString(`^\d+\.\s`, trimmed); matched {
			re := regexp.MustCompile(`^\d+\.\s(.*)`)
			m := re.FindStringSubmatch(trimmed)
			if m != nil {
				if !inList || listTag != "ol" {
					if inList {
						out = append(out, closeList(listTag))
					}
					out = append(out, "<ol>")
					inList = true
					listTag = "ol"
				}
				out = append(out, fmt.Sprintf("<li>%s</li>", inline(m[1])))
				continue
			}
		}

		// Close list if we hit a non-list line
		if inList {
			out = append(out, closeList(listTag))
			inList = false
		}

		// Pass through lines that are already HTML block elements (e.g. from convertMetadataBlock)
		if strings.HasPrefix(trimmed, "<table") || strings.HasPrefix(trimmed, "<div") || strings.HasPrefix(trimmed, "<ac:") {
			out = append(out, trimmed)
			continue
		}

		// Regular paragraph
		out = append(out, fmt.Sprintf("<p>%s</p>", inline(trimmed)))
	}

	// Close any open structures
	if inList {
		out = append(out, closeList(listTag))
	}
	if inTable {
		out = append(out, "</tbody></table>")
	}

	return strings.Join(out, "\n")
}

func closeList(tag string) string {
	return fmt.Sprintf("</%s>", tag)
}

func stripFrontmatter(md string) string {
	if !strings.HasPrefix(md, "---") {
		return md
	}
	// Find closing ---
	rest := md[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return md
	}
	return strings.TrimSpace(rest[idx+4:])
}

func parseHeading(line string) (string, int) {
	for level := 6; level >= 1; level-- {
		prefix := strings.Repeat("#", level) + " "
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(line[len(prefix):]), level
		}
	}
	return "", 0
}

func isTableRow(line string) bool {
	return strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|")
}

func isTableSeparator(line string) bool {
	cleaned := strings.ReplaceAll(line, "|", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ":", "")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned == ""
}

func parseTableRow(line string) []string {
	// Remove leading and trailing |
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	cells := make([]string, len(parts))
	for i, p := range parts {
		cells[i] = strings.TrimSpace(p)
	}
	return cells
}

var (
	reBold      = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalic    = regexp.MustCompile(`(?:^|[^*])_(.+?)_(?:[^*]|$)`)
	reItalicAlt = regexp.MustCompile(`\*([^*]+?)\*`)
	reCode      = regexp.MustCompile("`([^`]+)`")
	reStrike    = regexp.MustCompile(`~~(.+?)~~`)
	reLink      = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reImage     = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	reMetadata  = regexp.MustCompile(`^\*\*(.+?)\*\*:\s*(.+)$`)
)

// convertInline converts markdown inline formatting to Confluence storage format.
// linkMap is optional — when provided, relative .md links are resolved to exact
// Confluence page titles. When nil, the link label is used as the page title.
func convertInline(text string, linkMap map[string]string) string {
	// Images before links (images start with !)
	text = reImage.ReplaceAllString(text, `<ac:image><ri:url ri:value="$2" /></ac:image>`)
	// Links — handle relative .md links as Confluence page links
	text = reLink.ReplaceAllStringFunc(text, func(s string) string {
		m := reLink.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		label, href := m[1], m[2]
		// Anchor-only links (e.g. #section) — skip, not useful in storage format
		if strings.HasPrefix(href, "#") {
			return label
		}
		// Relative markdown links become Confluence page links
		if strings.HasSuffix(href, ".md") || strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../") {
			pageTitle := label // default: use the link label as page title
			if linkMap != nil {
				// Try to resolve from the link map
				// Normalize the href: strip ./ prefix, keep ../ paths
				normalized := href
				normalized = strings.TrimPrefix(normalized, "./")
				// Also try with the raw href
				if t, ok := linkMap[normalized]; ok {
					pageTitle = t
				} else if t, ok := linkMap[href]; ok {
					pageTitle = t
				} else {
					// Try just the filename
					base := filepath.Base(href)
					if t, ok := linkMap[base]; ok {
						pageTitle = t
					}
				}
			}
			return fmt.Sprintf(`<ac:link><ri:page ri:content-title="%s" /><ac:plain-text-link-body><![CDATA[%s]]></ac:plain-text-link-body></ac:link>`, pageTitle, label)
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, href, label)
	})
	// Bold before italic
	text = reBold.ReplaceAllString(text, "<strong>$1</strong>")
	// Code
	text = reCode.ReplaceAllString(text, "<code>$1</code>")
	// Strikethrough
	text = reStrike.ReplaceAllString(text, "<del>$1</del>")
	// Italic (using * — careful not to conflict with bold already converted)
	// Only convert remaining single * pairs that aren't inside <strong>
	text = reItalicAlt.ReplaceAllStringFunc(text, func(s string) string {
		if strings.Contains(s, "<strong>") || strings.Contains(s, "</strong>") {
			return s
		}
		m := reItalicAlt.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		return "<em>" + m[1] + "</em>"
	})
	return text
}

// convertMetadataBlock detects consecutive **Key**: Value lines near the top
// of the document (after frontmatter stripping, skipping the first heading and
// blank lines) and replaces them with a pre-rendered Confluence two-column
// table with the first column highlighted in gray.
func convertMetadataBlock(md string) string {
	lines := strings.Split(md, "\n")
	var metaRows [][]string
	metaStart := -1
	metaEnd := 0

	// Skip leading headings, blank lines, and blockquotes to find metadata block
	i := 0
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "> ") {
			i++
			continue
		}
		break
	}

	// Now look for consecutive **Key**: Value lines
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			if len(metaRows) > 0 {
				i++
				continue
			}
			break
		}
		m := reMetadata.FindStringSubmatch(trimmed)
		if m == nil {
			break
		}
		if metaStart == -1 {
			metaStart = i
		}
		metaRows = append(metaRows, []string{m[1], m[2]})
		metaEnd = i + 1
		i++
	}

	if len(metaRows) == 0 {
		return md
	}

	// Build the table with gray first column using Confluence's native highlight
	var tbl strings.Builder
	tbl.WriteString(`<table data-layout="full-width"><tbody>`)
	for _, row := range metaRows {
		tbl.WriteString(fmt.Sprintf(`<tr><td class="highlight-grey" data-highlight-colour="grey"><strong>%s</strong></td><td>%s</td></tr>`, row[0], row[1]))
	}
	tbl.WriteString(`</tbody></table>`)

	// Reconstruct: lines before metadata + table + lines after metadata
	before := lines[:metaStart]
	after := lines[metaEnd:]
	// Skip trailing blank lines after metadata
	for len(after) > 0 && strings.TrimSpace(after[0]) == "" {
		after = after[1:]
	}

	var result []string
	result = append(result, before...)
	result = append(result, tbl.String())
	result = append(result, after...)
	return strings.Join(result, "\n")
}

func escapeXML(s string) string {
	// Inside CDATA we don't need XML escaping, but we need to handle ]]>
	return strings.ReplaceAll(s, "]]>", "]]]]><![CDATA[>")
}

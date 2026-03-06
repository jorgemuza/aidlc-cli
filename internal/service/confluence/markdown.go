package confluence

import (
	"fmt"
	"regexp"
	"strings"
)

// MarkdownToStorage converts Markdown to Confluence storage format (XHTML).
// The first h1 heading is skipped because Confluence already displays the page title.
// Any "Table of Contents" section is replaced with the Confluence toc macro.
func MarkdownToStorage(md string) string {
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

			out = append(out, fmt.Sprintf("<h%d>%s</h%d>", level, convertInline(heading), level))
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
					out = append(out, fmt.Sprintf("<th>%s</th>", convertInline(cell)))
				}
				out = append(out, "</tr></thead><tbody>")
				inTable = true
			} else {
				out = append(out, "<tr>")
				for _, cell := range cells {
					out = append(out, fmt.Sprintf("<td>%s</td>", convertInline(cell)))
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
			quoteContent := convertInline(strings.Join(quoteLines, "<br />"))
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
			out = append(out, fmt.Sprintf("<li>%s</li>", convertInline(content)))
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
				out = append(out, fmt.Sprintf("<li>%s</li>", convertInline(m[1])))
				continue
			}
		}

		// Close list if we hit a non-list line
		if inList {
			out = append(out, closeList(listTag))
			inList = false
		}

		// Regular paragraph
		out = append(out, fmt.Sprintf("<p>%s</p>", convertInline(trimmed)))
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
	reBold       = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalic     = regexp.MustCompile(`(?:^|[^*])_(.+?)_(?:[^*]|$)`)
	reItalicAlt  = regexp.MustCompile(`\*([^*]+?)\*`)
	reCode       = regexp.MustCompile("`([^`]+)`")
	reStrike     = regexp.MustCompile(`~~(.+?)~~`)
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reImage      = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	reMetadata   = regexp.MustCompile(`^\*\*(.+?)\*\*:\s*(.+)$`)
)

func convertInline(text string) string {
	// Images before links (images start with !)
	text = reImage.ReplaceAllString(text, `<ac:image><ri:url ri:value="$2" /></ac:image>`)
	// Links — handle relative .md links
	text = reLink.ReplaceAllStringFunc(text, func(s string) string {
		m := reLink.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		label, href := m[1], m[2]
		// Skip relative markdown links (they'll be Confluence pages)
		if strings.HasSuffix(href, ".md") || strings.HasPrefix(href, "./") {
			return label
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

// convertMetadataBlock detects consecutive **Key**: Value lines at the start
// of the document (before any heading) and replaces them with a pre-rendered
// Confluence two-column table with the first column in gray.
func convertMetadataBlock(md string) string {
	lines := strings.Split(md, "\n")
	var metaRows [][]string
	consumed := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(metaRows) > 0 {
				consumed++
			}
			continue
		}
		m := reMetadata.FindStringSubmatch(trimmed)
		if m == nil {
			break
		}
		metaRows = append(metaRows, []string{m[1], m[2]})
		consumed++
	}

	if len(metaRows) == 0 {
		return md
	}

	// Build the table with gray first column
	var tbl strings.Builder
	tbl.WriteString(`<table data-layout="full-width"><tbody>`)
	for _, row := range metaRows {
		tbl.WriteString(fmt.Sprintf(`<tr><td style="background-color: #f4f5f7;"><strong>%s</strong></td><td>%s</td></tr>`, row[0], row[1]))
	}
	tbl.WriteString(`</tbody></table>`)

	// Skip consumed lines and any trailing blank lines
	rest := lines[consumed:]
	for len(rest) > 0 && strings.TrimSpace(rest[0]) == "" {
		rest = rest[1:]
	}

	return tbl.String() + "\n" + strings.Join(rest, "\n")
}

func escapeXML(s string) string {
	// Inside CDATA we don't need XML escaping, but we need to handle ]]>
	return strings.ReplaceAll(s, "]]>", "]]]]><![CDATA[>")
}

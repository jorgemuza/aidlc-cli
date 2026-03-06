package confluence

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	conflsvc "github.com/jorgemuza/orbit/internal/service/confluence"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish [directory]",
	Short: "Publish a directory of markdown files to Confluence",
	Long: `Publish recursively converts markdown files to Confluence storage format
and creates or updates pages preserving the directory structure. INDEX.md files
become the parent page for their directory. Other .md files become child pages.

If a markdown file has confluence_page_id in its YAML frontmatter, the existing
page is updated. Otherwise a new page is created. After creation, the page ID
and URL are written back to the frontmatter for future updates.

Relative markdown links (e.g. [text](./file.md)) are converted to Confluence
page links using the target file's title.`,
	Args: cobra.ExactArgs(1),
	Example: `  orbit confluence publish ./docs --space FO --parent 12345
  orbit confluence publish ./docs --space FO --parent 12345 --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		space, _ := cmd.Flags().GetString("space")
		parentID, _ := cmd.Flags().GetString("parent")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			fmt.Println("DRY RUN — no pages will be created or updated")
			fmt.Println()
		}

		var client *conflsvc.Client
		if !dryRun {
			var err error
			client, err = resolveConfluenceClient(cmd)
			if err != nil {
				return err
			}
		}

		baseURL := ""
		if client != nil {
			baseURL = client.Conn.BaseURL
		}

		// Build link map: relative path → page title for all .md files
		linkMap := buildLinkMap(dir)

		return publishDir(client, dir, space, parentID, baseURL, linkMap, dryRun, 0)
	},
}

func init() {
	publishCmd.Flags().String("space", "", "Confluence space key (required)")
	publishCmd.Flags().String("parent", "", "parent page ID (required)")
	publishCmd.Flags().Bool("dry-run", false, "preview without creating pages")
	_ = publishCmd.MarkFlagRequired("space")
	_ = publishCmd.MarkFlagRequired("parent")
}

// buildLinkMap scans all markdown files in the directory tree and maps
// relative paths to page titles. This allows the converter to resolve
// [link text](./file.md) to the correct Confluence page title.
func buildLinkMap(rootDir string) map[string]string {
	linkMap := make(map[string]string)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		title := titleFromMarkdown(path)

		// Get path relative to rootDir
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}

		// Store with multiple key variants so links from any subdirectory can resolve:
		// 1. Bare filename: "piv-loops.md"
		// 2. Relative from root: "workflow/piv-loops.md"
		// 3. With ./ prefix: "./workflow/piv-loops.md"
		linkMap[filepath.Base(path)] = title
		linkMap[relPath] = title
		linkMap["./"+relPath] = title

		// For INDEX.md files, also map the directory name
		if filepath.Base(path) == "INDEX.md" {
			linkMap[filepath.Dir(relPath)+"/INDEX.md"] = title
		}

		return nil
	})

	return linkMap
}

func publishDir(client *conflsvc.Client, dir, space, parentID, baseURL string, linkMap map[string]string, dryRun bool, depth int) error {
	indent := strings.Repeat("  ", depth)

	// Check for INDEX.md — it becomes the parent page for this directory
	indexPath := filepath.Join(dir, "INDEX.md")
	currentParentID := parentID

	if _, err := os.Stat(indexPath); err == nil {
		title := titleFromMarkdown(indexPath)
		pageID := frontmatterValue(indexPath, "confluence_page_id")

		if dryRun {
			if pageID != "" {
				fmt.Printf("%s📝 UPDATE page: %q (ID: %s, from INDEX.md)\n", indent, title, pageID)
			} else {
				fmt.Printf("%s📄 CREATE page: %q (from INDEX.md) under parent %s\n", indent, title, currentParentID)
			}
			currentParentID = fmt.Sprintf("<%s>", title)
		} else {
			content, err := readAndConvertMarkdownWithLinks(indexPath, linkMap)
			if err != nil {
				return fmt.Errorf("converting %s: %w", indexPath, err)
			}

			page, err := upsertPage(client, space, currentParentID, title, content, pageID)
			if err != nil {
				return fmt.Errorf("publishing page from %s: %w", indexPath, err)
			}
			currentParentID = page.ID
			if pageID != "" {
				fmt.Printf("%s📝 Updated: %s — %s (ID: %s)\n", indent, filepath.Base(dir), title, page.ID)
			} else {
				fmt.Printf("%s✅ Created: %s — %s (ID: %s)\n", indent, filepath.Base(dir), title, page.ID)
				// Write page ID and URL back to frontmatter
				updateFrontmatter(indexPath, page.ID, buildPageURL(baseURL, space, page.ID, title))
			}
		}
	}

	// Process subdirectories first (they have their own INDEX.md)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	// Subdirectories
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		subdir := filepath.Join(dir, entry.Name())
		if err := publishDir(client, subdir, space, currentParentID, baseURL, linkMap, dryRun, depth+1); err != nil {
			return err
		}
	}

	// Markdown files (skip INDEX.md — already processed)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "INDEX.md" {
			continue
		}
		filePath := filepath.Join(dir, entry.Name())
		title := titleFromMarkdown(filePath)
		pageID := frontmatterValue(filePath, "confluence_page_id")

		if dryRun {
			if pageID != "" {
				fmt.Printf("%s  📝 UPDATE page: %q (ID: %s, from %s)\n", indent, title, pageID, entry.Name())
			} else {
				fmt.Printf("%s  📄 CREATE page: %q (from %s) under parent %s\n", indent, title, entry.Name(), currentParentID)
			}
		} else {
			content, err := readAndConvertMarkdownWithLinks(filePath, linkMap)
			if err != nil {
				return fmt.Errorf("converting %s: %w", filePath, err)
			}

			page, err := upsertPage(client, space, currentParentID, title, content, pageID)
			if err != nil {
				return fmt.Errorf("publishing page from %s: %w", filePath, err)
			}
			if pageID != "" {
				fmt.Printf("%s  📝 Updated: %s — %s (ID: %s)\n", indent, entry.Name(), title, page.ID)
			} else {
				fmt.Printf("%s  ✅ Created: %s — %s (ID: %s)\n", indent, entry.Name(), title, page.ID)
				// Write page ID and URL back to frontmatter
				updateFrontmatter(filePath, page.ID, buildPageURL(baseURL, space, page.ID, title))
			}
		}
	}

	return nil
}

// updateExistingPage fetches the current version and updates the page.
func updateExistingPage(client *conflsvc.Client, pageID, title, content string) (*conflsvc.Page, error) {
	existing, err := client.GetPage(pageID)
	if err != nil {
		return nil, fmt.Errorf("fetching page %s: %w", pageID, err)
	}
	newVersion := 1
	if existing.Version != nil {
		newVersion = existing.Version.Number + 1
	}
	return client.UpdatePage(pageID, title, content, newVersion)
}

// upsertPage creates or updates a Confluence page. If pageID is set, it updates.
// If pageID is empty, it searches by title in the space — if found, updates the
// existing page; otherwise creates a new one.
func upsertPage(client *conflsvc.Client, space, parentID, title, content, pageID string) (*conflsvc.Page, error) {
	if pageID != "" {
		return updateExistingPage(client, pageID, title, content)
	}

	// Search for existing page by title in the space
	existing, err := client.FindPageByTitle(space, title)
	if err != nil {
		// Search failed — fall through to create
		fmt.Printf("    (search failed: %v, will create new page)\n", err)
	}
	if existing != nil {
		// Page exists — update it
		return updateExistingPage(client, existing.ID, title, content)
	}

	// No existing page — create new
	return client.CreatePage(space, parentID, title, content)
}

// titleFromMarkdown extracts a title from a markdown file.
// Uses the frontmatter title, or the first # heading, or the filename.
func titleFromMarkdown(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return filenameToTitle(path)
	}
	content := string(data)

	// Check frontmatter for title
	if strings.HasPrefix(content, "---") {
		rest := content[3:]
		idx := strings.Index(rest, "\n---")
		if idx != -1 {
			frontmatter := rest[:idx]
			for _, line := range strings.Split(frontmatter, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "title:") {
					title := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
					title = strings.Trim(title, `"'`)
					if title != "" {
						return title
					}
				}
			}
		}
	}

	// Check for first # heading
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}

	return filenameToTitle(path)
}

func filenameToTitle(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return strings.Title(name)
}

// readAndConvertMarkdown reads a markdown file and converts to Confluence storage format.
func readAndConvertMarkdown(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}
	return conflsvc.MarkdownToStorage(string(data)), nil
}

// readAndConvertMarkdownWithLinks reads a markdown file and converts to Confluence
// storage format, resolving relative links using the provided link map.
func readAndConvertMarkdownWithLinks(path string, linkMap map[string]string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}

	// Build a file-relative link map: resolve paths relative to the current file's directory
	fileDir := filepath.Dir(path)
	resolvedMap := make(map[string]string, len(linkMap))
	for k, v := range linkMap {
		resolvedMap[k] = v
	}

	// Also add entries resolved from the current file's directory
	// so that "../organization/governance.md" resolves correctly
	rootDir := findRootDir(path, linkMap)
	if rootDir != "" {
		for relFromRoot, title := range linkMap {
			// For each file in the link map, compute how it would be referenced
			// from the current file's directory
			absTarget := filepath.Join(rootDir, relFromRoot)
			relFromFile, err := filepath.Rel(fileDir, absTarget)
			if err == nil {
				resolvedMap[relFromFile] = title
				resolvedMap["./"+relFromFile] = title
			}
		}
	}

	return conflsvc.MarkdownToStorageWithLinks(string(data), resolvedMap), nil
}

// findRootDir determines the root directory by checking which linkMap key
// matches the current file's relative position.
func findRootDir(filePath string, linkMap map[string]string) string {
	base := filepath.Base(filePath)
	dir := filepath.Dir(filePath)

	// Try matching the filename in the link map to find the root
	for relPath := range linkMap {
		if filepath.Base(relPath) == base && strings.Contains(relPath, "/") {
			// relPath is like "workflow/bmad-agents.md"
			// filePath is like "/full/path/to/ai-dev/workflow/bmad-agents.md"
			// rootDir would be "/full/path/to/ai-dev/"
			suffix := relPath
			if strings.HasSuffix(filePath, suffix) {
				return strings.TrimSuffix(filePath, suffix)
			}
		}
	}

	// Fallback: walk up until we can't find more .md files
	// This is a rough heuristic
	return filepath.Dir(dir)
}

// frontmatterValue extracts a value from YAML frontmatter by key.
// Returns empty string if not found or empty.
func frontmatterValue(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return ""
	}
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return ""
	}
	frontmatter := rest[:idx]
	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+":") {
			val := strings.TrimSpace(strings.TrimPrefix(line, key+":"))
			val = strings.Trim(val, `"'`)
			return val
		}
	}
	return ""
}

// updateFrontmatter updates or adds confluence_page_id and confluence_url in the file's YAML frontmatter.
func updateFrontmatter(path, pageID, pageURL string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)

	if !strings.HasPrefix(content, "---") {
		// No frontmatter — add it
		fm := fmt.Sprintf("---\nconfluence_page_id: \"%s\"\nconfluence_url: \"%s\"\n---\n\n", pageID, pageURL)
		_ = os.WriteFile(path, []byte(fm+content), 0644)
		return
	}

	// Update existing frontmatter
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return
	}
	frontmatter := rest[:idx]
	afterFrontmatter := rest[idx+4:]

	// Replace or add confluence_page_id
	rePageID := regexp.MustCompile(`(?m)^confluence_page_id:.*$`)
	if rePageID.MatchString(frontmatter) {
		frontmatter = rePageID.ReplaceAllString(frontmatter, fmt.Sprintf(`confluence_page_id: "%s"`, pageID))
	} else {
		frontmatter += fmt.Sprintf("\nconfluence_page_id: \"%s\"", pageID)
	}

	// Replace or add confluence_url
	reURL := regexp.MustCompile(`(?m)^confluence_url:.*$`)
	if reURL.MatchString(frontmatter) {
		frontmatter = reURL.ReplaceAllString(frontmatter, fmt.Sprintf(`confluence_url: "%s"`, pageURL))
	} else {
		frontmatter += fmt.Sprintf("\nconfluence_url: \"%s\"", pageURL)
	}

	newContent := "---" + frontmatter + "\n---" + afterFrontmatter
	_ = os.WriteFile(path, []byte(newContent), 0644)
}

// buildPageURL constructs a Confluence page URL.
func buildPageURL(baseURL, spaceKey, pageID, title string) string {
	encodedTitle := strings.ReplaceAll(title, " ", "+")
	if strings.Contains(baseURL, "atlassian.net") {
		return fmt.Sprintf("%s/wiki/spaces/%s/pages/%s/%s", baseURL, spaceKey, pageID, encodedTitle)
	}
	return fmt.Sprintf("%s/spaces/%s/pages/%s/%s", baseURL, spaceKey, pageID, encodedTitle)
}

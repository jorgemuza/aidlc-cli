package confluence

import (
	"fmt"
	"net/url"

	"github.com/jorgemuza/aidlc-cli/internal/config"
	"github.com/jorgemuza/aidlc-cli/internal/service"
)

// Client provides Confluence REST API operations.
type Client struct {
	service.BaseService
}

// NewClient creates a Confluence API client from a BaseService.
func NewClient(base service.BaseService) *Client {
	return &Client{BaseService: base}
}

// ClientFromService extracts the Client from a Service.
func ClientFromService(s service.Service) (*Client, error) {
	cs, ok := s.(*svc)
	if !ok {
		return nil, fmt.Errorf("service is not a Confluence service")
	}
	return NewClient(cs.BaseService), nil
}

func (c *Client) apiPrefix() string {
	if c.Conn.Variant == config.VariantCloud {
		return "/wiki/rest/api"
	}
	return "/rest/api"
}

// Page represents a Confluence page.
type Page struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Title   string    `json:"title"`
	Version *Version  `json:"version,omitempty"`
	Body    *Body     `json:"body,omitempty"`
	Links   *Links    `json:"_links,omitempty"`
	Space   *Space    `json:"space,omitempty"`
}

// Version represents a page version.
type Version struct {
	Number  int    `json:"number"`
	Message string `json:"message,omitempty"`
}

// Body holds page content.
type Body struct {
	Storage *Storage `json:"storage,omitempty"`
}

// Storage holds Confluence storage format content.
type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

// Links holds page links.
type Links struct {
	WebUI string `json:"webui,omitempty"`
	Base  string `json:"base,omitempty"`
}

// Space represents a Confluence space.
type Space struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}

// ChildrenResult is the response from getting child pages.
type ChildrenResult struct {
	Results []Page `json:"results"`
	Size    int    `json:"size"`
}

// SearchResult is the response from a CQL search.
type SearchResult struct {
	Results []Page `json:"results"`
	Size    int    `json:"size"`
}

// FindPageByTitle searches for a page by exact title within a space.
// Returns nil if no page is found.
func (c *Client) FindPageByTitle(spaceKey, title string) (*Page, error) {
	cql := fmt.Sprintf(`space="%s" AND title="%s" AND type=page`, spaceKey, title)
	path := fmt.Sprintf("%s/content?cql=%s&expand=version,space", c.apiPrefix(), url.QueryEscape(cql))
	var result SearchResult
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("searching for page %q: %w", title, err)
	}
	if result.Size == 0 {
		return nil, nil
	}
	return &result.Results[0], nil
}

// GetPage fetches a page by ID with body content.
func (c *Client) GetPage(id string) (*Page, error) {
	path := fmt.Sprintf("%s/content/%s?expand=body.storage,version,space", c.apiPrefix(), url.PathEscape(id))
	var page Page
	if err := c.DoGet(path, &page); err != nil {
		return nil, fmt.Errorf("getting page %s: %w", id, err)
	}
	return &page, nil
}

// GetChildPages returns child pages of a given page.
func (c *Client) GetChildPages(parentID string) ([]Page, error) {
	path := fmt.Sprintf("%s/content/%s/child/page?limit=100&expand=version", c.apiPrefix(), url.PathEscape(parentID))
	var result ChildrenResult
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("getting children of %s: %w", parentID, err)
	}
	return result.Results, nil
}

// CreatePage creates a new page under a parent and sets it to wide width.
func (c *Client) CreatePage(spaceKey, parentID, title, storageBody string) (*Page, error) {
	req := map[string]any{
		"type":  "page",
		"title": title,
		"space": map[string]string{"key": spaceKey},
		"body": map[string]any{
			"storage": map[string]string{
				"value":          storageBody,
				"representation": "storage",
			},
		},
	}
	if parentID != "" {
		req["ancestors"] = []map[string]string{{"id": parentID}}
	}

	var page Page
	if err := c.DoPost(c.apiPrefix()+"/content", req, &page); err != nil {
		return nil, fmt.Errorf("creating page %q: %w", title, err)
	}

	// Set wide width by default
	_ = c.SetPageWidth(page.ID, "full-width")

	return &page, nil
}

// SetPageWidth sets the content appearance (width) of a page.
// Use "full-width" for wide or "fixed" for default.
func (c *Client) SetPageWidth(pageID, appearance string) error {
	path := fmt.Sprintf("%s/content/%s/property/content-appearance-draft", c.apiPrefix(), url.PathEscape(pageID))

	// Try to get existing property first
	var existing struct {
		Version struct {
			Number int `json:"number"`
		} `json:"version"`
	}
	verNum := 1
	if err := c.DoGet(path, &existing); err == nil {
		verNum = existing.Version.Number + 1
	}

	req := map[string]any{
		"key":   "content-appearance-draft",
		"value": appearance,
		"version": map[string]int{
			"number": verNum,
		},
	}

	if verNum == 1 {
		if err := c.DoPost(path, req, nil); err != nil {
			return fmt.Errorf("setting page width: %w", err)
		}
	} else {
		if err := c.DoPut(path, req, nil); err != nil {
			return fmt.Errorf("setting page width: %w", err)
		}
	}

	// Also set published appearance
	pubPath := fmt.Sprintf("%s/content/%s/property/content-appearance-published", c.apiPrefix(), url.PathEscape(pageID))
	verNum = 1
	if err := c.DoGet(pubPath, &existing); err == nil {
		verNum = existing.Version.Number + 1
	}
	pubReq := map[string]any{
		"key":   "content-appearance-published",
		"value": appearance,
		"version": map[string]int{
			"number": verNum,
		},
	}
	if verNum == 1 {
		_ = c.DoPost(pubPath, pubReq, nil)
	} else {
		_ = c.DoPut(pubPath, pubReq, nil)
	}

	return nil
}

// DeletePage moves a page to the trash.
func (c *Client) DeletePage(id string) error {
	path := fmt.Sprintf("%s/content/%s", c.apiPrefix(), url.PathEscape(id))
	if err := c.DoDelete(path); err != nil {
		return fmt.Errorf("deleting page %s: %w", id, err)
	}
	return nil
}

// UpdatePage updates an existing page.
func (c *Client) UpdatePage(id, title, storageBody string, version int) (*Page, error) {
	req := map[string]any{
		"type":  "page",
		"title": title,
		"version": map[string]int{
			"number": version,
		},
		"body": map[string]any{
			"storage": map[string]string{
				"value":          storageBody,
				"representation": "storage",
			},
		},
	}

	var page Page
	path := fmt.Sprintf("%s/content/%s", c.apiPrefix(), url.PathEscape(id))
	if err := c.DoPut(path, req, &page); err != nil {
		return nil, fmt.Errorf("updating page %s: %w", id, err)
	}
	return &page, nil
}

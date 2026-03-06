package bitbucket

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jorgemuza/orbit/internal/config"
	"github.com/jorgemuza/orbit/internal/service"
)

// Client provides Bitbucket REST API operations.
// Supports both Server/Data Center (/rest/api/latest/) and Cloud (/2.0/).
type Client struct {
	service.BaseService
}

func NewClient(base service.BaseService) *Client {
	return &Client{BaseService: base}
}

func ClientFromService(s service.Service) (*Client, error) {
	bs, ok := s.(*svc)
	if !ok {
		return nil, fmt.Errorf("service is not a Bitbucket service")
	}
	return NewClient(bs.BaseService), nil
}

func (c *Client) isCloud() bool {
	return c.Conn.Variant == config.VariantCloud
}

func (c *Client) apiPrefix() string {
	if c.isCloud() {
		return ""
	}
	return "/rest/api/latest"
}

// --- Types (Server/Data Center) ---

type PagedResponse[T any] struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	Start      int  `json:"start"`
	IsLastPage bool `json:"isLastPage"`
	Values     []T  `json:"values"`
}

type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links,omitempty"`
}

type Repository struct {
	ID            int      `json:"id"`
	Slug          string   `json:"slug"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	State         string   `json:"state"`
	Forkable      bool     `json:"forkable"`
	Public        bool     `json:"public"`
	Project       *Project `json:"project,omitempty"`
	ScmID         string   `json:"scmId"`
	StatusMessage string   `json:"statusMessage"`
	Links         Links    `json:"links,omitempty"`
}

type Branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}

type Tag struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	Hash            string `json:"hash"`
}

type Commit struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
	Message   string `json:"message"`
	Author    *User  `json:"author,omitempty"`
	Committer *User  `json:"committer,omitempty"`
	AuthorTS  int64  `json:"authorTimestamp"`
	Parents   []struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
	} `json:"parents"`
}

type PullRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Open        bool   `json:"open"`
	Closed      bool   `json:"closed"`
	CreatedDate int64  `json:"createdDate"`
	UpdatedDate int64  `json:"updatedDate"`
	ClosedDate  int64  `json:"closedDate,omitempty"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
	Author      *PRParticipant `json:"author,omitempty"`
	Reviewers   []PRParticipant `json:"reviewers"`
	Links       Links  `json:"links,omitempty"`
}

type PRParticipant struct {
	User     *User  `json:"user"`
	Role     string `json:"role"`
	Approved bool   `json:"approved"`
	Status   string `json:"status"`
}

type Ref struct {
	ID           string      `json:"id"`
	DisplayID    string      `json:"displayId"`
	LatestCommit string      `json:"latestCommit"`
	Repository   *Repository `json:"repository,omitempty"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type PRActivity struct {
	ID          int    `json:"id"`
	Action      string `json:"action"`
	CommentAction string `json:"commentAction,omitempty"`
	Comment     *PRComment `json:"comment,omitempty"`
	CreatedDate int64  `json:"createdDate"`
	User        *User  `json:"user,omitempty"`
}

type PRComment struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Author    *User  `json:"author,omitempty"`
	CreatedDate int64 `json:"createdDate"`
	UpdatedDate int64 `json:"updatedDate"`
}

type Links struct {
	Self  []Link `json:"self,omitempty"`
	Clone []Link `json:"clone,omitempty"`
}

type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// --- Project operations ---

func (c *Client) GetProject(projectKey string) (*Project, error) {
	var p Project
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s", c.apiPrefix(), projectKey), &p); err != nil {
		return nil, fmt.Errorf("getting project: %w", err)
	}
	return &p, nil
}

func (c *Client) ListProjects(limit int) ([]Project, error) {
	var resp PagedResponse[Project]
	if err := c.DoGet(fmt.Sprintf("%s/projects?limit=%d", c.apiPrefix(), limit), &resp); err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	return resp.Values, nil
}

// --- Repository operations ---

func (c *Client) GetRepository(projectKey, repoSlug string) (*Repository, error) {
	var r Repository
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos/%s", c.apiPrefix(), projectKey, repoSlug), &r); err != nil {
		return nil, fmt.Errorf("getting repository: %w", err)
	}
	return &r, nil
}

func (c *Client) ListRepositories(projectKey string, limit int) ([]Repository, error) {
	var resp PagedResponse[Repository]
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos?limit=%d", c.apiPrefix(), projectKey, limit), &resp); err != nil {
		return nil, fmt.Errorf("listing repositories: %w", err)
	}
	return resp.Values, nil
}

// --- Branch operations ---

func (c *Client) ListBranches(projectKey, repoSlug, filter string, limit int) ([]Branch, error) {
	path := fmt.Sprintf("%s/projects/%s/repos/%s/branches?limit=%d", c.apiPrefix(), projectKey, repoSlug, limit)
	if filter != "" {
		path += "&filterText=" + url.QueryEscape(filter)
	}
	var resp PagedResponse[Branch]
	if err := c.DoGet(path, &resp); err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}
	return resp.Values, nil
}

func (c *Client) GetDefaultBranch(projectKey, repoSlug string) (*Branch, error) {
	var b Branch
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos/%s/default-branch", c.apiPrefix(), projectKey, repoSlug), &b); err != nil {
		return nil, fmt.Errorf("getting default branch: %w", err)
	}
	return &b, nil
}

func (c *Client) CreateBranch(projectKey, repoSlug, name, startPoint string) error {
	body := map[string]string{"name": name, "startPoint": startPoint}
	return c.DoPost(fmt.Sprintf("%s/projects/%s/repos/%s/branches", c.apiPrefix(), projectKey, repoSlug), body, nil)
}

func (c *Client) DeleteBranch(projectKey, repoSlug, name string) error {
	body := map[string]any{"name": name, "dryRun": false}
	return c.DoRequest("DELETE", fmt.Sprintf("%s/projects/%s/repos/%s/branches", c.apiPrefix(), projectKey, repoSlug), body, nil)
}

// --- Tag operations ---

func (c *Client) ListTags(projectKey, repoSlug, filter string, limit int) ([]Tag, error) {
	path := fmt.Sprintf("%s/projects/%s/repos/%s/tags?limit=%d", c.apiPrefix(), projectKey, repoSlug, limit)
	if filter != "" {
		path += "&filterText=" + url.QueryEscape(filter)
	}
	var resp PagedResponse[Tag]
	if err := c.DoGet(path, &resp); err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	return resp.Values, nil
}

func (c *Client) CreateTag(projectKey, repoSlug, name, startPoint, message string) error {
	body := map[string]string{"name": name, "startPoint": startPoint}
	if message != "" {
		body["message"] = message
	}
	return c.DoPost(fmt.Sprintf("%s/projects/%s/repos/%s/tags", c.apiPrefix(), projectKey, repoSlug), body, nil)
}

// --- Commit operations ---

func (c *Client) ListCommits(projectKey, repoSlug, branch string, limit int) ([]Commit, error) {
	path := fmt.Sprintf("%s/projects/%s/repos/%s/commits?limit=%d", c.apiPrefix(), projectKey, repoSlug, limit)
	if branch != "" {
		path += "&until=" + url.QueryEscape(branch)
	}
	var resp PagedResponse[Commit]
	if err := c.DoGet(path, &resp); err != nil {
		return nil, fmt.Errorf("listing commits: %w", err)
	}
	return resp.Values, nil
}

func (c *Client) GetCommit(projectKey, repoSlug, commitID string) (*Commit, error) {
	var cm Commit
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos/%s/commits/%s", c.apiPrefix(), projectKey, repoSlug, commitID), &cm); err != nil {
		return nil, fmt.Errorf("getting commit: %w", err)
	}
	return &cm, nil
}

// --- Pull Request operations ---

func (c *Client) ListPullRequests(projectKey, repoSlug, state string, limit int) ([]PullRequest, error) {
	path := fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests?limit=%d", c.apiPrefix(), projectKey, repoSlug, limit)
	if state != "" {
		path += "&state=" + url.QueryEscape(strings.ToUpper(state))
	}
	var resp PagedResponse[PullRequest]
	if err := c.DoGet(path, &resp); err != nil {
		return nil, fmt.Errorf("listing pull requests: %w", err)
	}
	return resp.Values, nil
}

func (c *Client) GetPullRequest(projectKey, repoSlug string, prID int) (*PullRequest, error) {
	var pr PullRequest
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d", c.apiPrefix(), projectKey, repoSlug, prID), &pr); err != nil {
		return nil, fmt.Errorf("getting pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) CreatePullRequest(projectKey, repoSlug, title, description, fromBranch, toBranch string, reviewerSlugs []string) (*PullRequest, error) {
	body := map[string]any{
		"title":       title,
		"description": description,
		"fromRef":     map[string]string{"id": "refs/heads/" + fromBranch},
		"toRef":       map[string]string{"id": "refs/heads/" + toBranch},
	}
	if len(reviewerSlugs) > 0 {
		reviewers := make([]map[string]any, len(reviewerSlugs))
		for i, slug := range reviewerSlugs {
			reviewers[i] = map[string]any{"user": map[string]string{"name": slug}}
		}
		body["reviewers"] = reviewers
	}
	var pr PullRequest
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests", c.apiPrefix(), projectKey, repoSlug), body, &pr); err != nil {
		return nil, fmt.Errorf("creating pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) MergePullRequest(projectKey, repoSlug string, prID int, version int) (*PullRequest, error) {
	var pr PullRequest
	path := fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/merge?version=%d", c.apiPrefix(), projectKey, repoSlug, prID, version)
	if err := c.DoPost(path, nil, &pr); err != nil {
		return nil, fmt.Errorf("merging pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) DeclinePullRequest(projectKey, repoSlug string, prID int, version int) (*PullRequest, error) {
	var pr PullRequest
	path := fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/decline?version=%d", c.apiPrefix(), projectKey, repoSlug, prID, version)
	if err := c.DoPost(path, nil, &pr); err != nil {
		return nil, fmt.Errorf("declining pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) ListPRActivities(projectKey, repoSlug string, prID, limit int) ([]PRActivity, error) {
	var resp PagedResponse[PRActivity]
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/activities?limit=%d", c.apiPrefix(), projectKey, repoSlug, prID, limit), &resp); err != nil {
		return nil, fmt.Errorf("listing PR activities: %w", err)
	}
	return resp.Values, nil
}

func (c *Client) CommentPullRequest(projectKey, repoSlug string, prID int, text string) (*PRComment, error) {
	var comment PRComment
	body := map[string]string{"text": text}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/comments", c.apiPrefix(), projectKey, repoSlug, prID), body, &comment); err != nil {
		return nil, fmt.Errorf("commenting on pull request: %w", err)
	}
	return &comment, nil
}

// --- User operations ---

func (c *Client) CurrentUser() (*User, error) {
	// Bitbucket Server doesn't have a direct /user endpoint like Cloud.
	// We use the application-properties endpoint already, so let's use
	// the recently-used user list or just return nil for now.
	// Actually, there's no straightforward "who am I" in BB Server.
	// We'll use /plugins/servlet/applinks/whoami or /rest/api/latest/users?filter=...
	// The simplest approach: /rest/api/latest/users is available to admins.
	// Best bet for server: no direct endpoint; skip for now.
	return nil, fmt.Errorf("current user endpoint not available on Bitbucket Server")
}

func (c *Client) ListUsers(filter string, limit int) ([]User, error) {
	path := fmt.Sprintf("%s/users?limit=%d", c.apiPrefix(), limit)
	if filter != "" {
		path += "&filter=" + url.QueryEscape(filter)
	}
	var resp PagedResponse[User]
	if err := c.DoGet(path, &resp); err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return resp.Values, nil
}

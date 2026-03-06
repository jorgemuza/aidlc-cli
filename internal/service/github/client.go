package github

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jorgemuza/orbit/internal/service"
)

// Client provides GitHub REST API operations.
type Client struct {
	service.BaseService
}

// NewClient creates a GitHub API client from a BaseService.
func NewClient(base service.BaseService) *Client {
	return &Client{BaseService: base}
}

// ClientFromService extracts the Client from a Service.
func ClientFromService(s service.Service) (*Client, error) {
	gs, ok := s.(*svc)
	if !ok {
		return nil, fmt.Errorf("service is not a GitHub service")
	}
	return NewClient(gs.BaseService), nil
}

// --- Types ---

type Repository struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	SSHURL        string `json:"ssh_url"`
	DefaultBranch string `json:"default_branch"`
	Language      string `json:"language"`
	StargazersCount int  `json:"stargazers_count"`
	ForksCount    int    `json:"forks_count"`
	OpenIssuesCount int  `json:"open_issues_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
	Owner         *User  `json:"owner,omitempty"`
}

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

type Branch struct {
	Name      string  `json:"name"`
	Protected bool    `json:"protected"`
	Commit    *struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit,omitempty"`
}

type BranchDetail struct {
	Name      string  `json:"name"`
	Protected bool    `json:"protected"`
	Commit    *Commit `json:"commit,omitempty"`
}

type Commit struct {
	SHA       string       `json:"sha"`
	HTMLURL   string       `json:"html_url"`
	Commit    *CommitData  `json:"commit,omitempty"`
	Author    *User        `json:"author,omitempty"`
	Committer *User        `json:"committer,omitempty"`
	Parents   []CommitRef  `json:"parents"`
}

type CommitData struct {
	Message   string       `json:"message"`
	Author    *CommitActor `json:"author,omitempty"`
	Committer *CommitActor `json:"committer,omitempty"`
}

type CommitActor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type CommitRef struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Tag struct {
	Name       string  `json:"name"`
	ZipballURL string  `json:"zipball_url"`
	TarballURL string  `json:"tarball_url"`
	Commit     *struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit,omitempty"`
}

type PullRequest struct {
	ID        int      `json:"id"`
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	State     string   `json:"state"`
	Draft     bool     `json:"draft"`
	HTMLURL   string   `json:"html_url"`
	Head      *PRRef   `json:"head,omitempty"`
	Base      *PRRef   `json:"base,omitempty"`
	User      *User    `json:"user,omitempty"`
	Assignees []User   `json:"assignees"`
	Labels    []Label  `json:"labels"`
	MergedAt  string   `json:"merged_at"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	ClosedAt  string   `json:"closed_at"`
	Mergeable *bool    `json:"mergeable,omitempty"`
	Comments  int      `json:"comments"`
}

type PRRef struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Issue struct {
	ID        int     `json:"id"`
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	State     string  `json:"state"`
	HTMLURL   string  `json:"html_url"`
	User      *User   `json:"user,omitempty"`
	Assignees []User  `json:"assignees"`
	Labels    []Label `json:"labels"`
	Milestone *struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"milestone,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
}

type Comment struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	User      *User  `json:"user,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Release struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	HTMLURL     string `json:"html_url"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Author      *User  `json:"author,omitempty"`
}

type WorkflowRun struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	HeadBranch   string `json:"head_branch"`
	HeadSHA      string `json:"head_sha"`
	HTMLURL      string `json:"html_url"`
	Event        string `json:"event"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// --- Repository operations ---

func (c *Client) GetRepo(owner, repo string) (*Repository, error) {
	var r Repository
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s", owner, repo), &r); err != nil {
		return nil, fmt.Errorf("getting repo: %w", err)
	}
	return &r, nil
}

func (c *Client) ListUserRepos(perPage int) ([]Repository, error) {
	var repos []Repository
	if err := c.DoGet(fmt.Sprintf("/user/repos?per_page=%d&sort=pushed&direction=desc", perPage), &repos); err != nil {
		return nil, fmt.Errorf("listing repos: %w", err)
	}
	return repos, nil
}

func (c *Client) ListOrgRepos(org string, perPage int) ([]Repository, error) {
	var repos []Repository
	if err := c.DoGet(fmt.Sprintf("/orgs/%s/repos?per_page=%d&sort=pushed&direction=desc", url.PathEscape(org), perPage), &repos); err != nil {
		return nil, fmt.Errorf("listing org repos: %w", err)
	}
	return repos, nil
}

func ownerRepo(fullName string) (string, string, error) {
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo format %q, expected owner/repo", fullName)
	}
	return parts[0], parts[1], nil
}

// --- Branch operations ---

func (c *Client) ListBranches(owner, repo string, perPage int) ([]Branch, error) {
	var branches []Branch
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/branches?per_page=%d", owner, repo, perPage), &branches); err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}
	return branches, nil
}

func (c *Client) GetBranch(owner, repo, branch string) (*BranchDetail, error) {
	var b BranchDetail
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/branches/%s", owner, repo, url.PathEscape(branch)), &b); err != nil {
		return nil, fmt.Errorf("getting branch: %w", err)
	}
	return &b, nil
}

// --- Tag operations ---

func (c *Client) ListTags(owner, repo string, perPage int) ([]Tag, error) {
	var tags []Tag
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/tags?per_page=%d", owner, repo, perPage), &tags); err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	return tags, nil
}

// --- Commit operations ---

func (c *Client) ListCommits(owner, repo, sha string, perPage int) ([]Commit, error) {
	path := fmt.Sprintf("/repos/%s/%s/commits?per_page=%d", owner, repo, perPage)
	if sha != "" {
		path += "&sha=" + url.QueryEscape(sha)
	}
	var commits []Commit
	if err := c.DoGet(path, &commits); err != nil {
		return nil, fmt.Errorf("listing commits: %w", err)
	}
	return commits, nil
}

func (c *Client) GetCommit(owner, repo, ref string) (*Commit, error) {
	var cm Commit
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/commits/%s", owner, repo, url.PathEscape(ref)), &cm); err != nil {
		return nil, fmt.Errorf("getting commit: %w", err)
	}
	return &cm, nil
}

// --- Pull Request operations ---

func (c *Client) ListPullRequests(owner, repo, state string, perPage int) ([]PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls?per_page=%d&sort=updated&direction=desc", owner, repo, perPage)
	if state != "" {
		path += "&state=" + url.QueryEscape(state)
	}
	var prs []PullRequest
	if err := c.DoGet(path, &prs); err != nil {
		return nil, fmt.Errorf("listing pull requests: %w", err)
	}
	return prs, nil
}

func (c *Client) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	var pr PullRequest
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number), &pr); err != nil {
		return nil, fmt.Errorf("getting pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) CreatePullRequest(owner, repo, head, base, title, body string) (*PullRequest, error) {
	var pr PullRequest
	payload := map[string]string{
		"head":  head,
		"base":  base,
		"title": title,
	}
	if body != "" {
		payload["body"] = body
	}
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/pulls", owner, repo), payload, &pr); err != nil {
		return nil, fmt.Errorf("creating pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) MergePullRequest(owner, repo string, number int, mergeMethod string) error {
	payload := map[string]string{}
	if mergeMethod != "" {
		payload["merge_method"] = mergeMethod
	}
	if err := c.DoPut(fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", owner, repo, number), payload, nil); err != nil {
		return fmt.Errorf("merging pull request: %w", err)
	}
	return nil
}

func (c *Client) ListPRComments(owner, repo string, number, perPage int) ([]Comment, error) {
	var comments []Comment
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/issues/%d/comments?per_page=%d", owner, repo, number, perPage), &comments); err != nil {
		return nil, fmt.Errorf("listing PR comments: %w", err)
	}
	return comments, nil
}

func (c *Client) CreatePRComment(owner, repo string, number int, body string) (*Comment, error) {
	var comment Comment
	payload := map[string]string{"body": body}
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, number), payload, &comment); err != nil {
		return nil, fmt.Errorf("creating PR comment: %w", err)
	}
	return &comment, nil
}

// --- Issue operations ---

func (c *Client) ListIssues(owner, repo, state string, labels []string, perPage int) ([]Issue, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues?per_page=%d&sort=updated&direction=desc", owner, repo, perPage)
	if state != "" {
		path += "&state=" + url.QueryEscape(state)
	}
	if len(labels) > 0 {
		path += "&labels=" + url.QueryEscape(strings.Join(labels, ","))
	}
	var issues []Issue
	if err := c.DoGet(path, &issues); err != nil {
		return nil, fmt.Errorf("listing issues: %w", err)
	}
	return issues, nil
}

func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	var issue Issue
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number), &issue); err != nil {
		return nil, fmt.Errorf("getting issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) CreateIssue(owner, repo, title, body string, labels []string, assignees []string) (*Issue, error) {
	var issue Issue
	payload := map[string]any{
		"title": title,
	}
	if body != "" {
		payload["body"] = body
	}
	if len(labels) > 0 {
		payload["labels"] = labels
	}
	if len(assignees) > 0 {
		payload["assignees"] = assignees
	}
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/issues", owner, repo), payload, &issue); err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) UpdateIssue(owner, repo string, number int, updates map[string]any) (*Issue, error) {
	var issue Issue
	path := fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)
	if err := c.DoRequest("PATCH", path, updates, &issue); err != nil {
		return nil, fmt.Errorf("updating issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) ListIssueComments(owner, repo string, number, perPage int) ([]Comment, error) {
	var comments []Comment
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/issues/%d/comments?per_page=%d", owner, repo, number, perPage), &comments); err != nil {
		return nil, fmt.Errorf("listing issue comments: %w", err)
	}
	return comments, nil
}

func (c *Client) CreateIssueComment(owner, repo string, number int, body string) (*Comment, error) {
	var comment Comment
	payload := map[string]string{"body": body}
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, number), payload, &comment); err != nil {
		return nil, fmt.Errorf("creating issue comment: %w", err)
	}
	return &comment, nil
}

// --- Release operations ---

func (c *Client) ListReleases(owner, repo string, perPage int) ([]Release, error) {
	var releases []Release
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/releases?per_page=%d", owner, repo, perPage), &releases); err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}
	return releases, nil
}

func (c *Client) GetRelease(owner, repo string, id int) (*Release, error) {
	var r Release
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/releases/%d", owner, repo, id), &r); err != nil {
		return nil, fmt.Errorf("getting release: %w", err)
	}
	return &r, nil
}

func (c *Client) GetLatestRelease(owner, repo string) (*Release, error) {
	var r Release
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/releases/latest", owner, repo), &r); err != nil {
		return nil, fmt.Errorf("getting latest release: %w", err)
	}
	return &r, nil
}

// --- User operations ---

func (c *Client) CurrentUser() (*User, error) {
	var u User
	if err := c.DoGet("/user", &u); err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return &u, nil
}

func (c *Client) GetUser(username string) (*User, error) {
	var u User
	if err := c.DoGet(fmt.Sprintf("/users/%s", url.PathEscape(username)), &u); err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	return &u, nil
}

// --- Workflow Run operations ---

func (c *Client) ListWorkflowRuns(owner, repo, branch, status string, perPage int) ([]WorkflowRun, error) {
	path := fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=%d", owner, repo, perPage)
	if branch != "" {
		path += "&branch=" + url.QueryEscape(branch)
	}
	if status != "" {
		path += "&status=" + url.QueryEscape(status)
	}
	var result struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("listing workflow runs: %w", err)
	}
	return result.WorkflowRuns, nil
}

func (c *Client) GetWorkflowRun(owner, repo string, runID int) (*WorkflowRun, error) {
	var r WorkflowRun
	if err := c.DoGet(fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, repo, runID), &r); err != nil {
		return nil, fmt.Errorf("getting workflow run: %w", err)
	}
	return &r, nil
}

func (c *Client) CancelWorkflowRun(owner, repo string, runID int) error {
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/cancel", owner, repo, runID), nil, nil); err != nil {
		return fmt.Errorf("canceling workflow run: %w", err)
	}
	return nil
}

func (c *Client) RerunWorkflowRun(owner, repo string, runID int) error {
	if err := c.DoPost(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/rerun", owner, repo, runID), nil, nil); err != nil {
		return fmt.Errorf("rerunning workflow run: %w", err)
	}
	return nil
}

// OwnerRepo parses "owner/repo" into separate parts.
func OwnerRepo(fullName string) (string, string, error) {
	return ownerRepo(fullName)
}

// OwnerRepoFromArgs parses the first arg as "owner/repo" or uses separate args.
func OwnerRepoFromArgs(args []string) (string, string, error) {
	if len(args) < 1 {
		return "", "", fmt.Errorf("repo argument required (owner/repo)")
	}
	return ownerRepo(args[0])
}


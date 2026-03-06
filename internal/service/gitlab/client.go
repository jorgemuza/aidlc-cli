package gitlab

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jorgemuza/aidlc-cli/internal/service"
)

// Client provides GitLab REST API operations.
type Client struct {
	service.BaseService
}

// NewClient creates a GitLab API client from a BaseService.
func NewClient(base service.BaseService) *Client {
	return &Client{BaseService: base}
}

// ClientFromService extracts the Client from a Service.
func ClientFromService(s service.Service) (*Client, error) {
	gs, ok := s.(*svc)
	if !ok {
		return nil, fmt.Errorf("service is not a GitLab service")
	}
	return NewClient(gs.BaseService), nil
}

const apiV4 = "/api/v4"

// --- Types ---

type Project struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	Description       string     `json:"description"`
	WebURL            string     `json:"web_url"`
	DefaultBranch     string     `json:"default_branch"`
	Visibility        string     `json:"visibility"`
	CreatedAt         string     `json:"created_at"`
	LastActivityAt    string     `json:"last_activity_at"`
	Namespace         *Namespace `json:"namespace,omitempty"`
	StarCount         int        `json:"star_count"`
	ForksCount        int        `json:"forks_count"`
	OpenIssuesCount   int        `json:"open_issues_count"`
}

type Namespace struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	FullPath    string `json:"full_path"`
	Description string `json:"description"`
	WebURL      string `json:"web_url"`
	Visibility  string `json:"visibility"`
}

type Branch struct {
	Name      string  `json:"name"`
	Protected bool    `json:"protected"`
	Merged    bool    `json:"merged"`
	Default   bool    `json:"default"`
	WebURL    string  `json:"web_url"`
	Commit    *Commit `json:"commit,omitempty"`
}

type Commit struct {
	ID             string   `json:"id"`
	ShortID        string   `json:"short_id"`
	Title          string   `json:"title"`
	Message        string   `json:"message"`
	AuthorName     string   `json:"author_name"`
	AuthorEmail    string   `json:"author_email"`
	AuthoredDate   string   `json:"authored_date"`
	CommittedDate  string   `json:"committed_date"`
	CommitterName  string   `json:"committer_name"`
	CommitterEmail string   `json:"committer_email"`
	WebURL         string   `json:"web_url"`
	ParentIDs      []string `json:"parent_ids"`
}

type Tag struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Commit  *Commit `json:"commit,omitempty"`
	Release *struct {
		TagName     string `json:"tag_name"`
		Description string `json:"description"`
	} `json:"release,omitempty"`
}

type MergeRequest struct {
	ID             int      `json:"id"`
	IID            int      `json:"iid"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	State          string   `json:"state"`
	SourceBranch   string   `json:"source_branch"`
	TargetBranch   string   `json:"target_branch"`
	WebURL         string   `json:"web_url"`
	Author         *User    `json:"author,omitempty"`
	Assignee       *User    `json:"assignee,omitempty"`
	Reviewers      []User   `json:"reviewers,omitempty"`
	Labels         []string `json:"labels"`
	MergeStatus    string   `json:"merge_status"`
	HasConflicts   bool     `json:"has_conflicts"`
	Draft          bool     `json:"draft"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	MergedAt       string   `json:"merged_at"`
	ClosedAt       string   `json:"closed_at"`
	UserNotesCount int      `json:"user_notes_count"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	State    string `json:"state"`
	WebURL   string `json:"web_url"`
}

type Pipeline struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	SHA       string `json:"sha"`
	WebURL    string `json:"web_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Source    string `json:"source"`
}

type Job struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Stage      string `json:"stage"`
	Status     string `json:"status"`
	Ref        string `json:"ref"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	Duration   float64 `json:"duration"`
	WebURL     string `json:"web_url"`
	Pipeline   *struct {
		ID int `json:"id"`
	} `json:"pipeline,omitempty"`
	Runner *struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"runner,omitempty"`
}

type Issue struct {
	ID          int      `json:"id"`
	IID         int      `json:"iid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	Labels      []string `json:"labels"`
	Author      *User    `json:"author,omitempty"`
	Assignees   []User   `json:"assignees"`
	WebURL      string   `json:"web_url"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	ClosedAt    string   `json:"closed_at"`
	DueDate     string   `json:"due_date"`
	Milestone   *struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"milestone,omitempty"`
}

type Member struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	State       string `json:"state"`
	AccessLevel int    `json:"access_level"`
	WebURL      string `json:"web_url"`
}

type Note struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author    *User  `json:"author,omitempty"`
	CreatedAt string `json:"created_at"`
	System    bool   `json:"system"`
}

// --- Helper ---

func encodeProject(projectPath string) string {
	return url.PathEscape(projectPath)
}

// --- Project operations ---

func (c *Client) GetProject(idOrPath string) (*Project, error) {
	var p Project
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s", apiV4, encodeProject(idOrPath)), &p); err != nil {
		return nil, fmt.Errorf("getting project: %w", err)
	}
	return &p, nil
}

func (c *Client) ListProjects(search string, perPage int) ([]Project, error) {
	path := fmt.Sprintf("%s/projects?per_page=%d&order_by=last_activity_at&sort=desc&membership=true", apiV4, perPage)
	if search != "" {
		path += "&search=" + url.QueryEscape(search)
	}
	var projects []Project
	if err := c.DoGet(path, &projects); err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	return projects, nil
}

// --- Group operations ---

func (c *Client) GetGroup(idOrPath string) (*Group, error) {
	var g Group
	if err := c.DoGet(fmt.Sprintf("%s/groups/%s", apiV4, encodeProject(idOrPath)), &g); err != nil {
		return nil, fmt.Errorf("getting group: %w", err)
	}
	return &g, nil
}

func (c *Client) ListGroups(search string, perPage int) ([]Group, error) {
	path := fmt.Sprintf("%s/groups?per_page=%d&order_by=name&sort=asc", apiV4, perPage)
	if search != "" {
		path += "&search=" + url.QueryEscape(search)
	}
	var groups []Group
	if err := c.DoGet(path, &groups); err != nil {
		return nil, fmt.Errorf("listing groups: %w", err)
	}
	return groups, nil
}

func (c *Client) ListGroupProjects(groupIDOrPath string, perPage int) ([]Project, error) {
	var projects []Project
	if err := c.DoGet(fmt.Sprintf("%s/groups/%s/projects?per_page=%d&order_by=name&sort=asc&include_subgroups=true", apiV4, encodeProject(groupIDOrPath), perPage), &projects); err != nil {
		return nil, fmt.Errorf("listing group projects: %w", err)
	}
	return projects, nil
}

func (c *Client) ListSubgroups(groupIDOrPath string, perPage int) ([]Group, error) {
	var groups []Group
	if err := c.DoGet(fmt.Sprintf("%s/groups/%s/subgroups?per_page=%d", apiV4, encodeProject(groupIDOrPath), perPage), &groups); err != nil {
		return nil, fmt.Errorf("listing subgroups: %w", err)
	}
	return groups, nil
}

// --- Branch operations ---

func (c *Client) ListBranches(projectID string, search string, perPage int) ([]Branch, error) {
	path := fmt.Sprintf("%s/projects/%s/repository/branches?per_page=%d", apiV4, encodeProject(projectID), perPage)
	if search != "" {
		path += "&search=" + url.QueryEscape(search)
	}
	var branches []Branch
	if err := c.DoGet(path, &branches); err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}
	return branches, nil
}

func (c *Client) GetBranch(projectID, branchName string) (*Branch, error) {
	var b Branch
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repository/branches/%s", apiV4, encodeProject(projectID), url.PathEscape(branchName)), &b); err != nil {
		return nil, fmt.Errorf("getting branch: %w", err)
	}
	return &b, nil
}

func (c *Client) CreateBranch(projectID, branchName, ref string) (*Branch, error) {
	var b Branch
	body := map[string]string{"branch": branchName, "ref": ref}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/repository/branches", apiV4, encodeProject(projectID)), body, &b); err != nil {
		return nil, fmt.Errorf("creating branch: %w", err)
	}
	return &b, nil
}

func (c *Client) DeleteBranch(projectID, branchName string) error {
	if err := c.DoDelete(fmt.Sprintf("%s/projects/%s/repository/branches/%s", apiV4, encodeProject(projectID), url.PathEscape(branchName))); err != nil {
		return fmt.Errorf("deleting branch: %w", err)
	}
	return nil
}

// --- Tag operations ---

func (c *Client) ListTags(projectID string, perPage int) ([]Tag, error) {
	var tags []Tag
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repository/tags?per_page=%d", apiV4, encodeProject(projectID), perPage), &tags); err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}
	return tags, nil
}

func (c *Client) CreateTag(projectID, tagName, ref, message string) (*Tag, error) {
	var t Tag
	body := map[string]string{"tag_name": tagName, "ref": ref}
	if message != "" {
		body["message"] = message
	}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/repository/tags", apiV4, encodeProject(projectID)), body, &t); err != nil {
		return nil, fmt.Errorf("creating tag: %w", err)
	}
	return &t, nil
}

// --- Commit operations ---

func (c *Client) ListCommits(projectID, refName string, perPage int) ([]Commit, error) {
	path := fmt.Sprintf("%s/projects/%s/repository/commits?per_page=%d", apiV4, encodeProject(projectID), perPage)
	if refName != "" {
		path += "&ref_name=" + url.QueryEscape(refName)
	}
	var commits []Commit
	if err := c.DoGet(path, &commits); err != nil {
		return nil, fmt.Errorf("listing commits: %w", err)
	}
	return commits, nil
}

func (c *Client) GetCommit(projectID, sha string) (*Commit, error) {
	var cm Commit
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/repository/commits/%s", apiV4, encodeProject(projectID), url.PathEscape(sha)), &cm); err != nil {
		return nil, fmt.Errorf("getting commit: %w", err)
	}
	return &cm, nil
}

// --- Merge Request operations ---

func (c *Client) ListMergeRequests(projectID, state string, perPage int) ([]MergeRequest, error) {
	path := fmt.Sprintf("%s/projects/%s/merge_requests?per_page=%d&order_by=updated_at&sort=desc", apiV4, encodeProject(projectID), perPage)
	if state != "" {
		path += "&state=" + url.QueryEscape(state)
	}
	var mrs []MergeRequest
	if err := c.DoGet(path, &mrs); err != nil {
		return nil, fmt.Errorf("listing merge requests: %w", err)
	}
	return mrs, nil
}

func (c *Client) GetMergeRequest(projectID string, mrIID int) (*MergeRequest, error) {
	var mr MergeRequest
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/merge_requests/%d", apiV4, encodeProject(projectID), mrIID), &mr); err != nil {
		return nil, fmt.Errorf("getting merge request: %w", err)
	}
	return &mr, nil
}

func (c *Client) CreateMergeRequest(projectID, sourceBranch, targetBranch, title, description string) (*MergeRequest, error) {
	var mr MergeRequest
	body := map[string]string{
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
		"title":         title,
	}
	if description != "" {
		body["description"] = description
	}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/merge_requests", apiV4, encodeProject(projectID)), body, &mr); err != nil {
		return nil, fmt.Errorf("creating merge request: %w", err)
	}
	return &mr, nil
}

func (c *Client) UpdateMergeRequest(projectID string, mrIID int, updates map[string]any) (*MergeRequest, error) {
	var mr MergeRequest
	if err := c.DoPut(fmt.Sprintf("%s/projects/%s/merge_requests/%d", apiV4, encodeProject(projectID), mrIID), updates, &mr); err != nil {
		return nil, fmt.Errorf("updating merge request: %w", err)
	}
	return &mr, nil
}

func (c *Client) MergeMergeRequest(projectID string, mrIID int, squash bool) (*MergeRequest, error) {
	var mr MergeRequest
	body := map[string]any{}
	if squash {
		body["squash"] = true
	}
	if err := c.DoPut(fmt.Sprintf("%s/projects/%s/merge_requests/%d/merge", apiV4, encodeProject(projectID), mrIID), body, &mr); err != nil {
		return nil, fmt.Errorf("merging merge request: %w", err)
	}
	return &mr, nil
}

func (c *Client) ListMRNotes(projectID string, mrIID, perPage int) ([]Note, error) {
	var notes []Note
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/merge_requests/%d/notes?per_page=%d&sort=asc", apiV4, encodeProject(projectID), mrIID, perPage), &notes); err != nil {
		return nil, fmt.Errorf("listing MR notes: %w", err)
	}
	return notes, nil
}

func (c *Client) CreateMRNote(projectID string, mrIID int, body string) (*Note, error) {
	var note Note
	req := map[string]string{"body": body}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/merge_requests/%d/notes", apiV4, encodeProject(projectID), mrIID), req, &note); err != nil {
		return nil, fmt.Errorf("creating MR note: %w", err)
	}
	return &note, nil
}

// --- Pipeline operations ---

func (c *Client) ListPipelines(projectID, ref, status string, perPage int) ([]Pipeline, error) {
	path := fmt.Sprintf("%s/projects/%s/pipelines?per_page=%d&order_by=id&sort=desc", apiV4, encodeProject(projectID), perPage)
	if ref != "" {
		path += "&ref=" + url.QueryEscape(ref)
	}
	if status != "" {
		path += "&status=" + url.QueryEscape(status)
	}
	var pipelines []Pipeline
	if err := c.DoGet(path, &pipelines); err != nil {
		return nil, fmt.Errorf("listing pipelines: %w", err)
	}
	return pipelines, nil
}

func (c *Client) GetPipeline(projectID string, pipelineID int) (*Pipeline, error) {
	var p Pipeline
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/pipelines/%d", apiV4, encodeProject(projectID), pipelineID), &p); err != nil {
		return nil, fmt.Errorf("getting pipeline: %w", err)
	}
	return &p, nil
}

func (c *Client) ListPipelineJobs(projectID string, pipelineID, perPage int) ([]Job, error) {
	var jobs []Job
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/pipelines/%d/jobs?per_page=%d", apiV4, encodeProject(projectID), pipelineID, perPage), &jobs); err != nil {
		return nil, fmt.Errorf("listing pipeline jobs: %w", err)
	}
	return jobs, nil
}

func (c *Client) RetryPipeline(projectID string, pipelineID int) (*Pipeline, error) {
	var p Pipeline
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/pipelines/%d/retry", apiV4, encodeProject(projectID), pipelineID), nil, &p); err != nil {
		return nil, fmt.Errorf("retrying pipeline: %w", err)
	}
	return &p, nil
}

func (c *Client) CancelPipeline(projectID string, pipelineID int) (*Pipeline, error) {
	var p Pipeline
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/pipelines/%d/cancel", apiV4, encodeProject(projectID), pipelineID), nil, &p); err != nil {
		return nil, fmt.Errorf("canceling pipeline: %w", err)
	}
	return &p, nil
}

// --- Issue operations ---

func (c *Client) ListIssues(projectID, state string, labels []string, perPage int) ([]Issue, error) {
	path := fmt.Sprintf("%s/projects/%s/issues?per_page=%d&order_by=updated_at&sort=desc", apiV4, encodeProject(projectID), perPage)
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

func (c *Client) GetIssue(projectID string, issueIID int) (*Issue, error) {
	var issue Issue
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/issues/%d", apiV4, encodeProject(projectID), issueIID), &issue); err != nil {
		return nil, fmt.Errorf("getting issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) CreateIssue(projectID, title, description string, labels []string, assigneeIDs []int) (*Issue, error) {
	var issue Issue
	body := map[string]any{
		"title": title,
	}
	if description != "" {
		body["description"] = description
	}
	if len(labels) > 0 {
		body["labels"] = strings.Join(labels, ",")
	}
	if len(assigneeIDs) > 0 {
		body["assignee_ids"] = assigneeIDs
	}
	if err := c.DoPost(fmt.Sprintf("%s/projects/%s/issues", apiV4, encodeProject(projectID)), body, &issue); err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) UpdateIssue(projectID string, issueIID int, updates map[string]any) (*Issue, error) {
	var issue Issue
	if err := c.DoPut(fmt.Sprintf("%s/projects/%s/issues/%d", apiV4, encodeProject(projectID), issueIID), updates, &issue); err != nil {
		return nil, fmt.Errorf("updating issue: %w", err)
	}
	return &issue, nil
}

// --- Member operations ---

func (c *Client) ListProjectMembers(projectID string, perPage int) ([]Member, error) {
	var members []Member
	if err := c.DoGet(fmt.Sprintf("%s/projects/%s/members/all?per_page=%d", apiV4, encodeProject(projectID), perPage), &members); err != nil {
		return nil, fmt.Errorf("listing project members: %w", err)
	}
	return members, nil
}

// --- User operations ---

func (c *Client) CurrentUser() (*User, error) {
	var u User
	if err := c.DoGet(fmt.Sprintf("%s/user", apiV4), &u); err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return &u, nil
}

func (c *Client) ListUsers(search string, perPage int) ([]User, error) {
	path := fmt.Sprintf("%s/users?per_page=%d", apiV4, perPage)
	if search != "" {
		path += "&search=" + url.QueryEscape(search)
	}
	var users []User
	if err := c.DoGet(path, &users); err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return users, nil
}

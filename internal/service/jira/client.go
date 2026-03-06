package jira

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jorgemuza/orbit/internal/service"
)

// Client provides Jira REST API operations.
type Client struct {
	service.BaseService
}

// NewClient creates a Jira API client from a BaseService.
func NewClient(base service.BaseService) *Client {
	return &Client{BaseService: base}
}

// SearchIssues searches for issues using JQL.
func (c *Client) SearchIssues(jql string, startAt, maxResults int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("jql", jql)
	params.Set("startAt", fmt.Sprintf("%d", startAt))
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("fields", "summary,status,issuetype,priority,assignee,reporter,created,updated,labels,components,fixVersions,parent,subtasks,issuelinks,timetracking,description,resolution")

	var result SearchResult
	if err := c.DoGet("/rest/api/2/search?"+params.Encode(), &result); err != nil {
		return nil, fmt.Errorf("searching issues: %w", err)
	}
	return &result, nil
}

// GetIssue fetches a single issue by key.
func (c *Client) GetIssue(key string, commentCount int) (*Issue, error) {
	expand := "renderedFields"
	fields := "summary,status,issuetype,priority,assignee,reporter,created,updated,labels,components,fixVersions,parent,subtasks,issuelinks,timetracking,description,resolution,comment,worklog"

	path := fmt.Sprintf("/rest/api/2/issue/%s?fields=%s&expand=%s", url.PathEscape(key), fields, expand)
	var issue Issue
	if err := c.DoGet(path, &issue); err != nil {
		return nil, fmt.Errorf("getting issue %s: %w", key, err)
	}
	return &issue, nil
}

// CreateIssue creates a new issue.
func (c *Client) CreateIssue(req *CreateIssueRequest) (*CreatedIssue, error) {
	var result CreatedIssue
	if err := c.DoPost("/rest/api/2/issue", req, &result); err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}
	return &result, nil
}

// EditIssue updates an existing issue.
func (c *Client) EditIssue(key string, req *EditIssueRequest) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s", url.PathEscape(key))
	if err := c.DoPut(path, req, nil); err != nil {
		return fmt.Errorf("editing issue %s: %w", key, err)
	}
	return nil
}

// AssignIssue assigns an issue to a user.
func (c *Client) AssignIssue(key, assignee string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/assignee", url.PathEscape(key))
	body := map[string]string{"name": assignee}
	if err := c.DoPut(path, body, nil); err != nil {
		return fmt.Errorf("assigning issue %s: %w", key, err)
	}
	return nil
}

// UnassignIssue removes the assignee from an issue.
func (c *Client) UnassignIssue(key string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/assignee", url.PathEscape(key))
	body := map[string]*string{"name": nil}
	if err := c.DoPut(path, body, nil); err != nil {
		return fmt.Errorf("unassigning issue %s: %w", key, err)
	}
	return nil
}

// GetTransitions returns available transitions for an issue.
func (c *Client) GetTransitions(key string) ([]Transition, error) {
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", url.PathEscape(key))
	var result struct {
		Transitions []Transition `json:"transitions"`
	}
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("getting transitions for %s: %w", key, err)
	}
	return result.Transitions, nil
}

// TransitionIssue moves an issue to a new state.
func (c *Client) TransitionIssue(key string, req *TransitionRequest) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", url.PathEscape(key))
	if err := c.DoPost(path, req, nil); err != nil {
		return fmt.Errorf("transitioning issue %s: %w", key, err)
	}
	return nil
}

// DeleteIssue deletes an issue. If cascade is true, subtasks are deleted too.
func (c *Client) DeleteIssue(key string, cascade bool) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s", url.PathEscape(key))
	if cascade {
		path += "?deleteSubtasks=true"
	}
	if err := c.DoDelete(path); err != nil {
		return fmt.Errorf("deleting issue %s: %w", key, err)
	}
	return nil
}

// AddComment adds a comment to an issue.
func (c *Client) AddComment(key, body string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/comment", url.PathEscape(key))
	payload := map[string]string{"body": body}
	if err := c.DoPost(path, payload, nil); err != nil {
		return fmt.Errorf("adding comment to %s: %w", key, err)
	}
	return nil
}

// LinkIssues creates a link between two issues.
func (c *Client) LinkIssues(inwardKey, outwardKey, linkType string) error {
	req := &LinkRequest{
		Type:         map[string]string{"name": linkType},
		InwardIssue:  map[string]string{"key": inwardKey},
		OutwardIssue: map[string]string{"key": outwardKey},
	}
	if err := c.DoPost("/rest/api/2/issueLink", req, nil); err != nil {
		return fmt.Errorf("linking %s -> %s: %w", inwardKey, outwardKey, err)
	}
	return nil
}

// UnlinkIssues removes a link between two issues.
func (c *Client) UnlinkIssues(inwardKey, outwardKey string) error {
	// First get the issue to find the link ID
	issue, err := c.GetIssue(inwardKey, 0)
	if err != nil {
		return err
	}
	for _, link := range issue.Fields.IssueLinks {
		target := ""
		if link.OutwardIssue != nil {
			target = link.OutwardIssue.Key
		}
		if link.InwardIssue != nil {
			target = link.InwardIssue.Key
		}
		if target == outwardKey {
			path := fmt.Sprintf("/rest/api/2/issueLink/%s", link.ID)
			return c.DoDelete(path)
		}
	}
	return fmt.Errorf("no link found between %s and %s", inwardKey, outwardKey)
}

// AddWorklog adds a worklog entry to an issue.
func (c *Client) AddWorklog(key, timeSpent, comment string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/worklog", url.PathEscape(key))
	payload := map[string]string{"timeSpent": timeSpent}
	if comment != "" {
		payload["comment"] = comment
	}
	if err := c.DoPost(path, payload, nil); err != nil {
		return fmt.Errorf("adding worklog to %s: %w", key, err)
	}
	return nil
}

// CloneIssue clones an issue with optional field overrides.
func (c *Client) CloneIssue(key string, summaryOverride string, replace map[string]string) (*CreatedIssue, error) {
	issue, err := c.GetIssue(key, 0)
	if err != nil {
		return nil, err
	}

	summary := issue.Fields.Summary
	description := issue.Fields.Description

	if summaryOverride != "" {
		summary = summaryOverride
	}

	for find, repl := range replace {
		summary = strings.ReplaceAll(summary, find, repl)
		description = strings.ReplaceAll(description, find, repl)
	}

	req := &CreateIssueRequest{
		Fields: CreateIssueFields{
			Project:   map[string]string{"key": strings.Split(key, "-")[0]},
			IssueType: map[string]string{"name": issue.Fields.IssueType.Name},
			Summary:   summary,
		},
	}
	if description != "" {
		req.Fields.Description = description
	}
	if issue.Fields.Priority.Name != "" {
		req.Fields.Priority = map[string]string{"name": issue.Fields.Priority.Name}
	}
	if len(issue.Fields.Labels) > 0 {
		req.Fields.Labels = issue.Fields.Labels
	}

	return c.CreateIssue(req)
}

// ListBoards lists boards for a project.
func (c *Client) ListBoards(projectKey string) ([]Board, error) {
	path := "/rest/agile/1.0/board"
	if projectKey != "" {
		path += "?projectKeyOrId=" + url.QueryEscape(projectKey)
	}
	var result struct {
		Values []Board `json:"values"`
	}
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("listing boards: %w", err)
	}
	return result.Values, nil
}

// ListSprints lists sprints for a board.
func (c *Client) ListSprints(boardID int, state string) ([]Sprint, error) {
	path := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint", boardID)
	if state != "" {
		path += "?state=" + url.QueryEscape(state)
	}
	var result struct {
		Values []Sprint `json:"values"`
	}
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("listing sprints: %w", err)
	}
	return result.Values, nil
}

// GetSprintIssues lists issues in a sprint.
func (c *Client) GetSprintIssues(sprintID int) (*SearchResult, error) {
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue?maxResults=200", sprintID)
	var result SearchResult
	if err := c.DoGet(path, &result); err != nil {
		return nil, fmt.Errorf("getting sprint issues: %w", err)
	}
	return &result, nil
}

// AddIssuesToSprint moves issues into a sprint.
func (c *Client) AddIssuesToSprint(sprintID int, issueKeys []string) error {
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID)
	body := map[string][]string{"issues": issueKeys}
	if err := c.DoPost(path, body, nil); err != nil {
		return fmt.Errorf("adding issues to sprint %d: %w", sprintID, err)
	}
	return nil
}

// AddIssuesToEpic adds issues to an epic.
func (c *Client) AddIssuesToEpic(epicKey string, issueKeys []string) error {
	// Use the standard issue update to set the parent/epic link
	for _, key := range issueKeys {
		req := &EditIssueRequest{
			Fields: map[string]any{
				"parent": map[string]string{"key": epicKey},
			},
		}
		if err := c.EditIssue(key, req); err != nil {
			return fmt.Errorf("adding %s to epic %s: %w", key, epicKey, err)
		}
	}
	return nil
}

// RemoveIssuesFromEpic removes the epic/parent link from issues.
func (c *Client) RemoveIssuesFromEpic(issueKeys []string) error {
	for _, key := range issueKeys {
		req := &EditIssueRequest{
			Fields: map[string]any{
				"parent": nil,
			},
		}
		if err := c.EditIssue(key, req); err != nil {
			return fmt.Errorf("removing %s from epic: %w", key, err)
		}
	}
	return nil
}

// ListProjects lists all accessible projects.
func (c *Client) ListProjects() ([]Project, error) {
	var projects []Project
	if err := c.DoGet("/rest/api/2/project", &projects); err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	return projects, nil
}

// ListVersions lists versions for a project.
func (c *Client) ListVersions(projectKey string) ([]Version, error) {
	path := fmt.Sprintf("/rest/api/2/project/%s/versions", url.PathEscape(projectKey))
	var versions []Version
	if err := c.DoGet(path, &versions); err != nil {
		return nil, fmt.Errorf("listing versions for %s: %w", projectKey, err)
	}
	return versions, nil
}

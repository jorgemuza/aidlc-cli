package jira

import (
	"encoding/json"
	"strings"
)

// SearchResult is the response from the Jira search API.
// Server uses startAt/maxResults/total; Cloud v3 search/jql uses nextPageToken/isLast.
type SearchResult struct {
	StartAt       int     `json:"startAt"`
	MaxResults    int     `json:"maxResults"`
	Total         int     `json:"total"`
	Issues        []Issue `json:"issues"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
	IsLast        bool    `json:"isLast,omitempty"`
}

// Issue represents a Jira issue.
type Issue struct {
	Key    string      `json:"key"`
	Self   string      `json:"self"`
	Fields IssueFields `json:"fields"`
}

// IssueFields contains the fields of a Jira issue.
type IssueFields struct {
	Summary        string          `json:"summary"`
	Description    string          `json:"-"`
	RawDescription json.RawMessage `json:"description"`
	Status      NameField  `json:"status"`
	IssueType   NameField  `json:"issuetype"`
	Priority    NameField  `json:"priority"`
	Resolution  *NameField `json:"resolution"`
	Assignee    *UserField `json:"assignee"`
	Reporter    *UserField `json:"reporter"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Labels      []string   `json:"labels"`
	Components  []NameField `json:"components"`
	FixVersions []NameField `json:"fixVersions"`
	Parent      *ParentField `json:"parent"`
	Comment     *CommentList `json:"comment"`
	Worklog     *WorklogList `json:"worklog"`
	IssueLinks  []IssueLink  `json:"issuelinks"`
	Subtasks    []Issue      `json:"subtasks"`
	Epic        *Issue       `json:"epic"`
	Sprint      *SprintField `json:"sprint"`
	TimeTracking *TimeTracking `json:"timetracking"`
}

// ResolveDescription populates the Description string from RawDescription.
// Handles both plain text (Server API v2) and ADF (Cloud API v3).
func (f *IssueFields) ResolveDescription() {
	if f.RawDescription == nil {
		return
	}
	// Try plain string first (Server)
	var s string
	if err := json.Unmarshal(f.RawDescription, &s); err == nil {
		f.Description = s
		return
	}
	// Try ADF object (Cloud) — extract text nodes
	var doc map[string]any
	if err := json.Unmarshal(f.RawDescription, &doc); err == nil {
		f.Description = extractADFText(doc)
	}
}

// extractADFText recursively extracts plain text from an ADF document.
func extractADFText(node map[string]any) string {
	if node["type"] == "text" {
		if text, ok := node["text"].(string); ok {
			return text
		}
	}
	content, ok := node["content"].([]any)
	if !ok {
		return ""
	}
	var parts []string
	for _, item := range content {
		if child, ok := item.(map[string]any); ok {
			text := extractADFText(child)
			if text != "" {
				parts = append(parts, text)
			}
		}
	}
	sep := ""
	if node["type"] == "doc" || node["type"] == "paragraph" || node["type"] == "bulletList" || node["type"] == "orderedList" || node["type"] == "listItem" {
		sep = "\n"
	}
	return strings.Join(parts, sep)
}

// NameField is a generic field with a name.
type NameField struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

// UserField represents a Jira user.
type UserField struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Name         string `json:"name"`
	AccountID    string `json:"accountId,omitempty"`
}

// ParentField represents a parent issue reference.
type ParentField struct {
	Key    string      `json:"key"`
	Fields *IssueFields `json:"fields,omitempty"`
}

// CommentList is the comment container in an issue.
type CommentList struct {
	Total    int       `json:"total"`
	Comments []Comment `json:"comments"`
}

// Comment represents a single comment.
type Comment struct {
	ID      string     `json:"id"`
	Body    string     `json:"body"`
	Author  *UserField `json:"author"`
	Created string     `json:"created"`
	Updated string     `json:"updated"`
}

// WorklogList is the worklog container in an issue.
type WorklogList struct {
	Total    int       `json:"total"`
	Worklogs []Worklog `json:"worklogs"`
}

// Worklog represents a single worklog entry.
type Worklog struct {
	ID               string     `json:"id"`
	TimeSpent        string     `json:"timeSpent"`
	TimeSpentSeconds int        `json:"timeSpentSeconds"`
	Comment          string     `json:"comment"`
	Author           *UserField `json:"author"`
	Started          string     `json:"started"`
}

// IssueLink represents a link between two issues.
type IssueLink struct {
	ID           string     `json:"id"`
	Type         LinkType   `json:"type"`
	InwardIssue  *Issue     `json:"inwardIssue,omitempty"`
	OutwardIssue *Issue     `json:"outwardIssue,omitempty"`
}

// LinkType describes the link type.
type LinkType struct {
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
}

// Transition represents a workflow transition.
type Transition struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	To   NameField `json:"to"`
}

// SprintField represents sprint info on an issue.
type SprintField struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// Sprint represents a sprint from the Agile API.
type Sprint struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Goal      string `json:"goal,omitempty"`
}

// Board represents a Jira board.
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Project represents a Jira project.
type Project struct {
	Key  string    `json:"key"`
	Name string    `json:"name"`
	Lead *UserField `json:"lead,omitempty"`
}

// Version represents a project version/release.
type Version struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}

// TimeTracking holds time tracking info.
type TimeTracking struct {
	OriginalEstimate  string `json:"originalEstimate,omitempty"`
	RemainingEstimate string `json:"remainingEstimate,omitempty"`
	TimeSpent         string `json:"timeSpent,omitempty"`
}

// CreateIssueRequest is the payload for creating an issue.
// It marshals Fields as a flat map, merging structured fields with custom fields.
type CreateIssueRequest struct {
	Fields    CreateIssueFields
	CloudMode bool // When true, description is wrapped in ADF for Cloud API v3
}

// textToADF converts plain text to Atlassian Document Format (ADF).
func textToADF(text string) map[string]any {
	return map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []map[string]any{
			{
				"type": "paragraph",
				"content": []map[string]any{
					{"type": "text", "text": text},
				},
			},
		},
	}
}

// MarshalJSON produces a flat {"fields": {...}} with custom fields merged in.
func (r CreateIssueRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["project"] = r.Fields.Project
	m["issuetype"] = r.Fields.IssueType
	m["summary"] = r.Fields.Summary
	if r.Fields.Description != "" {
		if r.CloudMode {
			m["description"] = textToADF(r.Fields.Description)
		} else {
			m["description"] = r.Fields.Description
		}
	}
	if len(r.Fields.Priority) > 0 {
		m["priority"] = r.Fields.Priority
	}
	if len(r.Fields.Assignee) > 0 {
		m["assignee"] = r.Fields.Assignee
	}
	if len(r.Fields.Reporter) > 0 {
		m["reporter"] = r.Fields.Reporter
	}
	if len(r.Fields.Labels) > 0 {
		m["labels"] = r.Fields.Labels
	}
	if len(r.Fields.Components) > 0 {
		m["components"] = r.Fields.Components
	}
	if len(r.Fields.FixVersions) > 0 {
		m["fixVersions"] = r.Fields.FixVersions
	}
	if len(r.Fields.Parent) > 0 {
		m["parent"] = r.Fields.Parent
	}
	if r.Fields.TimeTracking != nil {
		m["timetracking"] = r.Fields.TimeTracking
	}
	for k, v := range r.Fields.CustomFields {
		m[k] = v
	}
	return json.Marshal(map[string]any{"fields": m})
}

// CreateIssueFields holds fields for issue creation.
type CreateIssueFields struct {
	Project      map[string]string   `json:"project"`
	IssueType    map[string]string   `json:"issuetype"`
	Summary      string              `json:"summary"`
	Description  string              `json:"description,omitempty"`
	Priority     map[string]string   `json:"priority,omitempty"`
	Assignee     map[string]string   `json:"assignee,omitempty"`
	Reporter     map[string]string   `json:"reporter,omitempty"`
	Labels       []string            `json:"labels,omitempty"`
	Components   []map[string]string `json:"components,omitempty"`
	FixVersions  []map[string]string `json:"fixVersions,omitempty"`
	Parent       map[string]string   `json:"parent,omitempty"`
	TimeTracking *TimeTracking       `json:"timetracking,omitempty"`
	CustomFields map[string]any      `json:"-"`
}

// EditIssueRequest is the payload for editing an issue.
type EditIssueRequest struct {
	Fields map[string]any `json:"fields,omitempty"`
	Update map[string]any `json:"update,omitempty"`
}

// CreatedIssue is the response from creating an issue.
type CreatedIssue struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// TransitionRequest is the payload for transitioning an issue.
type TransitionRequest struct {
	Transition map[string]string `json:"transition"`
	Update     map[string]any    `json:"update,omitempty"`
	Fields     map[string]any    `json:"fields,omitempty"`
}

// LinkRequest is the payload for linking issues.
type LinkRequest struct {
	Type         map[string]string `json:"type"`
	InwardIssue  map[string]string `json:"inwardIssue"`
	OutwardIssue map[string]string `json:"outwardIssue"`
}

// Field represents a Jira field definition.
type Field struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Custom      bool              `json:"custom"`
	Schema      *FieldSchema      `json:"schema,omitempty"`
}

// FieldSchema describes a field's type information.
type FieldSchema struct {
	Type     string `json:"type"`
	Custom   string `json:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty"`
	System   string `json:"system,omitempty"`
}

// CreateFieldRequest is the payload for creating a custom field (Cloud only).
type CreateFieldRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	SearcherKey string `json:"searcherKey,omitempty"`
}

// CreatedField is the response from creating a custom field.
type CreatedField struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      *FieldSchema `json:"schema,omitempty"`
}

// FieldContext represents a field configuration context.
type FieldContext struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsGlobal    bool   `json:"isGlobalContext"`
	IsAnyIssue  bool   `json:"isAnyIssueType"`
}

// FieldOption represents an option for a select/multi-select field.
type FieldOption struct {
	ID       string `json:"id"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

// Status represents a Jira workflow status.
type Status struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	StatusCategory *StatusCategory `json:"statusCategory,omitempty"`
}

// StatusCategory groups statuses into categories (To Do, In Progress, Done).
type StatusCategory struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	ColorName string `json:"colorName"`
}

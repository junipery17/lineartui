package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/shurcooL/graphql"
)

type Client interface {
	GetTeams(ctx context.Context) ([]TeamData, error)
	DisplayTeams(ctx context.Context) error
	GetTeamIssues(ctx context.Context, teamID string) (*TeamData, error)
	DisplayIssues(ctx context.Context, teamID string) error
	AddIssue(ctx context.Context, teamID string, title string, description ...string) (*IssueData, error)
	DeleteIssue(ctx context.Context, issueID string) error
	UpdateAssigneeOnIssue(ctx context.Context, issueID string, assignee string) error
}

type client struct {
	gql *graphql.Client
}

func NewClient(apiKey string, apiURL string) Client {
	httpClient := &http.Client{
		Transport: &authTransport{
			token: apiKey,
			base:  http.DefaultTransport,
		},
	}

	gqlClient := graphql.NewClient(apiURL, httpClient)

	return &client{
		gql: gqlClient,
	}
}

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	return t.base.RoundTrip(req)
}

type TeamData struct {
	ID     graphql.String
	Name   graphql.String
	Issues struct {
		Nodes []IssueData
	} `graphql:"issues"`
}

type IssueData struct {
	ID          graphql.String
	Title       graphql.String
	Description graphql.String
	Assignee    struct {
		ID   graphql.String
		Name graphql.String
	} `graphql:"assignee"`
}

func (c *client) GetTeams(ctx context.Context) ([]TeamData, error) {
	var query struct {
		Teams struct {
			Nodes []TeamData
		}
	}

	err := c.gql.Query(ctx, &query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams: %w", err)
	}

	return query.Teams.Nodes, nil
}

func (c *client) DisplayTeams(ctx context.Context) error {
	teams, err := c.GetTeams(ctx)
	if err != nil {
		return err
	}

	for _, team := range teams {
		fmt.Printf("Team: ID=%s, Name=%s\n", team.ID, team.Name)
	}

	return nil
}

func (c *client) GetTeamIssues(ctx context.Context, teamID string) (*TeamData, error) {
	var query struct {
		Team TeamData `graphql:"team(id: $teamId)"`
	}

	variables := map[string]any{
		"teamId": graphql.String(teamID),
	}

	err := c.gql.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch team issues: %w", err)
	}

	return &query.Team, nil
}

func (c *client) DisplayIssues(ctx context.Context, teamID string) error {
	team, err := c.GetTeamIssues(ctx, teamID)
	if err != nil {
		return err
	}

	fmt.Printf("Team: %s\n", team.Name)
	for _, issue := range team.Issues.Nodes {
		fmt.Printf("Issue: ID=%s, Title=%s\n", issue.ID, issue.Title)
		if issue.Description != "" {
			fmt.Printf("  Description: %s\n", issue.Description)
		}
		if issue.Assignee.Name != "" {
			fmt.Printf("  Assignee: %s\n", issue.Assignee.Name)
		}
	}

	return nil
}

func (c *client) AddIssue(ctx context.Context, teamID string, title string, description ...string) (*IssueData, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	var mutation struct {
		IssueCreate struct {
			Success graphql.Boolean `graphql:"success"`
			Issue   IssueData       `graphql:"issue"`
		} `graphql:"issueCreate(input: $input)"`
	}
	type IssueCreateInput struct {
		Title       graphql.String `json:"title"`
		Description graphql.String `json:"description"`
		TeamID      graphql.String `json:"teamId"`
	}

	variables := map[string]any{
		"input": IssueCreateInput{
			Title:       graphql.String(title),
			Description: graphql.String(desc),
			TeamID:      graphql.String(teamID),
		},
	}

	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	if !mutation.IssueCreate.Success {
		return nil, errors.New("issue creation was not successful")
	}

	fmt.Printf("Created issue: %s (ID: %s)\n", mutation.IssueCreate.Issue.Title, mutation.IssueCreate.Issue.ID)
	return &mutation.IssueCreate.Issue, nil
}

func (c *client) DeleteIssue(ctx context.Context, issueID string) error {
	var mutation struct {
		IssueDelete struct {
			Success graphql.Boolean
		} `graphql:"issueDelete(id: $id)"`
	}

	variables := map[string]any{
		"id": graphql.String(issueID),
	}

	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}

	if !mutation.IssueDelete.Success {
		return errors.New("issue deletion was not successful")
	}

	fmt.Printf("Successfully deleted issue: %s\n", issueID)
	return nil
}

func (c *client) UpdateAssigneeOnIssue(ctx context.Context, issueID string, assign string) error {
	var mutation struct {
		IssueUpdate struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueUpdate(id: $issueUpdateId, input: $input)"`
	}

	type IssueUpdateInput struct {
		AssigneeId graphql.String `json:"assigneeId"`
	}

	variables := map[string]any{
		"issueUpdateId": graphql.String(issueID),
		"input": IssueUpdateInput{
			AssigneeId: graphql.String(assign),
		},
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update issue assignee: %w", err)
	}
	if !mutation.IssueUpdate.Success {
		return errors.New("issue assignee update was not successful")
	}
	fmt.Printf("Successfully updated assignee %s to issue %s\n", assign, issueID)
	return nil
}

//a915ed54-6b89-4fb4-8361-5530ebe5783d <-- id of that one test issue u made

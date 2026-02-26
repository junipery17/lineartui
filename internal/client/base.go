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
	DisplayIssues(ctx context.Context, teamID string, titlesOnly bool) error
	FindIssueByTitle(ctx context.Context, teamID string, title string) (string, error)
	FindTeamByName(ctx context.Context, name string) (string, error)
	AddIssue(ctx context.Context, teamID string, title string, description ...string) (*IssueData, error)
	DeleteIssue(ctx context.Context, issueID string) error
	UpdateAssigneeOnIssue(ctx context.Context, issueID string, assignee string) error
	UpdateDescriptionOnIssue(ctx context.Context, issueID string, description string) error
	UpdatePriorityOnIssue(ctx context.Context, issueID string, priority float64) error
	UpdateStatusOnIssue(ctx context.Context, issueID string, status string) error
	SearchLabel(ctx context.Context, labelName string) (string, error)
	CreateNewLabel(ctx context.Context, labelName string) (string, error)
	AddLabeltoIssue(ctx context.Context, issueID string, labelName string) error
	RemoveLabelFromIssue(ctx context.Context, issueID string, labelName string) error
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

func (c *client) DisplayIssues(ctx context.Context, teamID string, titlesOnly bool) error {
	team, err := c.GetTeamIssues(ctx, teamID)
	if err != nil {
		return err
	}

	fmt.Printf("Team: %s\n", team.Name)
	if titlesOnly {
		for i, issue := range team.Issues.Nodes {
			fmt.Printf("%d: %s\n", i+1, issue.Title)
		}
	} else {
		for _, issue := range team.Issues.Nodes {
			fmt.Printf("Issue: ID=%s, Title=%s\n", issue.ID, issue.Title)
			if issue.Description != "" {
				fmt.Printf("  Description: %s\n", issue.Description)
			}
			if issue.Assignee.Name != "" {
				fmt.Printf("  Assignee: %s\n", issue.Assignee.Name)
			}
		}
	}

	return nil
}

func (c *client) FindIssueByTitle(ctx context.Context, teamID string, title string) (string, error) {
	var query struct {
		Issues struct {
			Nodes []struct {
				ID    graphql.String
				Title graphql.String
			}
		} `graphql:"issues(filter: {title: {containsIgnoreCase: $title}})"`
	}

	variables := map[string]any{
		"title": graphql.String(title),
	}

	err := c.gql.Query(ctx, &query, variables)
	if err != nil {
		return "", fmt.Errorf("Unable to find issue by title: %w\n", err)
	}
	issues := query.Issues.Nodes
	if len(issues) != 1 {
		return "", errors.New("Couldn't find one exact issue by title")
	}
	return string(issues[0].ID), nil
}

func (c *client) FindTeamByName(ctx context.Context, name string) (string, error) {
	var query struct {
		Teams struct {
			Nodes []struct {
				ID   graphql.String
				Name graphql.String
			}
		} `graphql:"teams(filter: {name: {containsIgnoreCase: $name}})"`
	}

	variables := map[string]any{
		"name": graphql.String(name),
	}
	err := c.gql.Query(ctx, &query, variables)
	if err != nil {
		return "", fmt.Errorf("Could not find Team by Name: %w\n", err)
	}
	names := query.Teams.Nodes
	if len(names) != 1 {
		return "", errors.New("Couldn't find one exact Team by name")
	}
	return string(names[0].ID), nil
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

func (c *client) UpdateDescriptionOnIssue(ctx context.Context, issueID string, description string) error {
	var mutation struct {
		IssueUpdate struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueUpdate(id: $issueUpdateId, input: $input)"`
	}
	type IssueUpdateInput struct {
		Description graphql.String `json:"description"`
	}
	variables := map[string]any{
		"issueUpdateId": graphql.String(issueID),
		"input": IssueUpdateInput{
			Description: graphql.String(description),
		},
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update issue description: %w", err)
	}
	if !mutation.IssueUpdate.Success {
		return errors.New("issue description update was not successful")
	}
	fmt.Printf("Successfully updated description to issue %s\n", issueID)
	return nil
}

func (c *client) UpdatePriorityOnIssue(ctx context.Context, issueID string, priority float64) error {
	var mutation struct {
		IssueUpdate struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueUpdate(id: $issueUpdateId, input: $input)"`
	}
	type IssueUpdateInput struct {
		Priority graphql.Float `json:"priority"`
	}
	variables := map[string]any{
		"issueUpdateId": graphql.String(issueID),
		"input": IssueUpdateInput{
			Priority: graphql.Float(priority),
		},
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update issue priority: %w\n", err)
	}
	if !mutation.IssueUpdate.Success {
		return errors.New("issue priority update was not successful")
	}
	fmt.Printf("Successfully updated priority %d to issue %s\n", int(priority), issueID)
	return nil
}

func (c *client) UpdateStatusOnIssue(ctx context.Context, issueID string, status string) error {
	var mutation struct {
		IssueUpdate struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueUpdate(id : $issueUpdateId, input : $input)"`
	}
	type IssueUpdateInput struct {
		StateId graphql.String `json:"stateId"`
	}
	variables := map[string]any{
		"issueUpdateId": graphql.String(issueID),
		"input": IssueUpdateInput{
			StateId: graphql.String(status),
		},
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("Could not update status:%w\n", err)
	}
	if !mutation.IssueUpdate.Success {
		return errors.New("Issue Status not updated\n")
	}
	fmt.Print("Successfully updated status\n")
	return nil
}

func (c *client) SearchLabel(ctx context.Context, labelName string) (string, error) {
	var query struct {
		IssueLabels struct {
			Nodes []struct {
				ID graphql.String `json:"id"`
			}
		} `graphql:"issueLabels(filter: {name: {eqIgnoreCase: $name}})"`
	}
	variables := map[string]any{
		"name": graphql.String(labelName),
	}
	err := c.gql.Query(ctx, &query, variables)
	if err != nil {
		return "", fmt.Errorf("Could not find label properly: %w\n", err)
	}
	if len(query.IssueLabels.Nodes) > 0 {
		return string(query.IssueLabels.Nodes[0].ID), nil
	}
	return "", nil
}

func (c *client) CreateNewLabel(ctx context.Context, labelName string) (string, error) {
	var mutation struct {
		IssueLabelCreate struct {
			Success    graphql.Boolean `json:"success"`
			IssueLabel struct {
				ID graphql.String `json:"id"`
			}
		} `graphql:"issueLabelCreate(input: $input)"`
	}
	type IssueLabelCreateInput struct {
		Name graphql.String `json:"name"`
	}
	variables := map[string]any{
		"input": IssueLabelCreateInput{
			Name: graphql.String(labelName),
		},
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return "", fmt.Errorf("unable to create new label: %w\n", err)
	}
	if !mutation.IssueLabelCreate.Success {
		return "", errors.New("Label not created :(\n")
	}
	return string(mutation.IssueLabelCreate.IssueLabel.ID), nil
}

func (c *client) AddLabeltoIssue(ctx context.Context, issueID string, labelName string) error {
	label, _ := c.SearchLabel(ctx, labelName)
	if label == "" {
		var err error
		label, err = c.CreateNewLabel(ctx, labelName)
		if err != nil {
			return err
		}
	}
	var mutation struct {
		IssueAddLabel struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueAddLabel(id: $issueAddLabelId, labelId: $labelId)"`
	}
	variables := map[string]any{
		"issueAddLabelId": graphql.String(issueID),
		"labelId":         graphql.String(label),
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to add label to issue: %w\n", err)
	}
	fmt.Printf("successfully added '%s' to issue!\n", labelName)
	return nil
}

func (c *client) RemoveLabelFromIssue(ctx context.Context, issueID string, labelName string) error {
	label, _ := c.SearchLabel(ctx, labelName)
	if label == "" {
		return errors.New("A label with this name does not exist\n")
	}
	var mutation struct {
		IssueRemoveLabel struct {
			Success graphql.Boolean `graphql:"success"`
		} `graphql:"issueRemoveLabel(id: $issueRemoveLabelId, labelId: $labelId)"`
	}
	variables := map[string]any{
		"issueRemoveLabelId": graphql.String(issueID),
		"labelId":            graphql.String(label),
	}
	err := c.gql.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to remove label: %w\n", err)
	}
	fmt.Printf("successfully removed '%s' from issue!\n", labelName)
	return nil
}

//a915ed54-6b89-4fb4-8361-5530ebe5783d <-- id of that one test issue u made

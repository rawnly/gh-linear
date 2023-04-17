package linear

import (
	"context"
	"fmt"
	"os"

	graphql "github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

func NewGraphqlClient() *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("LINEAR_API_KEY"),
		},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	return graphql.NewClient("https://api.linear.app/graphql", httpClient)
}

type LinearClient struct {
	client *graphql.Client
}

func NewClient() *LinearClient {
	return &LinearClient{
		client: NewGraphqlClient(),
	}
}

func (c *LinearClient) Query(query interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	err := c.client.Query(context.Background(), query, variables, options...)

	if err != nil {
		fmt.Println(err)
	}

	return err
}

type MeQuery struct {
	Viewer struct {
		Id    string
		Name  string
		Email string
	}
}

func (c *LinearClient) GetMe() (*MeQuery, error) {
	var query MeQuery

	err := c.Query(&query, nil)

	return &query, err
}

type IssueQuery struct {
	Issue struct {
		Id          string
		Title       string
		Description string
		BranchName  string
	} `grapqhl:"issue(id: $id)"`
}

func (c *LinearClient) GetIssue(issueId string) (*IssueQuery, error) {
	var query IssueQuery

	err := c.client.Query(context.Background(), &query, map[string]interface{}{
		"id": graphql.String(issueId),
	})

	return &query, err
}

type Issue struct {
	Id         string
	Identifier string
	Title      string
	BranchName string

	State struct {
		Name string
		Type string
	}
}

type TeamIssues struct {
	Team struct {
		Id   string
		Name string

		Issues struct {
			Nodes []Issue
		} `graphql:"issues(filter: $filter)"`
	} `graphql:"team(id: $teamId)"`
}

type StateFilter struct {
	Type map[string]string `json:"type"`
}

type IssueFilter struct {
	State StateFilter   `json:"state"`
	And   []IssueFilter `json:"and"`
	Or    []IssueFilter `json:"or"`
}

func (c *LinearClient) GetIssues(teamId string) (*TeamIssues, error) {
	var query TeamIssues

	variables := map[string]interface{}{
		"teamId": teamId,
		"filter": IssueFilter{
			And: []IssueFilter{
				{
					State: StateFilter{
						Type: map[string]string{
							"neq": "canceled",
						},
					},
				},
				{
					State: StateFilter{
						Type: map[string]string{
							"neq": "completed",
						},
					},
				},
			},
		},
	}

	err := c.client.Query(context.Background(), &query, variables)

	return &query, err
}

type Team struct {
	Name string
	Id   string

	ActiveCycle struct {
		StartsAt string
		EndsAt   string
		Name     string
	}
}

type Teams struct {
	Teams struct {
		Nodes []Team
	}
}

func (c *LinearClient) GetTeams() (*Teams, error) {
	var query Teams

	err := c.Query(&query, nil)

	return &query, err
}

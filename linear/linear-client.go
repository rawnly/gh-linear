package linear

import (
	"context"
	"fmt"
	"net/http"
	"os"

	graphql "github.com/hasura/go-graphql-client"
)

// @see https://stackoverflow.com/questions/54088660/add-headers-for-each-http-request-using-client
type transport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	base := t.base

	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(req)
}

func NewGraphqlClient(key string) *graphql.Client {
	client := &http.Client{
		Transport: &transport{
			headers: map[string]string{
				"Authorization": key,
			},
		},
	}

	return graphql.NewClient("https://api.linear.app/graphql", client)
}

type LinearClient struct {
	client *graphql.Client
}

func NewClient() *LinearClient {
	return &LinearClient{
		client: NewGraphqlClient(os.Getenv("LINEAR_API_KEY")),
	}
}

func (c *LinearClient) SetKey(key string) {
	c.client = NewGraphqlClient(key)
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

// GetMe returns the current user
func (c *LinearClient) GetMe() (*MeQuery, error) {
	var query MeQuery

	err := c.Query(&query, nil)

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

func (i *Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Identifier, i.Title)
}

type IssueQuery struct {
	Issue struct {
		Id string
	} `graphql:"issue(id: $issueId)"`
}

// GetIssue returns the issue with the given id
func (c *LinearClient) GetIssue(issueId string) (*Issue, error) {
	var query struct {
		Issue Issue `graphql:"issue(id: $issueId)"`
	}

	variables := map[string]interface{}{
		"issueId": issueId,
	}

	err := c.client.Query(context.Background(), &query, variables)

	return &query.Issue, err
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

// GetIssues returns all issues for the given team
func (c *LinearClient) GetIssues(teamId string) (*[]Issue, error) {
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

	return &query.Team.Issues.Nodes, err
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

var teamCacheTTL int64 = 60 * 12

// GetTeams returns all teams assciated to the  current user
func (c *LinearClient) GetTeams() (*Teams, error) {
	var query Teams

	err := c.client.Query(context.Background(), &query, nil)

	return &query, err
}

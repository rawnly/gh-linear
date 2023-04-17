package main

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Rawnly/gh-linear/linear"
	"github.com/rawnly/gitgud/run"
)

func main() {
	linearClient := linear.NewClient()

	indicator := NewSpinner("Loading user...")
	indicator.Spinner.Start()

	me, err := linearClient.GetMe()

	if err != nil {
		indicator.Fail("Error loading user.")
		return
	}

	indicator.Succeed("Welcome " + me.Viewer.Name + "!")

	indicator = NewSpinner("Loading teams...")
	indicator.Spinner.Start()

	teams, err := linearClient.GetTeams()

	if err != nil {
		indicator.Fail("Error loading teams.")
		return
	}

	teamsCount := len(teams.Teams.Nodes)
	indicator.Succeed("You are a member of " + fmt.Sprint(teamsCount) + " teams.")

	var teamNames []string

	for _, team := range teams.Teams.Nodes {
		teamNames = append(teamNames, team.Name)
	}

	qs := []*survey.Question{
		{
			Name: "team",
			Prompt: &survey.Select{
				Message: "Select a team:",
				Options: teamNames,
			},
		},
	}

	answer := struct {
		Team string `survey:"team"`
	}{}

	err = survey.Ask(qs, &answer)

	var selectedTeam linear.Team

	for _, team := range teams.Teams.Nodes {
		if team.Name == answer.Team {
			selectedTeam = team
		}
	}

	if selectedTeam.Id == "" {
		fmt.Println("Team not found.")
		return
	}

	indicator = NewSpinner("Loading issues...")
	indicator.Spinner.Start()
	issues, err := linearClient.GetIssues(selectedTeam.Id)

	if err != nil {
		indicator.Fail("Error loading issues.")

		printJson(err.Error())
		return
	}

	issuesCount := len(issues.Team.Issues.Nodes)
	indicator.Succeed("You have " + fmt.Sprint(issuesCount) + " issues.")

	var issueNames []string

	for _, issue := range issues.Team.Issues.Nodes {
		issueNames = append(issueNames, fmt.Sprintf("[%s] %s", issue.Identifier, issue.Title))
	}

	var issueName string
	err = survey.AskOne(&survey.Select{
		Message: "Select an issue:",
		Options: issueNames,
	}, &issueName)

	if err != nil {
		panic(err.Error())
	}

	var selectedIssue linear.Issue

	for _, issue := range issues.Team.Issues.Nodes {
		if fmt.Sprintf("[%s] %s", issue.Identifier, issue.Title) == issueName {
			selectedIssue = issue
		}
	}

	if selectedIssue.Id == "" {
		fmt.Println("Issue not found.")
		return
	}

	fmt.Println()
	fmt.Printf("Creating branch: feature/%s\n", selectedIssue.BranchName)

	branch := fmt.Sprintf("feature/%s", selectedIssue.BranchName)

	err = run.Git("checkout", "-b", branch).RunInTerminal()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go

func printJson(obj interface{}) {
	prettyJSON, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}

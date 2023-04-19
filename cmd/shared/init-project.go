package shared

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Rawnly/gh-linear/config"
	"github.com/Rawnly/gh-linear/linear"
	linearSdk "github.com/Rawnly/gh-linear/linear"
	"github.com/Rawnly/gh-linear/utils"
	"github.com/Rawnly/gh-linear/utils/slice"
)

var DefaultApiKey = os.Getenv("LINEAR_API_KEY")

func InitProject(ctx context.Context, projectKey string) error {
	conf := ctx.Value("config").(*config.Config)
	linear := ctx.Value("linear").(*linear.LinearClient)

	var apiKey string
	if err := survey.AskOne(&survey.Input{
		Message: "Enter your Linear API key:",
		Help:    "You can find your API key in your Linear profile settings.",
		Default: DefaultApiKey,
	}, &apiKey); err != nil {
		return err
	}

	linear.SetKey(apiKey)

	var teams []linearSdk.Team

	teamsCache := conf.Teams[apiKey]

	if teamsCache == nil {
		teamsCache = &config.CachedTeams{
			LastFetch: -1,
			Data:      nil,
		}
	}

	if teamsCache.IsExpired() || apiKey != DefaultApiKey {
		t, err := linear.GetTeams()
		if err != nil {
			return err
		}

		teams = t.Teams.Nodes

		teamsCache.Set(teams)

		if conf.Teams == nil {
			conf.Teams = make(map[string]*config.CachedTeams)
		}

		conf.Teams[apiKey] = teamsCache
		conf.Update()

		if err = conf.Save(); err != nil {
			return err
		}
	} else {
		teams = teamsCache.Data
	}

	var team string
	if err := survey.AskOne(&survey.Select{
		Message: "Choose a team:",
		Options: slice.Map(teams, func(team linearSdk.Team) string { return team.Name }),
		Help:    "This project will be associated with the selected team.",
		VimMode: true,
	}, &team); err != nil {
		return err
	}

	var project *config.Project = &config.Project{
		WorkDir: projectKey,
		ApiKey:  apiKey,
	}

	for _, t := range teams {
		if t.Name == team {
			project.TeamID = t.Id
		}
	}

	if project.TeamID == "" {
		return utils.NewError("Invalid team selected.")
	}

	if conf.Projects == nil {
		conf.Projects = make(map[string]*config.Project)
	}

	conf.Projects[projectKey] = project
	conf.Update()

	if err := conf.Save(); err != nil {
		return err
	}

	fmt.Println("Project created successfully.")

	return nil
}

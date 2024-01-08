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

func askForKey() (string, error) {
	var apiKey string

	if err := survey.AskOne(&survey.Input{
		Message: "Enter your Linear API key:",
		Help:    "You can find your API key in your Linear profile settings.",
		Default: DefaultApiKey,
	}, &apiKey); err != nil {
		return "", err
	}

	return apiKey, nil
}

const CREATE_NEW_KEY = "Add new workspace"

func InitProject(ctx context.Context, projectKey string) error {
	apiKey := ""
	conf := ctx.Value("config").(*config.Config)
	linear := ctx.Value("linear").(*linear.LinearClient)

	existingTeams, err := conf.GetTeams()
	if err != nil {
		return err
	}

	if len(existingTeams) > 0 {
		var selectedKey string

		options := []string{}
		for _, t := range existingTeams {
			options = append(options, t.Workspace.Name)
		}

		options = append(options, CREATE_NEW_KEY)

		if err := survey.AskOne(&survey.Select{
			Message: "Choose a workspace:",
			Options: options,
		}, &selectedKey); err != nil {
			return err
		}

		if selectedKey == CREATE_NEW_KEY {
			apiKey = ""
		} else {
			if item, ok := slice.Find(existingTeams, func(t config.CachedTeams) bool { return t.Workspace.Name == selectedKey }); ok {
				apiKey = item.ApiKey

				if item.ApiKey == "" {
					conf.ResetRefresh()
				}
			}
		}
	}

	if apiKey == "" {
		apiKey, err = askForKey()

		if err != nil {
			return err
		}
	}

	linear.SetKey(apiKey)

	var teams []linearSdk.Team
	var workspace *linearSdk.Workspace

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

		w, err := linear.GetWorkspace()
		if err != nil {
			return err
		}

		teams = t.Teams.Nodes
		workspace = w

		teamsCache.Set(apiKey, teams, w)

		if conf.Teams == nil {
			conf.Teams = make(map[string]*config.CachedTeams)
		}

		conf.SetTeam(apiKey, teamsCache)
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
		Description: func(value string, idx int) string {
			return fmt.Sprintf("(%s) %s", teams[idx].Id, workspace.Name)
		},
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

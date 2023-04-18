package project

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Rawnly/gh-linear/config"
	"github.com/Rawnly/gh-linear/linear"
	"github.com/Rawnly/gh-linear/utils"
	"github.com/Rawnly/gh-linear/utils/slice"
	"github.com/spf13/cobra"
)

var changeTeam = &cobra.Command{
	Use:   "change-team",
	Short: "Change team for current project",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		conf := ctx.Value("config").(*config.Config)
		client := ctx.Value("linear").(*linear.LinearClient)
		spinner := utils.NewSpinner("Loading teams...")

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		project := conf.GetProject(wd)

		if project == nil {
			return fmt.Errorf("Project not initialized")
		}

		spinner.Start()
		client.SetKey(project.ApiKey)
		teams, err := client.GetTeams()
		spinner.Stop()

		old_team, _ := slice.Find(teams.Teams.Nodes, func(team linear.Team) bool {
			return team.Id == project.TeamID
		})

		var team string
		err = survey.AskOne(&survey.Select{
			Message: "Select team",
			Default: old_team.Name,
			VimMode: true,
			Options: slice.Map(teams.Teams.Nodes, func(team linear.Team) string {
				return team.Name
			}),
		}, &team)

		var selectedTeam *linear.Team
		for _, t := range teams.Teams.Nodes {
			if t.Name == team {
				selectedTeam = &t
				break
			}
		}

		if selectedTeam == nil {
			return fmt.Errorf("Team not found")
		}

		project.TeamID = selectedTeam.Id

		conf.SetProject(project.WorkDir, project)
		conf.Update()

		fmt.Println("Team updated! " + selectedTeam.Name)

		return conf.Save()
	},
}

func init() {}

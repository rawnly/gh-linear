package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	linearSdk "github.com/Rawnly/gh-linear/linear"
	"github.com/Rawnly/gh-linear/utils"
	"github.com/rawnly/gitgud/run"

	"github.com/Rawnly/gh-linear/config"
	"github.com/Rawnly/gh-linear/utils/git"
	"github.com/Rawnly/gh-linear/utils/slice"
	"github.com/spf13/cobra"
)

func initProject(linear *linearSdk.LinearClient, dir string) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	var teams []linearSdk.Team

	if conf.Teams.IsExpired() {
		t, err := linear.GetTeams()
		if err != nil {
			return err
		}

		teams = t.Teams.Nodes

		conf.Teams.Set(teams)
		conf.Update()

		if err = conf.Save(); err != nil {
			return err
		}
	} else {
		teams = conf.Teams.Data
	}

	var team string
	err = survey.AskOne(&survey.Select{
		Message: "Choose a team:",
		Options: slice.Map(teams, func(team linearSdk.Team) string { return team.Name }),
		Help:    "This project will be associated with the selected team.",
		VimMode: true,
	}, &team)

	if err != nil {
		return err
	}

	var project *config.Project = &config.Project{
		WorkDir: dir,
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

	conf.Projects[dir] = project
	conf.Update()

	if err := conf.Save(); err != nil {
		return err
	}

	fmt.Println("Project created successfully.")

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "gh-linear",
	Short: "gh-linear is a tool to help you create new branches from Linear issues",
	Args:  cobra.NoArgs,
	Example: heredoc.Doc(`
    $ gh linear --issue <IDENTIFIER>
    $ gh linear
  `),
	RunE: func(cmd *cobra.Command, args []string) error {
		linear := linearSdk.NewClient()
		conf, err := config.Load()
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		wd = strings.ToLower(wd)
		if err != nil {
			return err
		}

		project := conf.Projects[wd]

		if len(conf.Projects) == 0 || project == nil {
			project_setup := false

			err = survey.AskOne(&survey.Confirm{
				Message: "Do you want to setup a project in this directory?",
				Default: false,
			}, &project_setup)

			if !project_setup {
				return utils.NewError("Operation aborted.")
			}

			return initProject(linear, wd)
		}

		// read the --issue flag
		issueId, err := cmd.Flags().GetString("issue")

		prompt := "Loading issues..."

		if issueId != "" {
			prompt = "Loading issue..."
		}

		spinner := utils.NewSpinner(prompt)
		spinner.Start()

		var selectedIssue *linearSdk.Issue
		if issueId != "" {
			issue, err := linear.GetIssue(issueId)
			if err != nil {
				return err
			}

			selectedIssue = issue

			spinner.Succeed(fmt.Sprintf("Loaded issue: %s", issue.Identifier))
		} else {
			issues, err := linear.GetIssues(project.TeamID)
			if err != nil {
				return err
			}

			issueCount := 0

			if issues != nil {
				issueCount = len(*issues)
			}

			spinner.Succeed(fmt.Sprintf("Loaded %d issues", issueCount))

			var issue string
			survey.AskOne(&survey.Select{
				Message: "Choose an issue:",
				VimMode: true,
				Options: slice.Map(*issues, func(issue linearSdk.Issue) string {
					return issue.String()
				}),
				Description: func(value string, index int) string {
					return (*issues)[index].BranchName
				},
			}, &issue)

			for _, i := range *issues {
				if i.String() == issue {
					selectedIssue = &i
				}
			}
		}

		if selectedIssue == nil {
			spinner.Fail("Invalid issue selected.")
			return nil
		}

		branches, err := git.GetBranches()
		if err != nil {
			return err
		}

		branch := selectedIssue.BranchName

		if slice.Includes(branches, branch) {
			return run.Git("checkout", branch).RunInTerminal()
		}

		should_continue := false

		if err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Do you want to create the branch %s?", branch),
			Default: false,
		}, &should_continue); err != nil {
			return err
		}

		if !should_continue {
			fmt.Println("Operation aborted.")
			return nil
		}

		return run.Git("checkout", "-b", branch).RunInTerminal()
	},
}

func init() {
	rootCmd.Flags().StringP("issue", "i", "", "The issue identifier")
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

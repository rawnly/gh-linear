package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/Rawnly/gh-linear/cmd/project"
	"github.com/Rawnly/gh-linear/cmd/shared"
	linearSdk "github.com/Rawnly/gh-linear/linear"
	"github.com/Rawnly/gh-linear/utils"
	"github.com/rawnly/gitgud/run"

	"github.com/Rawnly/gh-linear/config"
	"github.com/Rawnly/gh-linear/utils/git"
	"github.com/Rawnly/gh-linear/utils/slice"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-linear",
	Short: "gh-linear is a tool to help you create new branches from Linear issues",
	Args:  cobra.NoArgs,
	Example: heredoc.Doc(`
    $ gh linear --issue <IDENTIFIER>
    $ gh linear
  `),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		linear := ctx.Value("linear").(*linearSdk.LinearClient)
		conf := ctx.Value("config").(*config.Config)

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		wd = strings.ToLower(wd)
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

			return shared.InitProject(ctx, wd)
		}

		linear.SetKey(project.ApiKey)

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
	rootCmd.AddCommand(project.Cmd)
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

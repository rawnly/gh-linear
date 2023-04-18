package project

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects"},
	Short:   "Manage projects",
	Example: heredoc.Doc(`
		$ gh-linear project init
		$ gh-linear project set-key <api-key>
		$ gh-linear project change-team
	`),
}

func init() {
	Cmd.AddCommand(initProject)
	Cmd.AddCommand(setProjectKey)
	Cmd.AddCommand(changeTeam)
}

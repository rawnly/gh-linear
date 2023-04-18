package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/Rawnly/gh-linear/config"
	"github.com/spf13/cobra"
)

var setProjectKey = &cobra.Command{
	Use:     "set-key",
	Aliases: []string{"set-api-key"},
	Short:   "Set the API key for project",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		conf := ctx.Value("config").(*config.Config)

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		wd = strings.ToLower(wd)
		project := conf.Projects[wd]

		if project == nil {
			return fmt.Errorf("Project not initialized")
		}

		project.ApiKey = args[0]
		conf.Projects[wd] = project
		conf.Update()

		fmt.Println("Project API key set successfully")

		return conf.Save()
	},
}

func init() {}

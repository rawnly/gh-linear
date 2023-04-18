package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/Rawnly/gh-linear/cmd/shared"
	"github.com/Rawnly/gh-linear/config"
	"github.com/spf13/cobra"
)

var initProject = &cobra.Command{
	Use:     "init",
	Example: "gh-init project init",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		conf := ctx.Value("config").(*config.Config)

		forceInit, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		wd = strings.ToLower(wd)
		project := conf.Projects[wd]

		if project != nil && !forceInit {
			fmt.Println("Project already initialized")
			fmt.Println("Use --force to re-init project")
			return nil
		}

		if err := shared.InitProject(ctx, wd); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	initProject.Flags().BoolP("force", "f", false, "Force init project")
}

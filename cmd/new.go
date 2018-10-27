package cmd

import (
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// NewCmd - `tok new {project}`
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Drupal 8 project",
	Long:  "Creates a new Drupal project with `tok new {project-name}`, initialises git, installs Drupal and dependencies, installs Tokaido on the new site",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		tok.New(args)
	},
}

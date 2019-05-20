package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

var templateFlag string

// NewCmd - `tok new {project}`
var NewCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Drupal 8 project",
	Long: `The fastest way to launch new Drupal projects
using open source templates by Tokaido and the community.
Complete documentation is available at https://docs.tokaido.io/tokaido/starting-a-new-drupal-project`,
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("new")
		utils.CheckCmdHard("docker-compose")

		tok.New(args, templateFlag)
	},
}

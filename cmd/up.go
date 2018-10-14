package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/ironstar-io/tokaido/constants"

	"log"

	"github.com/spf13/cobra"
)

// UpCmd - `tok up`
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start a Tokaido local development environment",
	Long:  "Start a Tokaido local development environment",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		if conf.GetConfig().Tokaido.Project.Name == constants.ProjectRootNotFound {
			log.Fatal(constants.ProjectNotFoundError)
		}

		tok.Init()
		tok.InitMessage()
	},
}

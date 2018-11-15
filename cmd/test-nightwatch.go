package cmd

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// TestNightwatchCmd - `tok test:nightwatch`
var TestNightwatchCmd = &cobra.Command{
	Use:   "test:nightwatch",
	Short: "Run Nightwatch tests",
	Long:  "Run the Nightwatch test suite",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("test")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.BackgroundServiceWarning(conf.GetConfig().Tokaido.Project.Name)

		err := docker.StatusCheck()
		if err != nil {
			log.Fatalf("Tokaido containers must be running in order to start automated tests. Have you run `tok up`?")
		}

		nightwatch.RunDrupalTests()
	},
}

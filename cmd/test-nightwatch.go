package cmd

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch"
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
		telemetry.SendCommand("test:nightwatch")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
		if !ok {
			log.Fatalf("Tokaido containers must be running in order to start automated tests. Have you run `tok up`?")
		}

		nightwatch.RunDrupalTests()
	},
}

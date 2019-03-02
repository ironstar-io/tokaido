package cmd

import (
	"log"

	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// TestCmd - `tok test`
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all available tests",
	Long:  "Runs all available test suites; Drupal with Nightwatch, others TBD",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("test")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		err := docker.StatusCheck()
		if err != nil {
			log.Fatalf("Tokaido containers must be running in order to start automated tests. Have you run `tok up`?")
		}

		// phpcs.RunLinter()
		nightwatch.RunDrupalTests()
	},
}

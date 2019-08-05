package cmd

import (
	"fmt"
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/testing/testcafe"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

var useExistingDBFlag bool

// TestCmd - `tok test`
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all available tests",
	Long:  "Runs all available test suites for Drupal with TestCafe",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("test")
		telemetry.SendCommand("test")
		fmt.Println()
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
		if !ok {
			log.Fatalf("Tokaido containers must be running in order to start automated tests. Have you run `tok up`?")
		}

		testcafe.RunDrupalTests(useExistingDBFlag)
	},
}

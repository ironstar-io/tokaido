package cmd

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/testing/phpcs"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// TestPhpcbfCmd - `tok test:phpcbf`
var TestPhpcbfCmd = &cobra.Command{
	Use:   "test:phpcbf",
	Short: "Run PHPCodeSniffer linting and attempt to fix automatically",
	Long:  "Run the PHPCodeSniffer linter and attempt to fix automatically",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("test")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
		if !ok {
			log.Fatalf("Tokaido containers must be running in order to start automated tests. Have you run `tok up`?")
		}

		phpcs.CheckReqs()
		phpcs.RunLinterFix()
	},
}

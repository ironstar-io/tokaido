package cmd

import (
	"fmt"

	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// StatusCmd - `tok status`
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Have Tokaido perform a self-test",
	Long:  "Checks the status of your Tokaido environment",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("status")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		err := docker.StatusCheck()
		if err == nil {
			fmt.Println()
			console.Println(`🙂  All containers are running`, "√")
		}

		err = ssh.CheckKey()

		err = drupal.CheckContainer()

		if err == nil {
			fmt.Println()
			console.Println(`🍜  Checks have passed successfully`, "")
			fmt.Println()
			console.Println(`🌎  Run 'tok open' to open the environment in your default browser`, "")
			fmt.Println()
		} else {
			fmt.Println()
			console.Println("🙅  Some checks failed! You might be able to fix this by running `tok rebuild`", "")
			fmt.Println()
		}
	},
}

package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/unison"
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
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		docker.StatusCheck()

		ssh.CheckKey()

		unison.CheckSyncService(conf.GetConfig().Tokaido.Project.Name)

		drupal.CheckContainer()

		console.Println(`
🍜  Checks have passed successfully
		`, "")
		console.Println(`🌎  Run 'tok open' to open the environment in your default browser
		`, "")
	},
}

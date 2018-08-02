package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// StatusCmd - `tok status`
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of running containers",
	Long:  "Checks the status of containers lifted against the projects' docker-compose.yml - `docker-compose ps`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		docker.StatusCheck()

		ssh.CheckKey()

		unison.CheckSyncService()

		drupal.CheckContainer()

		console.Println(`
🍜  Checks have passed successfully
		`, "")
		console.Println(`🌎  Run 'tok open' to open the environment at 'https://localhost:`+docker.LocalPort("haproxy", "8443")+`' in your default browser
		`, "")
	},
}

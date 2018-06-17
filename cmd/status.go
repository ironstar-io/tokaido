package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
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

		conf.LoadConfig(cmd)

		docker.Status()

		ssh.CheckKey()

		drupal.CheckContainer()
	},
}

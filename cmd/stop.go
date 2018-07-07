package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// StopCmd - `tok stop`
var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all containers",
	Long:  "Gracefully stop your containers - `docker-compose stop`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		conf.LoadConfig(cmd)

		docker.Stop()

		unison.StopSyncService()
	},
}

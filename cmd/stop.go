package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"

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

		docker.Stop()

		unison.StopSyncService()
	},
}

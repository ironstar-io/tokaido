package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// StopCmd - `tok stop`
var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Gracefully stop/pause your Tokaido environment",
	Long:  "Gracefully stop/pause your Tokaido environment. Restart with `tok up`",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig()
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		docker.Stop()

		unison.StopSyncService(conf.GetConfig().Tokaido.Project.Name)
	},
}

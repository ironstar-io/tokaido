package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// SyncCmd - `tok sync`
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Perform a one-time sync of your Tokaido environment and local host",
	Long:  "Perform a one-time sync of your Tokaido environment and local host",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.Sync(conf.GetConfig().Tokaido.Project.Name)
	},
}

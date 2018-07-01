package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// DestroyCmd - `tok destroy`
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Stop and destroy all containers",
	Long:  "Gracefully stop and destroy your Tokaido containers - `docker-compose down`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.StopSyncService()

		conf.LoadConfig(cmd)

		docker.Down()
	},
}

package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// SyncCmd - `tok sync`
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Use unison to sync your files",
	Long:  "One-time sync to your container",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.Sync()
	},
}

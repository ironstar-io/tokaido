package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// WatchCmd - `tok watch`
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch and sync your files",
	Long:  "Watch your files for changes and sync them to your container",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.Watch()
	},
}

package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// WatchCmd - `tok watch`
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch and sync your files",
	Long:  "Watch your files for changes and sync them to your container",
	Run: func(cmd *cobra.Command, args []string) {
		conf.LoadConfig(cmd)

		utils.CheckCmdHard("docker-compose")

		unison.Watch()
	},
}

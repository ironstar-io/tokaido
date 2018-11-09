package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// RebuildCmd - `tok rebuild`
var RebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuilds your Tokaido environment",
	Long:  "Rebuilds your Tokaido environmnet",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		docker.Stop()
		tok.Init()
		tok.InitMessage()

	},
}

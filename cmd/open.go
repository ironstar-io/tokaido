package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// OpenCmd - `tok open`
var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the site in your default browser",
	Long:  "Opens your default browser pointing to the Tokaido HTTPS port",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.BackgroundServiceWarning(conf.GetConfig().Tokaido.Project.Name)

		system.OpenSite(args)
	},
}

package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// LogsCmd - `tok logs [x]`
var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Output container logs to the console",
	Long:  "Output container logs to the console for all or a single container. Example: tok logs fpm",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig()
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		unison.BackgroundServiceWarning(conf.GetConfig().Tokaido.Project.Name)

		docker.PrintLogs(args)
	},
}

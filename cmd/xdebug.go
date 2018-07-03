package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/xdebug"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// XdebugCmd - `tok xdebug <port>`
var XdebugCmd = &cobra.Command{
	Use:   "xdebug",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		conf.LoadConfig(cmd)

		xdebug.Configure()
	},
}

package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// OpenCmd - `tok open`
var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the site in your default browser",
	Long:  "Opens your default browser pointing to the Tokaido HTTPS port",
	Run: func(cmd *cobra.Command, args []string) {
		var GOOS = utils.CheckOS()
		conf.LoadConfig(cmd)

		utils.CheckCmdHard("docker-compose")

		if GOOS == "linux" {
			linux.OpenSite()
		}
	
		if GOOS == "osx" {
			osx.OpenSite()
		}
	},
}

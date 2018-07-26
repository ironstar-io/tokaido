package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/tok"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// RepairCmd - `tok repair`
var RepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Compose and attempt repair of your containers",
	Long:  "Functionally similar to `tok up`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		console.Println(`
🚅  Tokaido is attempting to repair your project!`, "")

		tok.Init()
		tok.InitMessage()
	},
}

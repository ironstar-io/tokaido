package cmd

import (
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// RepairCmd - `tok repair`
var RepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Compose and attempt repair of your containers",
	Long:  "Functionally similar to `tok up`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		console.Println(`
🚅  Tokaido is attempting to repair your project!`, "")

		tok.Init()
		tok.InitMessage()
	},
}

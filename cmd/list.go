package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// ListCmd - `tok list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the status of all Tokaido projects on this system",
	Long:  "Lists the status of all Tokaido projects on this system",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("list")
		utils.CheckCmdHard("docker-compose")

		tok.List()
	},
}

package cmd

import (
	"fmt"
	
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// SyscheckCmd - `tok syscheck`
var SyscheckCmd = &cobra.Command{
	Use:   "syscheck",
	Short: "Test if your local system is ready to run Tokaido",
	Long:  "Test if your local system is ready to run Tokaido",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("syscheck")
		telemetry.SendCommand("syscheck")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		fmt.Println()
		console.Println(`🚅  Checking Drupal for compatibility with Tokaido`, "")
		fmt.Println()

		drupal.CheckLocal()

		console.Println(`🚉  Drupal compatibility checks complete!`, "")
		fmt.Println()
	},
}

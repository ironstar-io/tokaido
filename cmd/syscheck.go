package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// SyscheckCmd - `tok syscheck`
var SyscheckCmd = &cobra.Command{
	Use:   "syscheck",
	Short: "Test if your local system is ready to run Tokadio and Drupal",
	Long:  "This will check if your system runs a supported version of Docker Machine, has appropriate permissions, and has the right version of PHP and Composer to maintain your Drupal site.",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		console.Println(`
🚅  Checking Drupal for compatibility with Tokaido
		`, "")

		drupal.CheckLocal()

		console.Println(`
🚉  Drupal compatibility checks complete!
		`, "")
	},
}

package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/utils"

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

package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// SyscheckCmd - `tok syscheck`
var SyscheckCmd = &cobra.Command{
	Use:   "syscheck",
	Short: "Test if your local system is ready to run Tokadio and Drupal",
	Long:  "This will check if your system runs a supported version of Docker Machine, has appropriate permissions, and has the right version of PHP and Composer to maintain your Drupal site.",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckPathHard("docker-compose")

		fmt.Println(`
🚅  Checking system for compatibility with Tokaido
		`)

		system.Check()
	},
}

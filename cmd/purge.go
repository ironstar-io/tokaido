package cmd

import (
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// PurgeCmd - `tok up`
var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge Varnish and Drupal cache",
	Long:  "Purge will purge the Varnish cache and run a drush cache-rebuild operation",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("purge")
		telemetry.SendCommand("purge")
		utils.CheckCmdHard("docker-compose")

		drupal.Purge()
	},
}

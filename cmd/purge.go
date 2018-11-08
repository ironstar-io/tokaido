package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/spf13/cobra"
)

// PurgeCmd - `tok up`
var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge Varnish and Drupal cache",
	Long:  "Purge will purge the Varnish cache and run a drush cache-rebuild operation",
	Run: func(cmd *cobra.Command, args []string) {
		conf.ValidProjectRoot()

		unison.BackgroundServiceWarning(conf.GetConfig().Tokaido.Project.Name)

		drupal.Purge()
	},
}

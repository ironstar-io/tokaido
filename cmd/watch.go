package cmd

import (
	"fmt"

	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/spf13/cobra"
)

// WatchCmd - `tok watch`
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Maintain a foreground sync service using Unison",
	Long:  "Watch your files for changes and sync them to your Tokaido environment, until you exit.",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("watch")
		fmt.Println(`🛠   This command was deprecated in Tokaido 1.6.0.

    Tokaido no longer relies on a background sync service, so you
    should not need to use this command. If you are having problems
    with sync please come visit us on the #Tokaido channel in the
    official Drupal Slack.`)

	},
}

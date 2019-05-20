package cmd

import (
	"fmt"

	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/spf13/cobra"
)

// SyncCmd - `tok sync`
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Perform a one-time sync of your Tokaido environment and local host",
	Long:  "Perform a one-time sync of your Tokaido environment and local host",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("sync")
		fmt.Println(`
🛠   This command was deprecated in Tokaido 1.6.0.

    Tokaido no longer relies on a background sync service, so you
    should not need to use this command. If you are having problems
    with sync please come visit us on the #Tokaido channel in the
    official Drupal Slack.
`)

	},
}

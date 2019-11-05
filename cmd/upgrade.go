package cmd

import (
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/version"

	"github.com/spf13/cobra"
)

// UpgradeCmd - `tok upgrade`
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version of Tokaido",
	Long:  "Upgrade to the latest version of Tokaido",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.LoadConfig("upgrade")
		telemetry.SendCommand("upgrade")

		version.Upgrade()
	},
}

package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/version"

	"github.com/spf13/cobra"
)

// InstallCmd - `tok install`
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and use this version of Tokaido",
	Long:  "Install and use this version of Tokaido",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("install")

		version.SelfInstall(true)
	},
}

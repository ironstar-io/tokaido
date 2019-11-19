package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/version"

	"github.com/spf13/cobra"
)

// InstallCmd - `tok install`
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the running Tokaido binary to your PATH",
	Long:  "Install the running Tokaido binary to your PATH",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("install")

		version.SelfInstall(true)
	},
}

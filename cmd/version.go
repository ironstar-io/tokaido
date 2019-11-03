package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/version"

	"github.com/spf13/cobra"
)

// VersionCmd - `tok version`
var VersionCmd = &cobra.Command{
	Use:   "version [version]",
	Short: "Print Tokdaido version information",
	Long:  "Print Tokdaido version information including 'Build Date', 'Compiler' and 'Platform'",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			version.Display()

			return
		}

		telemetry.SendCommand("version select (" + args[0] + ")")

		version.Select(args[0])
	},
}

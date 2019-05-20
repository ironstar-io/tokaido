package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/version"

	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd - `tok version`
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Tokdaido version information",
	Long:  "Print Tokdaido version information including 'Build Date', 'Compiler' and 'Platform'",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("version")
		info := version.Get()

		fmt.Println(`
Tokaido Version: v` + info.Version + `
Build Date:      ` + info.BuildDate + `
Compiler:        ` + info.GoVersion + `
Platform:        ` + info.Platform + `
		`)
	},
}

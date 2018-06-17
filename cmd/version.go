package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/system/version"

	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd - `tok version`
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "TODO",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()

		fmt.Println(`
Tokaido Version: v` + info.Version + `
Build Date:      ` + info.BuildDate + `
Compiler:        ` + info.GoVersion + `
Platform:        ` + info.Platform + `
		`)
	},
}

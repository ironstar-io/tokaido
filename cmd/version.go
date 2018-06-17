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
		fmt.Printf("%#v\n", version.Get())
	},
}

package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/spf13/cobra"
)

// ConfigCmd - `tok config`
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "An interactive Tokaido config editor",
	Long:  "An interactive Tokaido config editor",
	Run: func(cmd *cobra.Command, args []string) {
		conf.ValidProjectRoot()

		conf.MainMenu()
	},
}

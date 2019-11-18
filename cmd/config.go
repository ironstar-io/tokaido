package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// ConfigCmd - `tok config`
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "An interactive Tokaido config editor",
	Long:  "An interactive Tokaido config editor",
	Run: func(cmd *cobra.Command, args []string) {
		if system.CheckOS() == "windows" {
			fmt.Println(aurora.Red("Sorry! The 'tok config' command isn't available on Windows. Please use `tok config-set` instead"))
			fmt.Println("Please see the Full Tokaido Configuration Reference at https://docs.tokaido.io/en/docs/advanced/tokaido-config-reference for further detail.")
			os.Exit(1)
		}

		initialize.TokConfig("config")
		conf.ValidProjectRoot()
		telemetry.SendCommand("config")

		conf.MainMenu()
	},
}

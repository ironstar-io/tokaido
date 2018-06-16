package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"

	"fmt"

	"github.com/spf13/cobra"
)

// IronstarCmd - `tok ironstar`
var IronstarCmd = &cobra.Command{
	Use:   "ironstar",
	Short: "Information about Ironstar",
	Long:  "TODO: Pull company information from endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		conf.LoadConfig(cmd)

		fmt.Println("TODO: Pull company information from endpoint")
	},
}

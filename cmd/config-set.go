package cmd

import (
	"github.com/ironstar-io/tokaido/conf"

	"github.com/spf13/cobra"
)

// ConfigSetCmd - `tok config-set`
var ConfigSetCmd = &cobra.Command{
	Use:   "config-set",
	Short: "Set a config property value",
	Long:  "Set a config property value at a position defined in command arguments. Eg. `tok config-set drupal path`",
	Run: func(cmd *cobra.Command, args []string) {
		conf.SetConfigValueByArgs(args)
	},
}

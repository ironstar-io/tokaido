package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/spf13/cobra"
)

// ConfigSetCmd - `tok config-set`
var ConfigSetCmd = &cobra.Command{
	Use:   "config-set",
	Short: "Set a config property value",
	Long:  "Set a config property value. Eg. `tok config-set services solr enabled true`. See https://tokaido.io/docs/config for a full list of available options",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("config-get")

		conf.SetConfigValueByArgs(args, "project")
	},
}

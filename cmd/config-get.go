package cmd

import (
	"github.com/ironstar-io/tokaido/conf"

	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// ConfigGetCmd - `tok config-get`
var ConfigGetCmd = &cobra.Command{
	Use:   "config-get",
	Short: "Get a config property value",
	Long:  "Get a config property value at a position defined in command arguments. Eg. `tok config-get drupal path`",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := conf.GetConfigValueByArgs(args)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", c)
	},
}

package cmd

import (
	"github.com/ironstar-io/tokaido/conf"
	"reflect"

	"fmt"
	"log"
	"regexp"

	"github.com/sanity-io/litter"
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

		if c.Kind() == reflect.Struct {
			d := litter.Sdump(c.Interface())
			s := regexp.MustCompile(`(?s)struct {.*?}{`).ReplaceAllString(d, "{")

			fmt.Println(s)
			return
		}

		fmt.Printf("%+v\n", c)
	},
}

package auth

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/ironstar/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShowCmd - `tok auth show`
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show available credentials",
	Long:  "Show available credentials for the Ironstar API and how they map against your projects",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarShowCredentials(args)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

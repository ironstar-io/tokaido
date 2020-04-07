package auth

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/ironstar/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SetProjectCredentialsCmd - `tok auth [email]`
var SetDefaultCredentials = &cobra.Command{
	Use:   "default [email]",
	Short: "Set default credentials to be used",
	Long:  "Set default credentials to be used",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.SetDefaultCredentials(args)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}

		fmt.Println()
		color.Green("Default credentials successfully updated!")
	},
}

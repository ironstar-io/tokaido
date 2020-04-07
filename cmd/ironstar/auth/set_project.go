package auth

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/ironstar/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SetProjectCredentialsCmd - `tok auth [email]`
var SetProjectCredentialsCmd = &cobra.Command{
	Use:   "project [email]",
	Short: "Set credentials to be used for this project",
	Long:  "Set credentials to be used for this project",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.SetProjectCredentials(args)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}

		fmt.Println()
		color.Green("Project credentials successfully updated!")
	},
}

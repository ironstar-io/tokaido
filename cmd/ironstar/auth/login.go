package auth

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/ironstar/auth"

	"github.com/spf13/cobra"
)

// LoginCmd - `tok auth login`
var LoginCmd = &cobra.Command{
	Use:   "login [email]",
	Short: "Authenticate for the Ironstar API",
	Long:  "Authenticate and store credentials for use against the Ironstar API",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarAPILogin(args)
		if err != nil {
			fmt.Println("Ironstar API login failed!")
			fmt.Println(err)

			os.Exit(1)
		}

		fmt.Println("Ironstar API login succeeded!")
	},
}

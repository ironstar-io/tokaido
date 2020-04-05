package auth

import (
	"os"

	"github.com/ironstar-io/tokaido/ironstar/auth"

	"github.com/spf13/cobra"
)

var PasswordFlag string

// LoginCmd - `tok auth login`
var LoginCmd = &cobra.Command{
	Use:   "login [email]",
	Short: "Authenticate for the Ironstar API",
	Long:  "Authenticate and store credentials for use against the Ironstar API",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarAPILogin(args, PasswordFlag)
		if err != nil {
			os.Exit(1)
		}
	},
}

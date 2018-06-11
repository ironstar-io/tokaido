package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// IronstarCmd - `tok ironstar`
var IronstarCmd = &cobra.Command{
	Use:   "ironstar",
	Short: "Information about Ironstar",
	Long:  "TODO: Pull company information from endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		confirmCreate := utils.ConfirmationPrompt("Tokaido needs to create database connection settings for your site. May we add the file 'docroot/sites/default/settings.tok.php' and reference it from 'settings.php'?")
		if confirmCreate == false {
			fmt.Println(`
No problem! Please make sure that you manually configure your Drupal site to use the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido
		`)
			return
		}
		fmt.Println("TODO: Pull company information from endpoint")
	},
}

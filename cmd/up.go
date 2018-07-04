package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/tok"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// UpCmd - `tok up`
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Compose and run your containers",
	Long:  "Runs in unison in the background - `docker-compose up -d`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		tok.Init()

		fmt.Println(`
WELCOME TO TOKAIDO
==================

Your Drupal development environment is now up and running
		`)

		fmt.Println(`⌚  Run "tok watch" to keep files in your local system and the Tokaido environment synchronised`)
		fmt.Printf("💻  Run \"ssh %s.tok\" to access the Drush container\n", conf.GetConfig().Project)
		fmt.Println(`🌎  Run "tok open" to open the environment in your browser`)
		fmt.Println(`
Check out https://docs.tokaido.io/environments for tips on managing your Tokaido environment
		`)
	},
}

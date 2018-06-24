package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/tok"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// RepairCmd - `tok repair`
var RepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Compose and attempt repair of your containers",
	Long:  "Functionally similar to `tok up`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		fmt.Println(`
🚅  Tokaido is attempting to repair your project!`)

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

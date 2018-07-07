package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/tok"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// InitCmd - `tok init`
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Tokaido project",
	Long:  "Initialize a Tokaido project. Installs dependencies, starts unison synchronization, builds a configuration file, starts your container. Docker is required to be installed manually for this to work",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		fmt.Printf("The command 'tok init' has been deprecated. Please use 'tok up' instead.\n\n")

		tok.Init()

		fmt.Println(`
WELCOME TO TOKAIDO
==================

Your Drupal development environment is now up and running
		`)

		fmt.Printf("💻  Run \"ssh %s.tok\" to access the Drush container\n", conf.GetConfig().Project)
		fmt.Println(`🌎  Run "tok open" to open the environment in your browser`)
		fmt.Println(`🤔  Run "tok status" check the status of your environment`)
		fmt.Println(`
Check out https://docs.tokaido.io/environments for tips on managing your Tokaido environment
		`)

	},
}

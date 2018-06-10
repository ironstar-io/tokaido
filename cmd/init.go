package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
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

		fmt.Println(`🚅  Tokaido is initializing your project!`)

		system.CheckDependencies()

		ssh.GenerateKeys()

		unison.DockerUp()
		unison.CreateOrUpdatePrf()
		unison.Sync()

		docker.Up()
		docker.Status()

		drupal.ConfigureSSH()
		config := conf.GetConfig()

		fmt.Println(`🚉  Tokaido successfully initialised your environment!`)
		fmt.Println(`🌎  Check out https://docs.tokaido.io/environments for tips on managing your Tokaido environment`)
		fmt.Println(`⌚  Run "tok watch" to keep files in your local system and the Tokaido environment synchronised`)
		fmt.Printf("💻  To access Drush via SSH run ssh %s.tok", config.Project)
	},
}

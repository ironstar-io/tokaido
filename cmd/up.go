package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
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

		fmt.Println(`
🚅  Tokaido is pulling up your containers!
		`)

		unison.DockerUp()
		unison.CreateOrUpdatePrf()

		docker.Up()

		drupal.ConfigureSSH()

		fmt.Println(`
🚁  Tokaido lifted containers successfully!
		`)
	},
}

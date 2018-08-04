package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// SSHCmd - `tok status`
var SSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into your Drush container",
	Long:  "SSH into your Drush container using the Tokaido generated key",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		drupal.ConfigureSSH()
	},
}

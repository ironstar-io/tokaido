package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// SSHCmd - `tok status`
var SSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into your Drush container",
	Long:  "SSH into your Drush container using the Tokaido generated key",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		drupal.ConfigureSSH()
	},
}

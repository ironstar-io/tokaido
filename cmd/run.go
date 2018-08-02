package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"github.com/spf13/cobra"
)

// RunCmd - `tok run`
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command inside your Tokaido service",
	Long:  "Aliases `docker-compose -f x exec`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.Exec(args)
	},
}

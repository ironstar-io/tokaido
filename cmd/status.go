package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	"github.com/spf13/cobra"
)

// StatusCmd - `tok status`
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of running containers",
	Long:  "Checks the status of containers lifted against the projects' docker-compose.yml - `docker-compose ps`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckPath("docker-compose")

		utils.StdoutCmd("docker-compose", "ps")
	},
}

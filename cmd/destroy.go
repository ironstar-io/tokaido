package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// DestroyCmd - `tok destroy`
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Stop and destroy all containers",
	Long:  "Gracefully stop and destroy your Tokaido containers - `docker-compose down`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		fmt.Println(`
🚅  Tokaido is pulling down your containers!
		`)

		docker.Down()

		fmt.Println(`
🚉  Tokaido destroyed containers successfully!
		`)
	},
}

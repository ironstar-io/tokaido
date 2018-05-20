package cmd

import (
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
		utils.CheckPath("docker-compose")

		fmt.Println(`
🚅  Tokaido is pulling up your containers!
		`)

		utils.StdoutCmd("docker-compose", "up", "-d")

		fmt.Println(`
🚁  Tokaido lifted containers successfully!
		`)
	},
}

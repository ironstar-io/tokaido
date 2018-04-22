package cmd

import (
	"../utils"
	"fmt"
	"github.com/spf13/cobra"
)

// DownCmd - `tok down`
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all containers",
	Long:  "Gracefully stop your containers - `docker-compose down`",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckPath("docker-compose")

		fmt.Println(`
🚅  Tokaido is pulling down your containers!
		`)

		utils.StdoutCmd("docker-compose", "down")

		fmt.Println(`
🚉  Tokaido stopped containers successfully!
		`)
	},
}

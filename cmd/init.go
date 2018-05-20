package cmd

import (
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
		utils.CheckPath("docker-compose")

		fmt.Println(`
🚅  Tokaido is initializing your project!
		`)

		utils.StdoutCmd("docker-compose", "down")

		fmt.Println(`
🚉  Tokaido successfully initialized you project!
		`)
	},
}

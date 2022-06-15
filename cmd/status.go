package cmd

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// StatusCmd - `tok status`
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Have Tokaido perform a self-test",
	Long:  "Checks the status of your Tokaido environment",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("status")
		telemetry.SendCommand("status")
		utils.CheckCmdHard("docker-compose")

		fmt.Println()

		docker.HardCheckTokCompose()

		ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
		if ok {
			console.Println(`😊  All containers are running`, "√")
		} else {
			console.Println(`😓  Tokaido containers are not running`, "×")
			fmt.Println(`    It appears that some or all of the Tokaido containers are offline.

    View the status of your containers with 'tok ps'

    You can try to fix this by running 'tok up', or by running 'tok repair'.`)
		}

		ok = ssh.CheckKey()
		if ok {
			console.Println("😀  SSH access is configured", "√")
		} else {
			console.Println("😓  SSH access is not configured", "×")
			fmt.Println("    Your SSH access to the Drush container looks broken.")
			fmt.Println("    You should be able to run 'tok repair' to attempt to fix this automatically")
		}

		err := drupal.CheckContainer()

		if err == nil {
			fmt.Println()
			console.Println(`🍜  Checks have passed successfully`, "")
			fmt.Println()
			console.Println(`🌎  Run 'tok open' to open the environment in your default browser`, "")
			fmt.Println()
		} else {
			fmt.Println()
			fmt.Println(aurora.Yellow("🙅  Some checks failed! You might be able to fix this by running `tok rebuild`    "))
			fmt.Println()
		}
	},
}

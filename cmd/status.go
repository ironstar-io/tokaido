package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
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

		if conf.GetConfig().Global.Syncservice == "unison" {
			ok := unison.CheckBackgroundService(conf.GetConfig().Tokaido.Project.Name)
			if ok {
				console.Println(`🙂  Background sync service is running`, "√")
			} else {
				fmt.Println(aurora.Red(`😓  The Unison background sync service is not running    `))
				fmt.Println()
				pn := conf.GetConfig().Tokaido.Project.Name
				switch system.CheckOS() {
				case "macos":
					h := fs.HomeDir()
					lp := filepath.Join(h, "Library/Logs/tokaido.sync."+pn)
					fmt.Printf("You can check Unison logs in the...\n%s and \n%s files\n", aurora.Bold(lp+".out"), aurora.Bold(lp+".err"))
					fmt.Printf("Or you can run %s to reconfigure and restart it\n", "tok up")
				case "linux":
					sn := "tokaido-sync-" + pn + ".service"
					fmt.Printf("You can check Unison logs by running 'journalctl -u %s'", sn)
				}
				fmt.Println()
				os.Exit(1)
			}
		}

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

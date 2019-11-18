package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// ExecCmd - `tok exec`
var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command inside the Tokaido SSH container",
	Long:  "Execute a command inside the Tokaido SSH container. Alias to `ssh <project-name>.tok -C command`",
	Run: func(cmd *cobra.Command, args []string) {
		if system.CheckOS() == "windows" {
			fmt.Println(aurora.Red("Sorry! The 'tok exec' command isn't available on Windows. Please use SSH instead"))
			os.Exit(1)
		}

		initialize.TokConfig("exec")
		utils.CheckCmdHard("docker-compose")
		telemetry.SendCommand("exec")

		docker.HardCheckTokCompose()

		ssh.StreamConnectCommand(args)
	},
}

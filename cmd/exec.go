package cmd

import (
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// ExecCmd - `tok exec`
var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command inside the Tokaido shell (Drush) container",
	Long:  "Execute a command inside the Tokaido shell (Drush) container using SSH. Alias to `ssh <project-name>.tok -C command`",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("exec")
		utils.CheckCmdHard("docker-compose")
		telemetry.SendCommand("exec")

		docker.HardCheckTokCompose()

		ssh.StreamConnectCommand(args)
	},
}

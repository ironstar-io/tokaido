package cmd

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/spf13/cobra"
)

// PortsCmd - `tok ports x`
var PortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Display Docker port bindings",
	Long:  "Display all or a single local port binding for a Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		docker.PrintPorts(args)
	},
}

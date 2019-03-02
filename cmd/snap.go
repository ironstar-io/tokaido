package cmd

import (
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/snapshots"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

var nameFlag string

// SnapshotNewCmd - `tok snap new`
var SnapshotNewCmd = &cobra.Command{
	Use:   "snapshot new",
	Short: "Creates a new database snapshot",
	Long:  "Creates a new database snapshot and saves it to .tok/local/snapshots with the current UTC date",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("tokaido")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		snapshots.New(nameFlag)
	},
}

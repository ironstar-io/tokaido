package cmd

import (
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// SurveyCmd - `tok survey`
var SurveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "Take a 3 question survey to share your Tokaido feedback",
	Long:  "Take a 3 question survey to share your Tokaido feedback",
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.SendCommand("survey")
		utils.CheckCmdHard("docker-compose")

		telemetry.Survey()
	},
}

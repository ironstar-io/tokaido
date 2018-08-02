package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// IronstarCmd - `tok ironstar`
var IronstarCmd = &cobra.Command{
	Use:   "ironstar",
	Short: "Information about Ironstar",
	Long:  "TODO: Pull company information from endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Pull company information from endpoint")
	},
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd - `tok`
var RootCmd = &cobra.Command{
	Use:   "tok",
	Short: "Use Tokaido to bootstrap your Drupal applications",
	Long:  "Easily build out your Drupal application. Built in Go by the team at Ironstar.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("First tok command run!")
	},
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(IronstarCmd)
	RootCmd.AddCommand(UpCmd)
	RootCmd.AddCommand(StatusCmd)
}

// Execute - Root executable
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

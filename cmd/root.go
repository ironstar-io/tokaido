package cmd

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"github.com/spf13/cobra"
)

// RootCmd - `tok`
var rootCmd = cobra.Command{
	Use:   "tok",
	Short: "Use Tokaido to bootstrap your Drupal applications",
	Long:  "Easily build out your Drupal application. Built in Go by the team at Ironstar.",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(IronstarCmd)
	rootCmd.AddCommand(UpCmd)
	rootCmd.AddCommand(DownCmd)
	rootCmd.AddCommand(StatusCmd)
}

// Execute - Root executable
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// RootCmd will setup and return the root command
func RootCmd() *cobra.Command {
	rootCmd.PersistentFlags().StringP("config", "c", "", "Specify the Tokaido config file to use")
	rootCmd.Flags().IntP("port", "p", 5000, "The port to use for unison")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	fmt.Printf("Starting with config: %+v", config)
}

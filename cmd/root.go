package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os"

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
	rootCmd.AddCommand(DestroyCmd)
	rootCmd.AddCommand(OpenCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(SSHCmd)
	rootCmd.AddCommand(SyscheckCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(WatchCmd)
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
	rootCmd.PersistentFlags().StringP("port", "p", "5000", "The port to use for local development")
	rootCmd.PersistentFlags().StringP("project", "j", fs.Basename(), "The name of the project")
	rootCmd.PersistentFlags().StringP("path", "t", fs.WorkDir(), "The project path")
	rootCmd.PersistentFlags().BoolP("force", "", false, "Forcefully skip confirmation prompts with 'yes' response")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	fmt.Printf("Starting with config: %+v", config)
}

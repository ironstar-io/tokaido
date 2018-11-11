package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/system/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd - `tok`
var rootCmd = cobra.Command{
	Use:   "tok",
	Short: "Use Tokaido to bootstrap your Drupal applications",
	Long:  "Easily build out your Drupal application. Built in Go by the team at Ironstar.",
	Run:   run,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(ConfigCmd)
	rootCmd.AddCommand(ConfigGetCmd)
	rootCmd.AddCommand(ConfigSetCmd)
	rootCmd.AddCommand(DestroyCmd)
	rootCmd.AddCommand(ExecCmd)
	rootCmd.AddCommand(HashCmd)
	rootCmd.AddCommand(LogsCmd)
	rootCmd.AddCommand(NewCmd)
	rootCmd.AddCommand(OpenCmd)
	rootCmd.AddCommand(PortsCmd)
	rootCmd.AddCommand(PsCmd)
	rootCmd.AddCommand(PurgeCmd)
	rootCmd.AddCommand(RebuildCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(SyscheckCmd)
	rootCmd.AddCommand(TestCmd)
	rootCmd.AddCommand(UpCmd)
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
	rootCmd.PersistentFlags().BoolP("force", "", false, "Forcefully skip confirmation prompts with 'yes' response")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode, command output is printed to the console")

	OpenCmd.PersistentFlags().BoolVarP(&adminFlag, "admin", "", false, "Create a one-time admin login URL and open")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	if viper.GetBool("version") == true {
		fmt.Printf("v%s\n", version.Get().Version)
	} else {
		fmt.Printf("Tokaido v%s\n\n", version.Get().Version)
		fmt.Println("For help with Tokaido run `tok help` or take a look at our documentation at https://tokaido.io/docs")
	}
}

func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}

	viper.ReadInConfig()
}

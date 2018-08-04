package cmd

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"

	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

	rootCmd.AddCommand(ConfigGetCmd)
	rootCmd.AddCommand(ConfigSetCmd)
	rootCmd.AddCommand(DestroyCmd)
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(IronstarCmd)
	rootCmd.AddCommand(UpCmd)
	rootCmd.AddCommand(DestroyCmd)
	rootCmd.AddCommand(LogsCmd)
	rootCmd.AddCommand(OpenCmd)
	rootCmd.AddCommand(PortsCmd)
	rootCmd.AddCommand(PsCmd)
	rootCmd.AddCommand(RepairCmd)
	rootCmd.AddCommand(ExecCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(SyscheckCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(SyncCmd)
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
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Check the current Tokaido version (base command only)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode, command output is printed to the console")
	rootCmd.PersistentFlags().BoolP("customCompose", "", false, "When this is enabled, Tokaido will not modify the docker-compose.tok.yml file if it exists")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	if viper.GetBool("version") == true {
		fmt.Printf("v%s\n", version.Get().Version)
	} else {
		fmt.Printf("Tokaido v%s\n\n", version.Get().Version)
		fmt.Println("For help with Tokaido run `tok --help` or take a look at our documentation at https://docs.tokaido.io/")
	}
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	createDotTok()

	viper.SetEnvPrefix("TOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("Tokaido.Customcompose", viper.GetBool("customCompose"))
	viper.SetDefault("Tokaido.Debug", viper.GetBool("debug"))
	viper.SetDefault("Tokaido.Force", viper.GetBool("force"))
	viper.SetDefault("Tokaido.Betacontainers", false)
	viper.SetDefault("Tokaido.Dependencychecks", true)
	viper.SetDefault("Tokaido.Enableemoji", emojiDefaults())
	viper.SetDefault("Tokaido.Project.Name", fs.Basename())
	viper.SetDefault("Tokaido.Project.Path", fs.WorkDir())
	viper.SetDefault("System.Syncsvc.Enabled", true)

	if runtime.GOOS == "linux" {
		viper.SetDefault("System.Syncsvc.Systemdpath", filepath.Join(fs.HomeDir(), "/.config/systemd/user/"))
	}
	if runtime.GOOS == "darwin" {
		viper.SetDefault("System.Syncsvc.Launchdpath", filepath.Join(fs.HomeDir(), "/Library/LaunchAgents/"))
	}

	viper.SetDefault("Services.Memcache.Enabled", true)
	viper.SetDefault("Services.Solr.Enabled", false)

	viper.SetConfigType("yaml")

	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(filepath.Join(fs.WorkDir(), ".tok", "local"))
		viper.AddConfigPath(filepath.Join(fs.WorkDir(), ".tok"))
	}

	viper.ReadInConfig()

	// Check and error if trying to pass in invalid values
	_, err := conf.PopulateConfig(new(conf.Config))
	if err != nil {
		log.Fatalln("Unable to load your configuration\n", err)
	}
}

func emojiDefaults() bool {
	if runtime.GOOS == "windows" {
		return false
	}

	return true
}

func createDotTok() {
	d := filepath.Join(fs.WorkDir(), ".tok")
	if fs.CheckExists(d) == false {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			fmt.Println("There was an error creating the config directory:", err)
		}
	}
}

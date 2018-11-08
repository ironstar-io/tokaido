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
	rootCmd.AddCommand(ExecCmd)
	rootCmd.AddCommand(HashCmd)
	rootCmd.AddCommand(LogsCmd)
	rootCmd.AddCommand(OpenCmd)
	rootCmd.AddCommand(PortsCmd)
	rootCmd.AddCommand(PsCmd)
	rootCmd.AddCommand(PurgeCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(SyscheckCmd)
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

// LoadConfig loads the config from a file if specified, otherwise from the environment
func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	createDotTok()

	viper.SetEnvPrefix("TOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	pr := fs.ProjectRoot()

	viper.SetDefault("Tokaido.Customcompose", viper.GetBool("customCompose"))
	viper.SetDefault("Tokaido.Debug", viper.GetBool("debug"))
	viper.SetDefault("Tokaido.Force", viper.GetBool("force"))
	viper.SetDefault("Tokaido.Betacontainers", false)
	viper.SetDefault("Tokaido.Dependencychecks", true)
	viper.SetDefault("Tokaido.Enableemoji", emojiDefaults())
	viper.SetDefault("Tokaido.Project.Name", strings.Replace(filepath.Base(pr), ".", "", -1))
	viper.SetDefault("Tokaido.Project.Path", pr)
	viper.SetDefault("System.Syncsvc.Enabled", true)
	viper.SetDefault("System.Proxy.Enabled", true)

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
		viper.AddConfigPath(filepath.Join(pr, ".tok", "local"))
		viper.AddConfigPath(filepath.Join(pr, ".tok"))
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
	d := filepath.Join(fs.ProjectRoot(), ".tok")
	if fs.CheckExists(d) == false {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			fmt.Println("There was an error creating the config directory:", err)
		}
	}
}

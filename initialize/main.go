package initialize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/spf13/viper"
)

// TokConfig - loads the config from a file if specified, otherwise from the environment
func TokConfig() {
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
	viper.SetDefault("Drupal.FilePublicPath", "")
	viper.SetDefault("Drupal.FilePrivatePath", constants.DefaultDrupalPrivatePath)
	viper.SetDefault("Drupal.FileTemporaryPath", constants.DefaultDrupalTemporaryPath)
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

	// if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
	// 	viper.SetConfigFile(configFile)
	// } else {
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(pr, ".tok", "local"))
	viper.AddConfigPath(filepath.Join(pr, ".tok"))
	// }

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

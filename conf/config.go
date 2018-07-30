package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

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

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(filepath.Join(fs.WorkDir(), ".tok", "local"))
		viper.AddConfigPath(filepath.Join(fs.WorkDir(), ".tok"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return populateConfig(new(Config))
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

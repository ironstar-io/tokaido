package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the application's configuration
type Config struct {
	Port              string `yaml:"port,omitempty"`
	Config            string `yaml:"config,omitempty"`
	Project           string `yaml:"project,omitempty"`
	Path              string `yaml:"path,omitempty"`
	Force             bool   `yaml:"force,omitempty"`
	Debug             bool   `yaml:"debug,omitempty"`
	Version           bool   `yaml:"version,omitempty"`
	CustomCompose     bool   `yaml:"customcompose,omitempty"`
	SystemdPath       string `yaml:"systemdpath,omitempty"`
	LaunchdPath       string `yaml:"launchdpath,omitempty"`
	CreateSyncService bool   `yaml:"createsyncservice"`
	Drupal            struct {
		Path string `yaml:"path,omitempty"`
	} `yaml:"drupal,omitempty"`
	Xdebug struct {
		Port string `yaml:"port,omitempty"`
	} `yaml:"xdebug,omitempty"`
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("TOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("CreateSyncService", true)
	viper.SetDefault("CustomCompose", false)
	viper.SetDefault("SystemdPath", fs.HomeDir()+"/.config/systemd/user/")
	viper.SetDefault("LaunchdPath", fs.HomeDir()+"/Library/LaunchAgents/")
	viper.SetConfigType("yaml")

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./.tok/")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/.tok/")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return populateConfig(new(Config))
}

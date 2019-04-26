package initialize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// TokConfig - loads the config from a file if specified, otherwise from the environment
func TokConfig(command string) {
	createDotTok()
	createGlobalDotTok()

	viper.SetEnvPrefix("TOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	readProjectConfig(command)
	readGlobalConfig()
}

func emojiDefaults() bool {
	return true
}

func readGlobalConfig() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalln("Unable to resolve home directory so can't initialise Tokaido. Sorry!")
	}

	viper.SetDefault("Global.Syncservice", "fusion")

	utils.DebugString("looking for global config")

	// Check if the global config file exist, and read it in if it does
	gc := filepath.Join(h, ".tok/config.yml")
	_, err = os.Stat(gc)
	if err == nil {
		utils.DebugString("merging in global config file")
		viper.SetConfigFile(gc)
		err = viper.MergeInConfig()
		if err != nil {
			log.Fatalf("Unrecoverable error merging in global config file: %v", err)
		}
	}

	// Check and error if trying to pass in invalid values
	_, err = conf.PopulateConfig(new(conf.Config))
	if err != nil {
		log.Fatalln("Error parsing global configuration file\n", err)
	}
}

func readProjectConfig(command string) {
	utils.DebugString("reading project config")
	pr := fs.ProjectRoot()

	viper.SetDefault("Tokaido.Customcompose", viper.GetBool("customCompose"))
	viper.SetDefault("Tokaido.Debug", viper.GetBool("debug"))
	viper.SetDefault("Tokaido.Force", viper.GetBool("force"))
	viper.SetDefault("Tokaido.Yes", viper.GetBool("yes"))
	viper.SetDefault("Tokaido.Stability", "edge")
	viper.SetDefault("Tokaido.Dependencychecks", true)
	viper.SetDefault("Tokaido.Enableemoji", emojiDefaults())
	viper.SetDefault("Tokaido.Phpversion", "7.1")
	viper.SetDefault("Tokaido.Xdebug", false)
	viper.SetDefault("Tokaido.Project.Name", strings.Replace(filepath.Base(pr), ".", "", -1))
	viper.SetDefault("Tokaido.Project.Path", pr)
	viper.SetDefault("Drupal.FilePublicPath", "")
	viper.SetDefault("Drupal.FilePrivatePath", constants.DefaultDrupalPrivatePath)
	viper.SetDefault("Drupal.FileTemporaryPath", constants.DefaultDrupalTemporaryPath)
	viper.SetDefault("System.Syncsvc.Enabled", true)
	viper.SetDefault("System.Proxy.Enabled", true)

	if command == "new" {
		viper.SetDefault("Services.Adminer.Enabled", true)
		viper.SetDefault("Services.Mailhog.Enabled", true)
	}
	viper.SetDefault("Services.Memcache.Enabled", true)
	viper.SetDefault("Services.Solr.Enabled", false)

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(pr, ".tok"))

	viper.ReadInConfig()

	// Check and error if trying to pass in invalid values
	_, err := conf.PopulateConfig(new(conf.Config))
	if err != nil {
		log.Fatalln("Unable to load your configuration\n", err)
	}
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

func createGlobalDotTok() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalln("Unable to resolve home directory, unable to initialise. Sorry!")
	}

	d := filepath.Join(h, ".tok")
	if fs.CheckExists(d) == false {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			fmt.Println("There was an error creating the global config directory:", err)
		}
	}
}

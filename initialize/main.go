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

	. "github.com/logrusorgru/aurora"
)

// TokConfig - loads the config from a file if specified, otherwise from the environment
func TokConfig(command string) {
	createDotTok()
	createGlobalDotTok()

	viper.SetEnvPrefix("TOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	removeOldGlobalConfig()
	readProjectConfig(command)
	// readGlobalConfig()
}

func emojiDefaults() bool {
	return true
}

// Tokaido 1.7 used ~/.tok/config.yml as the global config file, which confused users
// when talking to them about "config.yml". 1.8 changes this to ~/.tok/global.yml.
// removeOldGlobalConfig will remove the old config file and advise the user.
func removeOldGlobalConfig() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalln("Unable to resolve home directory so can't initialise Tokaido. Sorry!")
	}

	oc := filepath.Join(h, ".tok/config.yml")
	if fs.CheckExists(oc) {
		fs.Remove(oc)
		fmt.Println(Magenta("Tokaido has removed your legacy global config file in $HOME/.tok/config.yml. You don't need it anymore"))
	}
}

func readProjectConfig(command string) {
	utils.DebugString("reading project config")
	pr := fs.ProjectRoot()

	viper.SetDefault("Global.Syncservice", "docker")

	viper.SetDefault("Tokaido.Customcompose", viper.GetBool("customCompose"))
	viper.SetDefault("Tokaido.Debug", viper.GetBool("debug"))
	viper.SetDefault("Tokaido.Force", viper.GetBool("force"))
	viper.SetDefault("Tokaido.Yes", viper.GetBool("yes"))
	viper.SetDefault("Tokaido.Stability", "edge")
	viper.SetDefault("Tokaido.Dependencychecks", true)
	viper.SetDefault("Tokaido.Enableemoji", emojiDefaults())
	viper.SetDefault("Tokaido.Phpversion", "7.2")
	viper.SetDefault("Tokaido.Xdebug", false)
	viper.SetDefault("Tokaido.Project.Name", strings.Replace(filepath.Base(pr), ".", "", -1))
	viper.SetDefault("Drupal.FilePublicPath", "")
	viper.SetDefault("Drupal.FilePrivatePath", constants.DefaultDrupalPrivatePath)
	viper.SetDefault("Drupal.FileTemporaryPath", constants.DefaultDrupalTemporaryPath)
	viper.SetDefault("System.Syncsvc.Enabled", true)
	viper.SetDefault("System.Proxy.Enabled", true)

	if command == "new" {
		viper.SetDefault("Services.Adminer.Enabled", true)
	}
	viper.SetDefault("Services.Memcache.Enabled", false)
	viper.SetDefault("Services.Mailhog.Enabled", true)
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

package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/system/hash"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var configPath = filepath.Join(fs.WorkDir(), "/.tok/config.yml")

// SetConfigValueByArgs ...
func SetConfigValueByArgs(args []string) {
	yt := argsToYaml(args)
	cp := getConfigPath()

	c := GetConfig()

	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file: %v", err)
	}

	err = yaml.Unmarshal(yf, &c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file: %v", err)
	}

	err = yaml.Unmarshal([]byte(yt), &c)
	if err != nil {
		log.Fatalf("There was an issue updating your config file: %v", err)
	}

	// Stop these values leaking into config
	c.Tokaido.Debug = false
	c.Tokaido.Force = false

	fc, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("There was an issue building your config file: %v", err)
	}

	fs.Replace(cp, fc)

	compareFiles(yf, cp)
}

func compareFiles(original []byte, newPath string) {
	o, err := hash.BytesMD5(original)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file: %v", err)
	}

	n, err := hash.FileMD5(newPath)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file: %v", err)
	}

	if o == n {
		fmt.Println("The operation completed without error, but the file is unchanged")
		fmt.Println("Are you sure you selected the right values?")

		return
	}

	fmt.Println("Config updated successfully")
}

func unmarshalConfig(cp string) *Config {
	c := &Config{}

	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file: %v", err)
	}

	err = yaml.Unmarshal(yf, c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file: %v", err)
	}

	return c
}

func getConfigPath() string {
	vc := viper.ConfigFileUsed()
	if vc != "" {
		return vc
	}

	var _, errf = os.Stat(configPath)
	if os.IsNotExist(errf) {
		console.Println(`🏯  Generating a new .tok/config.yml file`, "")
		fs.TouchEmpty(configPath)
	}

	return configPath
}

func argsToYaml(args []string) string {
	var y string
	for i, a := range args {
		if i == len(args)-1 {
			y = y + " " + a
			continue
		}
		y = y + calcWhitespace(i) + a + ":"
	}

	return strings.ToLower(y)
}

func calcWhitespace(i int) string {
	if i == 0 {
		return ""
	}

	w := "\n"
	for x := 1; x <= i; x++ {
		w = w + "  "
	}

	return w
}

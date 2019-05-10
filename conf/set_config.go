package conf

import (
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hash"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	. "github.com/logrusorgru/aurora"
)

// SetConfigValueByArgs updates the config file by merging a single new value with the
// current in memory configuration.
// - args are a slice of new values such as `[]string{"nginx", "workerconnections", "30"}`
// - file is either 'project' or 'global' and will determine which file is updated
func SetConfigValueByArgs(args []string, configFile string) {
	if configFile != "project" && configFile != "global" {
		fmt.Println(Sprintf("The config file %s is unknown", Bold(configFile)))
		os.Exit(1)
	}

	validateArgs(args)

	yt := argsToYaml(args)
	cp := getConfigPath(configFile)

	// 'c' initially carries in-memory config from Viper, which does not differentiate
	//  between our "project" and "global" config files
	// later on we merge our yaml config into this in-memory config
	c := GetConfig()

	// Read the saved config file from disk
	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file\n%v", err)
	}

	// Unmarshal the saved config
	err = yaml.Unmarshal(yf, &c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file\n%v", err)
	}

	// Merge the saved config into our in-memory config
	err = yaml.Unmarshal([]byte(yt), &c)
	if err != nil {
		log.Fatalf("There was an issue updating your config file\n%v", err)
	}

	// Viper doesn't split config in memory so 'c' now contains merged
	// project and global config settings. We need to split them out.
	if configFile == "project" {
		// These values must not be written to the project config file so we reset them to nil or empty
		c.Global.Syncservice = ""
	}

	// Now that we've merged the config, we'll write that merged config to disk
	fc, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("There was an issue building your config file\n%v", err)
	}

	fs.Replace(cp, fc)

	compareFiles(yf, cp)
}

// compareFiles checks original and new config files to identify if any values were changed
func compareFiles(original []byte, newPath string) {
	o, err := hash.BytesMD5(original)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file\n%v", err)
	}

	n, err := hash.FileMD5(newPath)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file:\n%v", err)
	}

	if o == n {
		fmt.Println("Action resulted in no change to config")
		return
	}

	fmt.Println("Config updated successfully")
}

func unmarshalConfig(cp string) *Config {
	c := &Config{}

	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file\n%v", err)
	}

	err = yaml.Unmarshal(yf, c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file\n%v", err)
	}

	return c
}

func getConfigPath(configFile string) string {
	var cp string
	if configFile == "project" {
		cp = filepath.Join(GetConfig().Tokaido.Project.Path, "/.tok/config.yml")
	} else if configFile == "global" {
		cp = filepath.Join(fs.HomeDir(), "/.tok/global.yml")
	} else {
		fmt.Println(Sprintf("The config file %s is unknown", Bold(configFile)))
	}

	var _, errf = os.Stat(cp)
	if os.IsNotExist(errf) {
		fs.TouchEmpty(cp)
	}

	return cp
}

func argsToYaml(args []string) string {
	var y string
	for i, a := range args {
		if i == len(args)-1 {
			y = y + " " + a
			continue
		}
		y = y + calcWhitespace(i) + mapEdgeKeys(a) + ":"
	}

	return y
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

func mapEdgeKeys(a string) string {
	var keyMap = map[string]string{
		"volumesfrom": "volumes_from",
		"dependson":   "depends_on",
		"workingdir":  "working_dir",
	}

	if keyMap[a] != "" {
		return keyMap[a]
	}

	return a
}

func validateArgs(args []string) {
	if len(args) < 2 {
		log.Fatal("At least two arguments are required in order to set a config value")
	}

	ca := args[:len(args)-1]

	_, err := GetConfigValueByArgs(ca)
	if err != nil {
		log.Fatal(err)
	}
}

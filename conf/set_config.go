package conf

import (
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hash"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// SetConfigValueByArgs ...
func SetConfigValueByArgs(args []string) {
	validateArgs(args)

	yt := argsToYaml(args)
	cp := getConfigPath()

	c := GetConfig()

	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file\n%v", err)
	}

	err = yaml.Unmarshal(yf, &c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file\n%v", err)
	}

	err = yaml.Unmarshal([]byte(yt), &c)
	if err != nil {
		log.Fatalf("There was an issue updating your config file\n%v", err)
	}

	// Stop these values leaking into config
	c.Tokaido.Debug = false
	c.Tokaido.Force = false
	c.Tokaido.Yes = false
	c.Tokaido.Project.Path = ""
	c.System.Syncsvc.Launchdpath = ""
	c.System.Syncsvc.Systemdpath = ""

	fc, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("There was an issue building your config file\n%v", err)
	}

	fs.Replace(cp, fc)

	compareFiles(yf, cp)
}

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

func getConfigPath() string {
	cp := filepath.Join(GetConfig().Tokaido.Project.Path, "/.tok/config.yml")

	var _, errf = os.Stat(cp)
	if os.IsNotExist(errf) {
		console.Println(`🏯  Generating a new .tok/config.yml file`, "")
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

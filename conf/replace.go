package conf

import (
	"fmt"

	"github.com/ironstar-io/tokaido/system/fs"

	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func drupalVars(path, version string) []byte {
	return []byte(`drupal:
  path: ` + path + `
  majorVersion: ` + version)
}

// CreateOrReplaceDrupalVars ...
func CreateOrReplaceDrupalVars(path, version string) {
	viper.Set("drupal.path", path)
	viper.Set("drupal.majorversion", version)
	cf := viper.ConfigFileUsed()
	if cf == "" {
		fs.TouchByteArray(filepath.Join(fs.WorkDir(), "/.tok/config.yml"), drupalVars(path, version))
		return
	}

	replaceDrupalVars(cf, path, version)
}

// replaceDrupalVars ...
func replaceDrupalVars(cf, path, version string) {
	confStruct := Config{}
	confStruct.System.Syncsvc.Enabled = true

	iocf, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	errTwo := yaml.Unmarshal(iocf, &confStruct)
	if errTwo != nil {
		log.Fatalf("error: %v", errTwo)
	}

	confStruct.Drupal.Path = path
	confStruct.Drupal.Majorversion = version

	confYml, err := yaml.Marshal(&confStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fs.Replace(cf, confYml)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`There was an error reading your config file. This command may need to be run again.`)
	}
}

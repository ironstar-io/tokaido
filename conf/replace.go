package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"io/ioutil"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func drupalPath(path string) []byte {
	return []byte(`drupal:
  path: ` + path)
}

// CreateOrReplaceDrupalPath ...
func CreateOrReplaceDrupalPath(path string) {
	viper.Set("drupal.path", path)
	cf := viper.ConfigFileUsed()
	if cf == "" {
		fs.TouchByteArray(fs.WorkDir()+"/.tok/config.yml", drupalPath(path))
		return
	}

	replaceDrupalPath(cf, path)
}

// replaceDrupalPath ...
func replaceDrupalPath(cf string, path string) {
	confStruct := Config{
		CreateSyncService: true,
	}
	iocf, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	errTwo := yaml.Unmarshal(iocf, &confStruct)
	if errTwo != nil {
		log.Fatalf("error: %v", errTwo)
	}

	confStruct.Drupal.Path = path

	confYml, err := yaml.Marshal(&confStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fs.Replace(cf, confYml)
}

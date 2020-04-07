package auth

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"

	yaml "gopkg.in/yaml.v2"
)

func ReadInGlobals() (conf.Global, error) {
	globals := conf.Global{}
	gp := filepath.Join(fs.HomeDir(), ".tok", "global.yml")

	err := SafeTouchConfigYAML("global")
	if err != nil {
		return globals, err
	}

	gBytes, err := ioutil.ReadFile(gp)
	if err != nil {
		return globals, err
	}

	err = yaml.Unmarshal(gBytes, &globals)
	if err != nil {
		return globals, err
	}

	return globals, nil
}

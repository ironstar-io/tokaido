package conf

import (
	"github.com/spf13/viper"
)

var tokConfPath = ""

// ReplaceDrupalPath ...
func ReplaceDrupalPath(path string) {
	viper.Set("drupal.path", path)
	viper.WriteConfig()
}

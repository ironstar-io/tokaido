package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"github.com/spf13/viper"
)

func drupalPath(path string) []byte {
	return []byte(`drupal:
	path: ` + path)
}

// ReplaceDrupalPath ...
func ReplaceDrupalPath(path string) {
	viper.Set("drupal.path", path)
	cf := viper.ConfigFileUsed()
	if cf == "" {
		fs.TouchByteArray(fs.WorkDir()+"/.tok/config.yml", drupalPath(path))
		return
	}

}

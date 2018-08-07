package conf

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/spf13/viper"
)

func drupalVars(path, version string) []byte {
	return []byte(`drupal:
  path: ` + path + `
  majorVersion: ` + version)
}

func commandSubstitution(name string, args ...string) {

	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, _ := cmd.CombinedOutput()

	fmt.Println(strings.TrimSpace(string(stdoutStderr)))

}

// CreateOrReplaceDrupalConfig ...
func CreateOrReplaceDrupalConfig(path, version string) {
	viper.Set("drupal.path", path)
	viper.Set("drupal.majorversion", version)
	cf := viper.ConfigFileUsed()

	if cf == "" {
		fs.TouchByteArray(filepath.Join(fs.WorkDir(), "/.tok/config.yml"), drupalVars(path, version))
	}
}

package conf

import (
	"fmt"
	"os"
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
	cmd.Dir = GetConfig().Tokaido.Project.Path
	stdoutStderr, _ := cmd.CombinedOutput()

	fmt.Println(strings.TrimSpace(string(stdoutStderr)))
}

// CreateOrReplaceDrupalConfig ...
func CreateOrReplaceDrupalConfig(path, version string) {
	viper.Set("drupal.path", path)
	viper.Set("drupal.majorversion", version)

	configPath := filepath.Join(GetConfig().Tokaido.Project.Path, "/.tok/config.yml")

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		fs.TouchByteArray(configPath, drupalVars(path, version))
	}
}

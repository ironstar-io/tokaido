package nightwatch

import (
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
)

func calcSiteURL() string {
	if proxy.CheckProxyUp() == true {
		return proxy.GetProxyURL()
	}

	return "https://localhost:" + docker.LocalPort("haproxy", "8443")
}

// Copy and amend .env file if it doesn't already exist
func configureEnvFile() {
	env := filepath.Join(conf.CoreDrupal8Path(), ".env")
	if fs.CheckExists(env) == false {
		generateEnvFile(env)
	}
}

func yarnInstall() error {
	cp := conf.CoreDrupal8Path()

	fmt.Println("    Ensuring Nightwatch dependencies are installed")
	ssh.StreamConnectCommand([]string{"cd", cp, "&&", "yarn", "install"})

	// _, err := utils.CommandSubSplitOutputContext(d, "yarn", "check", "--verify-tree")
	// if err != nil {
	// 	fmt.Println()
	// 	console.Println("🧶   Nightwatch dependencies haven't been installed yet. Attempting to install with yarn", "")
	// 	utils.StreamOSCmdContext(d, "yarn", "install")
	// }

	return nil
}

func yarnTestNightwatch() {
	cp := conf.CoreDrupal8Path()

	console.Println("👩‍💻  Tokaido is starting a Nightwatch test run with the command `yarn test:nightwatch`\n", "")

	ssh.StreamConnectCommand([]string{"cd", cp, "&&", "yarn", "test:nightwatch"})

}

func checkCompatibility() error {
	goos.CheckDependencies()
	return nil
}

func enableChromedriver() {
	conf.SetConfigValueByArgs([]string{"services", "chromedriver", "enabled", "true"}, "project")
}

// RunDrupalTests - Run Drupal Nightwatch tests. See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#nightwatch-tests
func RunDrupalTests() error {
	err := CheckLocalDrupal()
	if err != nil {
		console.Println("⚠️   "+err.Error(), "")
		return err
	}

	err = checkCompatibility()
	if err != nil {
		return err
	}

	enableChromedriver()

	checkPHPUnit()
	yarnInstall()

	configureEnvFile()

	yarnTestNightwatch()

	return nil
}

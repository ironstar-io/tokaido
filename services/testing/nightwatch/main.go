package nightwatch

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

func updateEnvFile() {

}

// Copy and amend .env file if it doesn't already exist
func configureEnvFile() {
	env := filepath.Join(conf.CoreDrupal8Path(), ".env")
	if fs.CheckExists(env) == false {
		fs.Copy(env+".example", env)
		updateEnvFile()
	}
}

func yarnInstall() error {
	d := conf.CoreDrupal8Path()

	_, err := utils.CommandSubSplitOutputContext(d, "yarn", "check", "--verify-tree")
	if err != nil {
		fmt.Println()
		console.Println("🧶  Nightwatch dependencies missing! Attempting to install with yarn", "")
		utils.StreamOSCmdContext(d, "yarn", "install")
	}

	return nil
}

func yarnTestNightwatch() {
	d := conf.CoreDrupal8Path()

	console.Println("👩‍💻  Tokaido is starting a Nightwatch test run with the command `yarn test:nightwatch`\n", "")
	utils.StreamOSCmdContext(d, "yarn", "test:nightwatch")
}

func upgradePHPUnit() {
	// TODO: Check that a compatible version of phpunit is being used
	// TODO: Check that user is OK with upgrading their version of phpunit before commencing.
	// ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "run-script", "drupal-phpunit-upgrade"})
}

func checkCompatibility() error {
	goos.CheckDependencies()
	return nil
}

// RunDrupalTests - Run Drupal Nightwatch tests. See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#nightwatch-tests
func RunDrupalTests() error {
	err := CheckLocalDrupal()
	if err != nil {
		console.Println("⚠️    "+err.Error(), "")
		log.Fatalf("Local Drupal checks failed. Unable to continue.")
	}

	err = checkCompatibility()
	if err != nil {
		return err
	}

	upgradePHPUnit()
	yarnInstall()

	configureEnvFile()

	yarnTestNightwatch()

	return nil
}
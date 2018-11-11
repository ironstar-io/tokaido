package nightwatch

import (
	"fmt"
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch/goos"
	"github.com/ironstar-io/tokaido/system/console"
)

func yarnInstall() error {
	d := conf.CoreDrupal8Path()

	// yarn check --verify-tree
	// if exit code > 0
	// yarn install

	return nil
}

func upgradePHPUnit() {
	// TODO: Check that a compatible version of phpunit is being used
	// TODO: Check that user is OK with upgrading their version of phpunit before commencing.
	// ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "run-script", "drupal-phpunit-upgrade"})
}

func checkCompatibility() error {
	fmt.Println("\nNote that you will need to have Google Chrome installed in order to be able to run Drupal Nightwatch tests.")

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

	return nil
}

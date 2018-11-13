package nightwatch

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/ironstar-io/tokaido/services/composer"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
)

var validPHPUnitRange = ">=6.x.x"
var phpUnit = "phpunit/phpunit"

func upgradePHPUnit(currentVersion string) {
	console.Println("🐘  Tokaido detected you're running version "+currentVersion+" of PHPUnit, which isn't compatible with Drupal Nightwatch.\n", "")
	confirmUpdate := utils.ConfirmationPrompt("Would you like us to attempt to update it for you automatically with composer?", "y")
	if confirmUpdate == false {
		fmt.Println("We'll still attempt to run the tests, but they make not behave as expected!")
		return
	}

	ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "update", phpUnit})
}

func installPHPUnit() {
	console.Println("🐘  Tokaido was unable to detect an installation of PHPUnit, which is required for Drupal Nightwatch.\n", "")
	confirmUpdate := utils.ConfirmationPrompt("Would you like us to attempt to install it for you automatically with composer?", "y")
	if confirmUpdate == false {
		fmt.Println("We'll still attempt to run the tests, but they make not behave as expected!")
		return
	}

	// TODO: Check that user is OK with installing phpunit before commencing.
	ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "install", phpUnit})
}

func checkPHPUnit() {
	v, err := composer.FindPackageVersion(phpUnit)
	if err != nil {
		installPHPUnit()
		return
	}
	if v == "" {
		installPHPUnit()
		return
	}

	ps, _ := semver.Parse(v)
	validRange, _ := semver.ParseRange(validPHPUnitRange)
	if !validRange(ps) {
		upgradePHPUnit(v)
		return
	}
}

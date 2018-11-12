package nightwatch

import (
	"github.com/blang/semver"
	"github.com/ironstar-io/tokaido/services/composer"
	"github.com/ironstar-io/tokaido/system/ssh"
)

var validPHPUnitRange = ">=6.x.x"

func upgradePHPUnit() {
	// TODO: Check that user is OK with upgrading their version of phpunit before commencing.
	ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "update", "phpunit/phpunit"})
}

func installPHPUnit() {
	// TODO: Check that user is OK with installing phpunit before commencing.
	ssh.StreamConnectCommand([]string{"cd", "/tokaido/site;", "composer", "install", "phpunit/phpunit"})
}

func checkPHPUnit() {
	d := ssh.ConnectCommand([]string{"cd", "/tokaido/site;", "composer", "show", "phpunit/phpunit"})
	if d == "" {
		installPHPUnit()
		return
	}

	v := composer.FindPackageVersion(d)
	if v == "" {
		installPHPUnit()
		return
	}

	ps, _ := semver.Parse(v)
	validRange, _ := semver.ParseRange(validPHPUnitRange)
	if !validRange(ps) {
		upgradePHPUnit()
		return
	}
}

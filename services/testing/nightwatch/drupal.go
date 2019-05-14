package nightwatch

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/blang/semver"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/fs"
)

var rgx = regexp.MustCompile("'(.*?)'")
var validDrupalRange = ">=8.6.x"

// CheckLocalDrupal ...
func CheckLocalDrupal() error {
	if fs.CheckExists(filepath.Join(conf.GetProjectPath(), conf.CoreDrupalFile())) == false {
		return errors.New("A valid Drupal installation was not found. Exiting...")
	}

	e := checkDrupalVersion()
	if e != nil {
		return e
	}

	return nil
}

// checkDrupalVersion ...
func checkDrupalVersion() error {
	drupalVersion := drupal.GetDrupalVersion()
	if drupalVersion == "" {
		return errors.New("Tokaido was unable to determine the Drupal version, it needs to be " + validDrupalRange + " in order to run Nightwatch tests")
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		return nil
	}

	return errors.New("The Drupal version (" + drupalVersion + ") is not supported for Nightwatch tests. Drupal version must be " + validDrupalRange + ". See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#nightwatch-tests for more information.")
}

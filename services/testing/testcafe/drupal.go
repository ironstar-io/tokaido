package testcafe

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/blang/semver"
)

var rgx = regexp.MustCompile("'(.*?)'")
var validDrupalRange = ">=8.0.x"

// CheckLocalDrupal ...
func CheckLocalDrupal() error {
	if fs.CheckExists(filepath.Join(conf.GetProjectPath(), conf.CoreDrupalFile())) == false {
		return errors.New("a valid Drupal installation was not found")
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
		return errors.New("Tokaido was unable to determine the Drupal version, it needs to be " + validDrupalRange + " in order to run TestCafe tests")
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		return nil
	}

	return errors.New("The Drupal version (" + drupalVersion + ") is not supported for TestCafe tests. Drupal version must be " + validDrupalRange + ".")
}

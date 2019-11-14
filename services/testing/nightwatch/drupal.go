package nightwatch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/blang/semver"
	"github.com/logrusorgru/aurora"
)

var rgx = regexp.MustCompile("'(.*?)'")
var validDrupalRange = ">=8.6.x"

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
		return errors.New("Tokaido was unable to determine the Drupal version, it needs to be " + validDrupalRange + " in order to run Nightwatch tests")
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		return nil
	}

	return errors.New("The Drupal version (" + drupalVersion + ") is not supported for Nightwatch tests. Drupal version must be " + validDrupalRange + ". See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#nightwatch-tests for more information.")
}

// checkNightwatchConfig scans nightwatch config files in Drupal Core to see
// if they contain settings that we need. If not, we provide info for the user
// to manually fix those, until such time as we have Drupal patches ready
func checkNightwatchConfig() error {
	ok := true // If false at the end of checking, the program will exit with an error
	pr := conf.GetProjectPath()
	dr := conf.GetConfig().Drupal.Path

	nightwatchConfPath := filepath.Join(pr, dr, "core", "tests", "Drupal", "Nightwatch", "nightwatch.conf.js")
	utils.DebugString("nightwatchConfPath: " + nightwatchConfPath)
	if !fs.CheckExists(nightwatchConfPath) {
		fmt.Println(aurora.Red("Expected nightwatch configuration file was not found: "), nightwatchConfPath)
		fmt.Println(aurora.Red("Tokaido cannot continue. Is your Drupal Core installation working?"))
		os.Exit(1)
	}

	nightwatchOk := fs.Contains(nightwatchConfPath, "acceptInsecureCerts: true")
	utils.DebugString("nightwatchOk: " + strconv.FormatBool(nightwatchOk))
	if !nightwatchOk {
		ok = false
		fmt.Println(aurora.Red("Tokaido tests via HTTPS, but Druapl's Nightwatch config only supports HTTP."))
		fmt.Println(aurora.Red("To fix this, you need to add 'acceptInsecureCerts: true' to the desiredCapabilities "))
		fmt.Println(aurora.Red("section of "), nightwatchConfPath)
		fmt.Println()
		fmt.Println(aurora.Yellow("Please see docs.tokaido.io for more information"))
		fmt.Println()
	}

	drupalUserIsLoggedInPath := filepath.Join(pr, dr, "core", "tests", "Drupal", "Nightwatch", "Commands", "drupalUserIsLoggedIn.js")
	utils.DebugString("drupalUserIsLoggedInPath: " + drupalUserIsLoggedInPath)
	if !fs.CheckExists(drupalUserIsLoggedInPath) {
		fmt.Println(aurora.Red("Expected nightwatch command file was not found: "), drupalUserIsLoggedInPath)
		fmt.Println(aurora.Red("Tokaido cannot continue. Is your Drupal Core installation working?"))
		os.Exit(1)
	}

	drupalUserIsLoggedInOk := fs.Contains(drupalUserIsLoggedInPath, "/^SSESS/")
	utils.DebugString("drupalUserIsLoggedInOk: " + strconv.FormatBool(drupalUserIsLoggedInOk))
	if !drupalUserIsLoggedInOk {
		ok = false
		fmt.Println(aurora.Red("Tokaido tests via HTTPS, but Druapl's Nightwatch config only supports HTTP."))
		fmt.Println(aurora.Red("To fix this, you need change the cookie name in"), drupalUserIsLoggedInPath)
		fmt.Println(aurora.Red("from 'SESS' to 'SSESS' (adding one extra 'S') "))
		fmt.Println()
		fmt.Println(aurora.Yellow("Please see docs.tokaido.io for more information"))
		fmt.Println()
	}

	if !ok {
		fmt.Println(aurora.Red("'tok test' cannot continue until Drupal core config is updated"))
		os.Exit(0)
	}

	return nil
}

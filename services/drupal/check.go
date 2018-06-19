package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

var rgx = regexp.MustCompile("'(.*?)'")
var drupalDir = fs.WorkDir() + "/docroot/core/lib/Drupal.php"
var validDrupalRange = ">=8.5.x"
var checkFailMsg = "\n⚠️  There were some problems detected during the system checks. This won't stop you from running any Tokaido commands, but they may not behave as you were expecting.\n\n"

// CheckLocal ...
func CheckLocal() {
	if _, err := os.Stat(drupalDir); os.IsNotExist(err) {
		fmt.Println("  ✘  A Drupal installation was not found")
		fmt.Printf(checkFailMsg)
		return
	}

	fmt.Println("  ✓  A Drupal installation was found")

	checkDrupalVersion()
}

// CheckContainer ...
func CheckContainer() {
	haproxyPort := docker.LocalPort("haproxy", "8443")

	drupalStatus := utils.SilentBashStringCmd(`curl -sko /dev/null -I -w"%{http_code}" https://localhost:` + haproxyPort + ` | grep 200`)
	if drupalStatus == "200" {
		fmt.Println("✅  Drupal is listening on HTTPS")
		return
	}

	fmt.Println(`😓  Drupal is not installed

Tokaido is running but it looks like your Drupal site isn't installed.

You can install Drupal by using the web interface at https://project-name:port_number. Note that your database credentials should be:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido

It might be easier to use Drush to install your site, which you can do by connecting to SSH with 'tok ssh' and running 'drush site-install'`)
	os.Exit(1)
}

// checkDrupalVersion ...
func checkDrupalVersion() {
	drupalVersion := getDrupalVersion()
	if drupalVersion == "" {
		fmt.Printf("  ✘  Tokaido was unable to determine the Drupal version, it should be %s\n", validDrupalRange)
		fmt.Printf(checkFailMsg)
		return
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		fmt.Printf("  ✓  The Drupal version (%s) is supported by Tokaido\n\n", drupalVersion)
	} else {
		fmt.Printf("  ✘  The Drupal version (%s) is not supported by Tokaido. Drupal version %s has been tested and is working with Tokaido\n", drupalVersion, validDrupalRange)
		fmt.Printf(checkFailMsg)
	}
}

// getDrupalVersion ...
func getDrupalVersion() string {
	drupalVersionString := versionStr()
	if drupalVersionString == "" {
		return ""
	}

	drupalVersion := rgx.FindStringSubmatch(drupalVersionString)[0]

	return strings.Replace(drupalVersion, "'", "", 2)
}

// versionStr ...
func versionStr() string {
	f, err := os.Open(drupalDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "const VERSION = '") {
			return scanner.Text()
		}
	}

	return ""
}
